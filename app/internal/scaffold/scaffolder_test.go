package scaffold

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScaffolder_Run(t *testing.T) {
	ctx := context.Background()

	t.Run("successful scaffolding with default repo and target dir", func(t *testing.T) {
		mockGit := &MockGitClient{}
		mockFS := NewMockFS()

		scaffolder := NewScaffolder(mockGit, mockFS)

		opts := Options{
			ProjectName: "my-project",
		}

		err := scaffolder.Run(ctx, opts)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Verify Clone call
		if len(mockGit.CloneCalls) != 1 {
			t.Fatalf("expected 1 clone call, got: %d", len(mockGit.CloneCalls))
		}
		if mockGit.CloneCalls[0].RepoURL != DefaultRepoURL {
			t.Errorf("expected repoURL %s, got: %s", DefaultRepoURL, mockGit.CloneCalls[0].RepoURL)
		}
		if mockGit.CloneCalls[0].TargetDir != "my-project" {
			t.Errorf("expected targetDir 'my-project', got: %s", mockGit.CloneCalls[0].TargetDir)
		}

		// Verify removals: app directory, quickstart.md, and .git directory
		expectedAppDir := filepath.Join("my-project", "app")
		expectedQuickstart := filepath.Join("my-project", "quickstart.md")
		expectedGitDir := filepath.Join("my-project", ".git")

		hasRemovedApp := false
		hasRemovedQuickstart := false
		hasRemovedGit := false
		for _, p := range mockFS.RemovedPaths {
			if p == expectedAppDir {
				hasRemovedApp = true
			}
			if p == expectedQuickstart {
				hasRemovedQuickstart = true
			}
			if p == expectedGitDir {
				hasRemovedGit = true
			}
		}
		if !hasRemovedApp {
			t.Errorf("expected app directory %s to be removed, removed paths: %v", expectedAppDir, mockFS.RemovedPaths)
		}
		if !hasRemovedQuickstart {
			t.Errorf("expected quickstart.md %s to be removed, removed paths: %v", expectedQuickstart, mockFS.RemovedPaths)
		}
		if !hasRemovedGit {
			t.Errorf("expected .git directory %s to be removed, removed paths: %v", expectedGitDir, mockFS.RemovedPaths)
		}

		// Verify Git Init
		if len(mockGit.InitCalls) != 1 {
			t.Fatalf("expected 1 git init call, got: %d", len(mockGit.InitCalls))
		}
		if mockGit.InitCalls[0] != "my-project" {
			t.Errorf("expected git init in 'my-project', got: %s", mockGit.InitCalls[0])
		}

		// Verify README.md content
		expectedReadme := filepath.Join("my-project", "README.md")
		content, ok := mockFS.WrittenFiles[expectedReadme]
		if !ok {
			t.Fatalf("expected %s to be written", expectedReadme)
		}
		expectedContent := "# my-project\n"
		if string(content) != expectedContent {
			t.Errorf("expected README content %q, got: %q", expectedContent, string(content))
		}
	})

	t.Run("successful scaffolding with explicit repo URL and target dir", func(t *testing.T) {
		mockGit := &MockGitClient{}
		mockFS := NewMockFS()

		scaffolder := NewScaffolder(mockGit, mockFS)

		opts := Options{
			ProjectName: "awesome-app",
			TargetDir:   "custom-path/awesome-app",
			RepoURL:     "https://github.com/custom/repo",
		}

		err := scaffolder.Run(ctx, opts)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if mockGit.CloneCalls[0].RepoURL != "https://github.com/custom/repo" {
			t.Errorf("expected custom repo URL, got: %s", mockGit.CloneCalls[0].RepoURL)
		}
		expectedTarget := filepath.Clean("custom-path/awesome-app")
		if mockGit.CloneCalls[0].TargetDir != expectedTarget {
			t.Errorf("expected custom target dir %s, got: %s", expectedTarget, mockGit.CloneCalls[0].TargetDir)
		}
	})

	t.Run("rejects invalid project names", func(t *testing.T) {
		invalidNames := []string{
			"",
			"   \t\n  ",
			"-invalid-flag",
			"--bare",
			"foo/../bar",
			".",
			"..",
			"project\ninjection",
			"project\x00nullbyte",
		}

		for _, name := range invalidNames {
			scaffolder := NewScaffolder(&MockGitClient{}, NewMockFS())
			err := scaffolder.Run(ctx, Options{ProjectName: name})
			if err == nil {
				t.Fatalf("expected error for invalid project name %q, got nil", name)
			}
			if !errors.Is(err, ErrInvalidProjectName) {
				t.Errorf("expected ErrInvalidProjectName for %q, got: %v", name, err)
			}
		}
	})

	t.Run("rejects invalid target directories", func(t *testing.T) {
		invalidDirs := []string{
			"-flag-dir",
			"--option",
			"..",
			"../outside",
			"/",
			".",
		}

		for _, dir := range invalidDirs {
			scaffolder := NewScaffolder(&MockGitClient{}, NewMockFS())
			err := scaffolder.Run(ctx, Options{
				ProjectName: "valid-name",
				TargetDir:   dir,
			})
			if err == nil {
				t.Fatalf("expected error for invalid target dir %q, got nil", dir)
			}
			if !errors.Is(err, ErrInvalidTargetDir) {
				t.Errorf("expected ErrInvalidTargetDir for %q, got: %v", dir, err)
			}
		}
	})

	t.Run("rejects unsafe repo URLs", func(t *testing.T) {
		unsafeURLs := []string{
			"--upload-pack=exploit",
			"-v",
			"ext::sh -c evil",
		}

		for _, url := range unsafeURLs {
			scaffolder := NewScaffolder(&MockGitClient{}, NewMockFS())
			err := scaffolder.Run(ctx, Options{
				ProjectName: "valid-name",
				RepoURL:     url,
			})
			if err == nil {
				t.Fatalf("expected error for unsafe repo URL %q, got nil", url)
			}
			if !errors.Is(err, ErrInvalidRepoURL) {
				t.Errorf("expected ErrInvalidRepoURL for %q, got: %v", url, err)
			}
		}
	})

	t.Run("returns ErrDirectoryNotEmpty when target exists and is not empty", func(t *testing.T) {
		mockFS := NewMockFS()
		mockFS.ExistsFunc = func(path string) (bool, error) {
			return true, nil
		}
		mockFS.IsEmptyDirFunc = func(path string) (bool, error) {
			return false, nil
		}

		scaffolder := NewScaffolder(&MockGitClient{}, mockFS)
		err := scaffolder.Run(ctx, Options{ProjectName: "existing-dir"})
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, ErrDirectoryNotEmpty) {
			t.Errorf("expected ErrDirectoryNotEmpty, got: %v", err)
		}
	})

	t.Run("allows scaffolding if target directory exists and is empty", func(t *testing.T) {
		mockFS := NewMockFS()
		mockFS.ExistsFunc = func(path string) (bool, error) {
			return true, nil
		}
		mockFS.IsEmptyDirFunc = func(path string) (bool, error) {
			return true, nil
		}

		scaffolder := NewScaffolder(&MockGitClient{}, mockFS)
		err := scaffolder.Run(ctx, Options{ProjectName: "empty-dir"})
		if err != nil {
			t.Fatalf("expected success on empty dir, got: %v", err)
		}
	})

	t.Run("handles filesystem exists check failure", func(t *testing.T) {
		mockFS := NewMockFS()
		mockFS.ExistsFunc = func(path string) (bool, error) {
			return false, errors.New("disk failure")
		}

		scaffolder := NewScaffolder(&MockGitClient{}, mockFS)
		err := scaffolder.Run(ctx, Options{ProjectName: "proj"})
		if err == nil || !strings.Contains(err.Error(), "disk failure") {
			t.Fatalf("expected disk failure error, got: %v", err)
		}
	})

	t.Run("handles filesystem is-empty check failure", func(t *testing.T) {
		mockFS := NewMockFS()
		mockFS.ExistsFunc = func(path string) (bool, error) {
			return true, nil
		}
		mockFS.IsEmptyDirFunc = func(path string) (bool, error) {
			return false, errors.New("permission denied")
		}

		scaffolder := NewScaffolder(&MockGitClient{}, mockFS)
		err := scaffolder.Run(ctx, Options{ProjectName: "proj"})
		if err == nil || !strings.Contains(err.Error(), "permission denied") {
			t.Fatalf("expected permission denied error, got: %v", err)
		}
	})

	t.Run("handles git clone failure and rolls back target dir if newly created", func(t *testing.T) {
		mockFS := NewMockFS()
		mockFS.ExistsFunc = func(path string) (bool, error) {
			return false, nil // target does not exist initially
		}
		mockGit := &MockGitClient{
			CloneFunc: func(ctx context.Context, repoURL, targetDir string) error {
				return errors.New("network unreachable")
			},
		}

		scaffolder := NewScaffolder(mockGit, mockFS)
		err := scaffolder.Run(ctx, Options{ProjectName: "proj"})
		if err == nil || !strings.Contains(err.Error(), "network unreachable") {
			t.Fatalf("expected network unreachable error, got: %v", err)
		}

		// Verify target dir was rolled back / removed
		rolledBack := false
		for _, p := range mockFS.RemovedPaths {
			if p == "proj" {
				rolledBack = true
			}
		}
		if !rolledBack {
			t.Errorf("expected target dir 'proj' to be removed on failure rollback, removed paths: %v", mockFS.RemovedPaths)
		}
	})

	t.Run("handles app directory removal failure", func(t *testing.T) {
		mockFS := NewMockFS()
		mockFS.RemoveAllFunc = func(path string) error {
			if strings.HasSuffix(path, "app") {
				return errors.New("cannot remove app")
			}
			return nil
		}

		scaffolder := NewScaffolder(&MockGitClient{}, mockFS)
		err := scaffolder.Run(ctx, Options{ProjectName: "proj"})
		if err == nil || !strings.Contains(err.Error(), "cannot remove app") {
			t.Fatalf("expected cannot remove app error, got: %v", err)
		}
	})

	t.Run("handles quickstart.md removal failure", func(t *testing.T) {
		mockFS := NewMockFS()
		mockFS.RemoveAllFunc = func(path string) error {
			if strings.HasSuffix(path, "quickstart.md") {
				return errors.New("cannot remove quickstart.md")
			}
			return nil
		}

		scaffolder := NewScaffolder(&MockGitClient{}, mockFS)
		err := scaffolder.Run(ctx, Options{ProjectName: "proj"})
		if err == nil || !strings.Contains(err.Error(), "cannot remove quickstart.md") {
			t.Fatalf("expected cannot remove quickstart.md error, got: %v", err)
		}
	})

	t.Run("handles .git directory removal failure", func(t *testing.T) {
		mockFS := NewMockFS()
		mockFS.RemoveAllFunc = func(path string) error {
			if strings.HasSuffix(path, ".git") {
				return errors.New("cannot remove .git")
			}
			return nil
		}

		scaffolder := NewScaffolder(&MockGitClient{}, mockFS)
		err := scaffolder.Run(ctx, Options{ProjectName: "proj"})
		if err == nil || !strings.Contains(err.Error(), "cannot remove .git") {
			t.Fatalf("expected cannot remove .git error, got: %v", err)
		}
	})

	t.Run("handles git init failure", func(t *testing.T) {
		mockGit := &MockGitClient{
			InitFunc: func(ctx context.Context, targetDir string) error {
				return errors.New("git init error")
			},
		}

		scaffolder := NewScaffolder(mockGit, NewMockFS())
		err := scaffolder.Run(ctx, Options{ProjectName: "proj"})
		if err == nil || !strings.Contains(err.Error(), "git init error") {
			t.Fatalf("expected git init error, got: %v", err)
		}
	})

	t.Run("handles README write failure", func(t *testing.T) {
		mockFS := NewMockFS()
		mockFS.WriteFileFunc = func(path string, data []byte, perm os.FileMode) error {
			return errors.New("write failure")
		}

		scaffolder := NewScaffolder(&MockGitClient{}, mockFS)
		err := scaffolder.Run(ctx, Options{ProjectName: "proj"})
		if err == nil || !strings.Contains(err.Error(), "write failure") {
			t.Fatalf("expected write failure error, got: %v", err)
		}
	})
}

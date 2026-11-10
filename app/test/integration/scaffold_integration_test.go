package integration

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"new-dev-ai/internal/scaffold"
)

// setupMockRemoteRepo creates a local git repository that mirrors the template structure.
func setupMockRemoteRepo(t *testing.T) string {
	t.Helper()

	repoDir := t.TempDir()

	// Create files simulating the template repo
	files := map[string]string{
		"README.md":              "# Original Template Readme\nSome description",
		"AGENTS.md":              "# Agent Instructions",
		"quickstart.md":          "# Quickstart Guide",
		".agents/plugin.json":    `{"name": "egm"}`,
		"app/spec.md":            "# Spec",
		"app/main.go":            "package main\nfunc main() {}",
	}

	for relPath, content := range files {
		fullPath := filepath.Join(repoDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("failed to create directory for %s: %v", relPath, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write %s: %v", relPath, err)
		}
	}

	// Initialize git repo and make initial commit
	runCmd(t, repoDir, "git", "init")
	runCmd(t, repoDir, "git", "config", "user.name", "TestUser")
	runCmd(t, repoDir, "git", "config", "user.email", "test@example.com")
	runCmd(t, repoDir, "git", "add", ".")
	runCmd(t, repoDir, "git", "commit", "-m", "initial template commit")

	return repoDir
}

func runCmd(t *testing.T, dir string, name string, args ...string) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %s %v in %s: %v\nOutput: %s", name, args, dir, err, string(output))
	}
}

func TestScaffoldIntegration_LocalRepo(t *testing.T) {
	sourceRepo := setupMockRemoteRepo(t)
	targetBaseDir := t.TempDir()
	targetDir := filepath.Join(targetBaseDir, "my-test-project")

	gitClient, err := scaffold.NewExecGitClient()
	if err != nil {
		t.Fatalf("git client initialization failed: %v", err)
	}

	fs := scaffold.NewOSFS()
	scaffolder := scaffold.NewScaffolder(gitClient, fs)

	opts := scaffold.Options{
		ProjectName: "my-test-project",
		TargetDir:   targetDir,
		RepoURL:     sourceRepo,
	}

	ctx := context.Background()
	if err := scaffolder.Run(ctx, opts); err != nil {
		t.Fatalf("scaffolder.Run failed: %v", err)
	}

	// 1. Target directory must exist
	exists, err := fs.Exists(targetDir)
	if err != nil || !exists {
		t.Fatalf("target directory %s does not exist", targetDir)
	}

	// 2. app directory and quickstart.md must NOT exist
	appDir := filepath.Join(targetDir, "app")
	appExists, _ := fs.Exists(appDir)
	if appExists {
		t.Errorf("expected app/ directory to NOT exist in scaffolded project")
	}

	quickstartFile := filepath.Join(targetDir, "quickstart.md")
	quickstartExists, _ := fs.Exists(quickstartFile)
	if quickstartExists {
		t.Errorf("expected quickstart.md to NOT exist in scaffolded project")
	}

	// 3. git repository must be initialized
	gitDir := filepath.Join(targetDir, ".git")
	gitExists, _ := fs.Exists(gitDir)
	if !gitExists {
		t.Errorf("expected .git directory to exist in %s", targetDir)
	}

	// Verify that it is a valid git repository
	cmd := exec.Command("git", "-C", targetDir, "rev-parse", "--is-inside-work-tree")
	output, err := cmd.CombinedOutput()
	if err != nil || strings.TrimSpace(string(output)) != "true" {
		t.Errorf("target directory is not a valid git work tree: %v, output: %s", err, string(output))
	}

	// 4. README.md must contain only the project name passed as parameter as title
	readmePath := filepath.Join(targetDir, "README.md")
	readmeBytes, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("failed to read README.md: %v", err)
	}
	expectedReadme := "# my-test-project\n"
	if string(readmeBytes) != expectedReadme {
		t.Errorf("README.md content mismatch.\nExpected: %q\nGot:      %q", expectedReadme, string(readmeBytes))
	}

	// 5. Template files (like AGENTS.md, .agents/) must still be present
	for _, relPath := range []string{"AGENTS.md", filepath.Join(".agents", "plugin.json")} {
		p := filepath.Join(targetDir, relPath)
		fileFound, _ := fs.Exists(p)
		if !fileFound {
			t.Errorf("expected template file %s to be preserved", relPath)
		}
	}
}

func TestScaffoldIntegration_RealRemote(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping real remote clone integration test in -short mode")
	}

	targetBaseDir := t.TempDir()
	targetDir := filepath.Join(targetBaseDir, "real-repo-test")

	gitClient, err := scaffold.NewExecGitClient()
	if err != nil {
		t.Fatalf("git client initialization failed: %v", err)
	}

	fs := scaffold.NewOSFS()
	scaffolder := scaffold.NewScaffolder(gitClient, fs)

	opts := scaffold.Options{
		ProjectName: "real-repo-test",
		TargetDir:   targetDir,
		RepoURL:     scaffold.DefaultRepoURL,
	}

	ctx := context.Background()
	if err := scaffolder.Run(ctx, opts); err != nil {
		t.Fatalf("scaffolder.Run failed on real remote: %v", err)
	}

	// Verify app/ is gone
	appDir := filepath.Join(targetDir, "app")
	if exists, _ := fs.Exists(appDir); exists {
		t.Errorf("expected app/ directory to NOT exist")
	}

	// Verify quickstart.md is gone
	quickstartPath := filepath.Join(targetDir, "quickstart.md")
	if exists, _ := fs.Exists(quickstartPath); exists {
		t.Errorf("expected quickstart.md to NOT exist")
	}

	// Verify README.md has solely "# real-repo-test\n"
	readmePath := filepath.Join(targetDir, "README.md")
	readmeBytes, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("failed to read README.md: %v", err)
	}
	if string(readmeBytes) != "# real-repo-test\n" {
		t.Errorf("unexpected README.md content: %q", string(readmeBytes))
	}

	// Verify .git exists
	gitDir := filepath.Join(targetDir, ".git")
	if exists, _ := fs.Exists(gitDir); !exists {
		t.Errorf("expected .git to exist")
	}

	// Verify AGENTS.md exists
	agentsPath := filepath.Join(targetDir, "AGENTS.md")
	if exists, _ := fs.Exists(agentsPath); !exists {
		t.Errorf("expected AGENTS.md to exist")
	}
}

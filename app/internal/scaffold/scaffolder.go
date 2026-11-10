package scaffold

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
)

// Scaffolder handles orchestrating project template creation.
type Scaffolder struct {
	git GitClient
	fs  FS
}

// NewScaffolder creates a new Scaffolder instance.
func NewScaffolder(git GitClient, fs FS) *Scaffolder {
	return &Scaffolder{
		git: git,
		fs:  fs,
	}
}

// Run executes the full scaffolding process based on the provided options.
func (s *Scaffolder) Run(ctx context.Context, opts Options) (err error) {
	normalizedOpts, err := s.validateAndNormalizeOptions(opts)
	if err != nil {
		return err
	}

	targetExisted, err := s.fs.Exists(normalizedOpts.TargetDir)
	if err != nil {
		return fmt.Errorf("failed to check target directory %q: %w", normalizedOpts.TargetDir, err)
	}

	if targetExisted {
		empty, err := s.fs.IsEmptyDir(normalizedOpts.TargetDir)
		if err != nil {
			return fmt.Errorf("failed to inspect target directory %q: %w", normalizedOpts.TargetDir, err)
		}
		if !empty {
			return fmt.Errorf("%w: %q", ErrDirectoryNotEmpty, normalizedOpts.TargetDir)
		}
	}

	// If directory was newly created during clone and an error occurs later, roll back
	defer func() {
		if err != nil && !targetExisted {
			_ = s.fs.RemoveAll(normalizedOpts.TargetDir)
		}
	}()

	if err = s.git.Clone(ctx, normalizedOpts.RepoURL, normalizedOpts.TargetDir); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	if err = s.removeExcludedPaths(normalizedOpts.TargetDir); err != nil {
		return err
	}

	if err = s.git.Init(ctx, normalizedOpts.TargetDir); err != nil {
		return fmt.Errorf("failed to initialize git repository: %w", err)
	}

	if err = s.writeReadme(normalizedOpts.TargetDir, normalizedOpts.ProjectName); err != nil {
		return err
	}

	return nil
}

func (s *Scaffolder) validateAndNormalizeOptions(opts Options) (Options, error) {
	rawName := strings.TrimSpace(opts.ProjectName)
	if !isValidProjectName(rawName) {
		return Options{}, ErrInvalidProjectName
	}

	targetDir := strings.TrimSpace(opts.TargetDir)
	if targetDir == "" {
		targetDir = rawName
	}
	targetDir = filepath.Clean(targetDir)

	if !isValidTargetDir(targetDir) {
		return Options{}, ErrInvalidTargetDir
	}

	repoURL := strings.TrimSpace(opts.RepoURL)
	if repoURL == "" {
		repoURL = DefaultRepoURL
	}
	if !isValidRepoURL(repoURL) {
		return Options{}, ErrInvalidRepoURL
	}

	projectName := filepath.Base(rawName)
	if projectName == "." || projectName == string(filepath.Separator) {
		projectName = rawName
	}

	return Options{
		ProjectName: projectName,
		TargetDir:   targetDir,
		RepoURL:     repoURL,
	}, nil
}

func isValidProjectName(name string) bool {
	if name == "" || name == "." || name == ".." {
		return false
	}
	if strings.HasPrefix(name, "-") {
		return false
	}
	if strings.ContainsAny(name, "\r\n\x00") {
		return false
	}
	if strings.Contains(name, "..") {
		return false
	}
	return true
}

func isValidTargetDir(dir string) bool {
	if dir == "" || dir == "." || dir == ".." || dir == "/" || dir == `\` {
		return false
	}
	if strings.HasPrefix(dir, "-") {
		return false
	}
	if strings.HasPrefix(dir, "..") || strings.Contains(dir, filepath.FromSlash("/../")) {
		return false
	}
	return true
}

func isValidRepoURL(repo string) bool {
	if strings.HasPrefix(repo, "-") {
		return false
	}
	if strings.HasPrefix(repo, "ext::") {
		return false
	}
	return true
}

func (s *Scaffolder) removeExcludedPaths(targetDir string) error {
	appDir := filepath.Join(targetDir, "app")
	if err := s.fs.RemoveAll(appDir); err != nil {
		return fmt.Errorf("failed to remove app directory: %w", err)
	}

	quickstartPath := filepath.Join(targetDir, "quickstart.md")
	if err := s.fs.RemoveAll(quickstartPath); err != nil {
		return fmt.Errorf("failed to remove quickstart.md: %w", err)
	}

	gitDir := filepath.Join(targetDir, ".git")
	if err := s.fs.RemoveAll(gitDir); err != nil {
		return fmt.Errorf("failed to remove existing git metadata: %w", err)
	}

	return nil
}

func (s *Scaffolder) writeReadme(targetDir, projectName string) error {
	readmePath := filepath.Join(targetDir, "README.md")
	readmeContent := fmt.Sprintf("# %s\n", projectName)
	if err := s.fs.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		return fmt.Errorf("failed to write README.md: %w", err)
	}
	return nil
}

package scaffold

import "errors"

const (
	// DefaultRepoURL is the canonical source template repository.
	DefaultRepoURL = "https://github.com/willycornelissen/new-dev-ai"
)

var (
	// ErrInvalidProjectName indicates the provided project name is empty, contains path separators, or is invalid.
	ErrInvalidProjectName = errors.New("invalid project name")

	// ErrInvalidTargetDir indicates the target directory path is invalid, root, or unsafe.
	ErrInvalidTargetDir = errors.New("invalid target directory")

	// ErrInvalidRepoURL indicates the repository URL uses an unsafe or unsupported protocol.
	ErrInvalidRepoURL = errors.New("invalid repository URL")

	// ErrDirectoryNotEmpty indicates the destination directory exists and is not empty.
	ErrDirectoryNotEmpty = errors.New("target directory exists and is not empty")

	// ErrGitNotInstalled indicates git binary is not found in PATH.
	ErrGitNotInstalled = errors.New("git executable not found in PATH")
)

// Options holds configuration options for scaffolding a new project.
type Options struct {
	ProjectName string
	TargetDir   string
	RepoURL     string
}

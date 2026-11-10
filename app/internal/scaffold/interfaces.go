package scaffold

import (
	"context"
	"os"
)

// GitClient defines git operations required for scaffolding.
type GitClient interface {
	Clone(ctx context.Context, repoURL, targetDir string) error
	Init(ctx context.Context, targetDir string) error
}

// FS defines filesystem operations required for scaffolding.
type FS interface {
	Exists(path string) (bool, error)
	IsEmptyDir(path string) (bool, error)
	RemoveAll(path string) error
	WriteFile(path string, data []byte, perm os.FileMode) error
}

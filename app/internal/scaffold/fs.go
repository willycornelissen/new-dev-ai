package scaffold

import (
	"errors"
	"io"
	"os"
)

// OSFS implements FS using standard os package functions.
type OSFS struct{}

// NewOSFS creates a new OSFS instance.
func NewOSFS() *OSFS {
	return &OSFS{}
}

// Exists checks if a file, symlink, or directory exists at the given path.
func (f *OSFS) Exists(path string) (bool, error) {
	_, err := os.Lstat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

// IsEmptyDir checks if a directory exists and is empty.
func (f *OSFS) IsEmptyDir(path string) (bool, error) {
	dir, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer dir.Close()

	_, err = dir.Readdirnames(1)
	if errors.Is(err, io.EOF) {
		return true, nil
	}
	return false, err
}

// RemoveAll removes path and any children it contains.
func (f *OSFS) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

// WriteFile writes data to a file at path with the specified permissions.
func (f *OSFS) WriteFile(path string, data []byte, perm os.FileMode) error {
	return os.WriteFile(path, data, perm)
}

package scaffold

import (
	"context"
	"os"
)

type MockGitClient struct {
	CloneFunc func(ctx context.Context, repoURL, targetDir string) error
	InitFunc  func(ctx context.Context, targetDir string) error

	CloneCalls []struct {
		RepoURL   string
		TargetDir string
	}
	InitCalls []string
}

func (m *MockGitClient) Clone(ctx context.Context, repoURL, targetDir string) error {
	m.CloneCalls = append(m.CloneCalls, struct {
		RepoURL   string
		TargetDir string
	}{RepoURL: repoURL, TargetDir: targetDir})

	if m.CloneFunc != nil {
		return m.CloneFunc(ctx, repoURL, targetDir)
	}
	return nil
}

func (m *MockGitClient) Init(ctx context.Context, targetDir string) error {
	m.InitCalls = append(m.InitCalls, targetDir)
	if m.InitFunc != nil {
		return m.InitFunc(ctx, targetDir)
	}
	return nil
}

type MockFS struct {
	ExistsFunc     func(path string) (bool, error)
	IsEmptyDirFunc func(path string) (bool, error)
	RemoveAllFunc  func(path string) error
	WriteFileFunc  func(path string, data []byte, perm os.FileMode) error

	WrittenFiles map[string][]byte
	RemovedPaths []string
}

func NewMockFS() *MockFS {
	return &MockFS{
		WrittenFiles: make(map[string][]byte),
		RemovedPaths: make([]string, 0),
	}
}

func (m *MockFS) Exists(path string) (bool, error) {
	if m.ExistsFunc != nil {
		return m.ExistsFunc(path)
	}
	return false, nil
}

func (m *MockFS) IsEmptyDir(path string) (bool, error) {
	if m.IsEmptyDirFunc != nil {
		return m.IsEmptyDirFunc(path)
	}
	return true, nil
}

func (m *MockFS) RemoveAll(path string) error {
	m.RemovedPaths = append(m.RemovedPaths, path)
	if m.RemoveAllFunc != nil {
		return m.RemoveAllFunc(path)
	}
	return nil
}

func (m *MockFS) WriteFile(path string, data []byte, perm os.FileMode) error {
	m.WrittenFiles[path] = data
	if m.WriteFileFunc != nil {
		return m.WriteFileFunc(path, data, perm)
	}
	return nil
}

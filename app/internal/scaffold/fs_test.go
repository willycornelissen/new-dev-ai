package scaffold

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOSFS(t *testing.T) {
	fs := NewOSFS()
	tempDir := t.TempDir()

	t.Run("Exists on non-existent path", func(t *testing.T) {
		exists, err := fs.Exists(filepath.Join(tempDir, "does-not-exist"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if exists {
			t.Fatalf("expected exists to be false")
		}
	})

	t.Run("Exists on existing path", func(t *testing.T) {
		testFile := filepath.Join(tempDir, "test.txt")
		if err := os.WriteFile(testFile, []byte("hello"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		exists, err := fs.Exists(testFile)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !exists {
			t.Fatalf("expected exists to be true")
		}
	})

	t.Run("IsEmptyDir", func(t *testing.T) {
		emptyDir := filepath.Join(tempDir, "empty")
		if err := os.Mkdir(emptyDir, 0755); err != nil {
			t.Fatalf("failed to create dir: %v", err)
		}

		empty, err := fs.IsEmptyDir(emptyDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !empty {
			t.Fatalf("expected directory to be empty")
		}

		// Add a file
		fileInside := filepath.Join(emptyDir, "file.txt")
		if err := os.WriteFile(fileInside, []byte("data"), 0644); err != nil {
			t.Fatalf("failed to write file: %v", err)
		}

		empty, err = fs.IsEmptyDir(emptyDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if empty {
			t.Fatalf("expected directory to NOT be empty")
		}

		// Non-existent dir error
		_, err = fs.IsEmptyDir(filepath.Join(tempDir, "not-found"))
		if err == nil {
			t.Fatalf("expected error on non-existent dir")
		}
	})

	t.Run("WriteFile and RemoveAll", func(t *testing.T) {
		nestedDir := filepath.Join(tempDir, "nested")
		if err := os.Mkdir(nestedDir, 0755); err != nil {
			t.Fatalf("failed to create dir: %v", err)
		}

		targetFile := filepath.Join(nestedDir, "sample.md")
		err := fs.WriteFile(targetFile, []byte("# Title\n"), 0644)
		if err != nil {
			t.Fatalf("failed to write file: %v", err)
		}

		data, err := os.ReadFile(targetFile)
		if err != nil {
			t.Fatalf("failed to read written file: %v", err)
		}
		if string(data) != "# Title\n" {
			t.Errorf("expected '# Title\n', got %q", string(data))
		}

		// RemoveAll
		if err := fs.RemoveAll(nestedDir); err != nil {
			t.Fatalf("failed to remove dir: %v", err)
		}

		exists, err := fs.Exists(nestedDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if exists {
			t.Errorf("expected directory to be removed")
		}
	})
}

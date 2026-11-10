package scaffold

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestExecGitClient(t *testing.T) {
	ctx := context.Background()

	t.Run("NewExecGitClient succeeds when git exists", func(t *testing.T) {
		client, err := NewExecGitClient()
		if err != nil {
			t.Fatalf("expected git to be found in test environment: %v", err)
		}
		if client == nil {
			t.Fatalf("expected client to be non-nil")
		}
	})

	t.Run("Zero value ExecGitClient uses default runner safely", func(t *testing.T) {
		client := &ExecGitClient{}
		runner := client.getRunner()
		if runner == nil {
			t.Fatalf("expected runner to default to DefaultCommandRunner")
		}
	})

	t.Run("Clone invokes runner with '--' delimiter to prevent argument injection", func(t *testing.T) {
		var invokedName string
		var invokedArgs []string

		mockRunner := func(ctx context.Context, name string, args ...string) ([]byte, error) {
			invokedName = name
			invokedArgs = args
			return []byte("cloned"), nil
		}

		client := NewExecGitClientWithRunner(mockRunner)
		err := client.Clone(ctx, "https://github.com/example/repo", "/tmp/dest")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if invokedName != "git" {
			t.Errorf("expected command 'git', got: %s", invokedName)
		}
		expectedArgs := []string{"clone", "--depth", "1", "--", "https://github.com/example/repo", "/tmp/dest"}
		if len(invokedArgs) != len(expectedArgs) {
			t.Fatalf("expected %d args, got %d", len(expectedArgs), len(invokedArgs))
		}
		for i, arg := range expectedArgs {
			if invokedArgs[i] != arg {
				t.Errorf("expected arg %d to be %s, got: %s", i, arg, invokedArgs[i])
			}
		}
	})

	t.Run("Clone propagates runner error", func(t *testing.T) {
		mockRunner := func(ctx context.Context, name string, args ...string) ([]byte, error) {
			return nil, errors.New("command failed")
		}

		client := NewExecGitClientWithRunner(mockRunner)
		err := client.Clone(ctx, "https://github.com/example/repo", "/tmp/dest")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	t.Run("Init invokes runner with correct arguments", func(t *testing.T) {
		var invokedName string
		var invokedArgs []string

		mockRunner := func(ctx context.Context, name string, args ...string) ([]byte, error) {
			invokedName = name
			invokedArgs = args
			return []byte("Initialized"), nil
		}

		client := NewExecGitClientWithRunner(mockRunner)
		err := client.Init(ctx, "/tmp/dest")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if invokedName != "git" {
			t.Errorf("expected command 'git', got: %s", invokedName)
		}
		expectedArgs := []string{"-C", "/tmp/dest", "init"}
		if len(invokedArgs) != len(expectedArgs) {
			t.Fatalf("expected %d args, got %d", len(expectedArgs), len(invokedArgs))
		}
		for i, arg := range expectedArgs {
			if invokedArgs[i] != arg {
				t.Errorf("expected arg %d to be %s, got: %s", i, arg, invokedArgs[i])
			}
		}
	})

	t.Run("Init propagates runner error", func(t *testing.T) {
		mockRunner := func(ctx context.Context, name string, args ...string) ([]byte, error) {
			return nil, errors.New("git init failed")
		}

		client := NewExecGitClientWithRunner(mockRunner)
		err := client.Init(ctx, "/tmp/dest")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	t.Run("DefaultCommandRunner executes real command and scrubs credentials on error", func(t *testing.T) {
		out, err := DefaultCommandRunner(ctx, "echo", "hello")
		if err != nil {
			t.Fatalf("unexpected error running echo: %v", err)
		}
		if string(out) != "hello\n" {
			t.Errorf("expected 'hello\\n', got %q", string(out))
		}

		// Test credential scrubbing when command returns error output containing token
		_, err = DefaultCommandRunner(ctx, "sh", "-c", "echo 'fatal: https://user:secret123@github.com/repo failed' >&2; exit 1")
		if err == nil {
			t.Fatalf("expected error from sh command")
		}
		if strings.Contains(err.Error(), "secret123") {
			t.Errorf("expected secret to be scrubbed, but got: %v", err)
		}
		if !strings.Contains(err.Error(), "://***:***@") {
			t.Errorf("expected masked credential placeholder in error, got: %v", err)
		}
	})
}

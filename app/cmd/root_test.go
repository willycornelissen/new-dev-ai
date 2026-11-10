package cmd

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"new-dev-ai/internal/scaffold"
)

type mockScaffolderRunner struct {
	runFunc func(ctx context.Context, opts scaffold.Options) error
	calls   []scaffold.Options
}

func (m *mockScaffolderRunner) Run(ctx context.Context, opts scaffold.Options) error {
	m.calls = append(m.calls, opts)
	if m.runFunc != nil {
		return m.runFunc(ctx, opts)
	}
	return nil
}

func TestRootCmd(t *testing.T) {
	t.Run("fails when no arguments are provided", func(t *testing.T) {
		runner := &mockScaffolderRunner{}
		cmd := NewRootCmd(runner)

		var outBuf, errBuf bytes.Buffer
		cmd.SetOut(&outBuf)
		cmd.SetErr(&errBuf)
		cmd.SetArgs([]string{})

		err := cmd.Execute()
		if err == nil {
			t.Fatalf("expected error when no args provided")
		}
		if len(runner.calls) != 0 {
			t.Errorf("scaffolder should not have been called")
		}
	})

	t.Run("succeeds with project name and default flags", func(t *testing.T) {
		runner := &mockScaffolderRunner{}
		cmd := NewRootCmd(runner)

		var outBuf, errBuf bytes.Buffer
		cmd.SetOut(&outBuf)
		cmd.SetErr(&errBuf)
		cmd.SetArgs([]string{"hello-world"})

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(runner.calls) != 1 {
			t.Fatalf("expected 1 call, got: %d", len(runner.calls))
		}
		call := runner.calls[0]
		if call.ProjectName != "hello-world" {
			t.Errorf("expected project name 'hello-world', got: %s", call.ProjectName)
		}
		if call.RepoURL != scaffold.DefaultRepoURL {
			t.Errorf("expected default repo URL %s, got: %s", scaffold.DefaultRepoURL, call.RepoURL)
		}

		output := outBuf.String()
		if !strings.Contains(output, "Project 'hello-world' created successfully!") {
			t.Errorf("expected success message in stdout, got: %s", output)
		}
	})

	t.Run("supports custom --repo and --dir flags", func(t *testing.T) {
		runner := &mockScaffolderRunner{}
		cmd := NewRootCmd(runner)

		var outBuf, errBuf bytes.Buffer
		cmd.SetOut(&outBuf)
		cmd.SetErr(&errBuf)
		cmd.SetArgs([]string{
			"custom-app",
			"--repo", "https://github.com/custom/template",
			"--dir", "/custom/target/dir",
		})

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(runner.calls) != 1 {
			t.Fatalf("expected 1 call, got: %d", len(runner.calls))
		}
		call := runner.calls[0]
		if call.ProjectName != "custom-app" {
			t.Errorf("expected project name 'custom-app', got: %s", call.ProjectName)
		}
		if call.RepoURL != "https://github.com/custom/template" {
			t.Errorf("expected custom repo, got: %s", call.RepoURL)
		}
		if call.TargetDir != "/custom/target/dir" {
			t.Errorf("expected custom dir, got: %s", call.TargetDir)
		}
	})

	t.Run("returns error when scaffolder fails", func(t *testing.T) {
		runner := &mockScaffolderRunner{
			runFunc: func(ctx context.Context, opts scaffold.Options) error {
				return errors.New("simulated scaffold failure")
			},
		}
		cmd := NewRootCmd(runner)

		var outBuf, errBuf bytes.Buffer
		cmd.SetOut(&outBuf)
		cmd.SetErr(&errBuf)
		cmd.SetArgs([]string{"failed-app"})

		err := cmd.Execute()
		if err == nil {
			t.Fatalf("expected error from scaffolder")
		}
		if !strings.Contains(err.Error(), "simulated scaffold failure") {
			t.Errorf("expected error message to contain failure details, got: %v", err)
		}
	})
}

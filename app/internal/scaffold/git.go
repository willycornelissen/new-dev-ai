package scaffold

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var credentialRegex = regexp.MustCompile(`://([^:]+):([^@]+)@`)

// CommandRunner abstracts execution of external commands.
type CommandRunner func(ctx context.Context, name string, args ...string) ([]byte, error)

// DefaultCommandRunner runs a command using os/exec with disabled terminal prompts.
func DefaultCommandRunner(ctx context.Context, name string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	output, err := cmd.Output()
	if err != nil {
		errOutput := strings.TrimSpace(stderr.String())
		if errOutput != "" {
			sanitizedErr := credentialRegex.ReplaceAllString(errOutput, "://***:***@")
			return nil, fmt.Errorf("%w: %s", err, sanitizedErr)
		}
		return nil, err
	}
	return output, nil
}

// ExecGitClient implements GitClient by executing the git CLI.
type ExecGitClient struct {
	runner CommandRunner
}

// NewExecGitClient creates a new ExecGitClient.
func NewExecGitClient() (*ExecGitClient, error) {
	if _, err := exec.LookPath("git"); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGitNotInstalled, err)
	}
	return &ExecGitClient{
		runner: DefaultCommandRunner,
	}, nil
}

// NewExecGitClientWithRunner creates an ExecGitClient with a custom CommandRunner.
func NewExecGitClientWithRunner(runner CommandRunner) *ExecGitClient {
	return &ExecGitClient{
		runner: runner,
	}
}

func (g *ExecGitClient) getRunner() CommandRunner {
	if g == nil || g.runner == nil {
		return DefaultCommandRunner
	}
	return g.runner
}

// Clone clones a git repository into the destination directory safely using '--'.
func (g *ExecGitClient) Clone(ctx context.Context, repoURL, targetDir string) error {
	_, err := g.getRunner()(ctx, "git", "clone", "--depth", "1", "--", repoURL, targetDir)
	if err != nil {
		return fmt.Errorf("git clone execution failed: %w", err)
	}
	return nil
}

// Init initializes a fresh git repository in the destination directory.
func (g *ExecGitClient) Init(ctx context.Context, targetDir string) error {
	_, err := g.getRunner()(ctx, "git", "-C", targetDir, "init")
	if err != nil {
		return fmt.Errorf("git init execution failed: %w", err)
	}
	return nil
}

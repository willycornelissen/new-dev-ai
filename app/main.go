package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"new-dev-ai/cmd"
	"new-dev-ai/internal/scaffold"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	gitClient, err := scaffold.NewExecGitClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fs := scaffold.NewOSFS()
	scaffolder := scaffold.NewScaffolder(gitClient, fs)
	rootCmd := cmd.NewRootCmd(scaffolder)

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}

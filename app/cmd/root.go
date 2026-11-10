package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"new-dev-ai/internal/scaffold"
)

// ScaffolderRunner defines the interface for running the scaffolding process.
type ScaffolderRunner interface {
	Run(ctx context.Context, opts scaffold.Options) error
}

// NewRootCmd constructs and returns the root cobra command for new-dev-ai.
func NewRootCmd(runner ScaffolderRunner) *cobra.Command {
	var (
		repoURL   string
		targetDir string
	)

	cmd := &cobra.Command{
		Use:   "new-dev-ai <nome-do-projeto>",
		Short: "Scaffold a new development project tailored for AI workflows",
		Long: `new-dev-ai is a CLI tool that creates a new project directory cloned from
https://github.com/willycornelissen/new-dev-ai without the app/ directory and quickstart.md,
with a fresh git repository initialized and a clean README.md.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			projectName := args[0]

			cmd.Printf("🚀 Scaffolding project '%s'...\n", projectName)

			opts := scaffold.Options{
				ProjectName: projectName,
				TargetDir:   targetDir,
				RepoURL:     repoURL,
			}

			if err := runner.Run(cmd.Context(), opts); err != nil {
				return err
			}

			cmd.Println("📦 Cloned template repository.")
			cmd.Println("🧹 Excluded app/ directory and quickstart.md.")
			cmd.Println("🔧 Initialized fresh Git repository.")
			cmd.Println("📝 Initialized README.md.")
			cmd.Printf("✨ Project '%s' created successfully!\n", projectName)
			return nil
		},
	}

	cmd.Flags().StringVar(&repoURL, "repo", scaffold.DefaultRepoURL, "URL of the template repository to clone")
	cmd.Flags().StringVar(&targetDir, "dir", "", "Target directory path (defaults to project name)")

	return cmd
}

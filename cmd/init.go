package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(initCmd)
	initFlags()
}

var defaultBranch string

var initCmd = &cobra.Command{
	Use:   "init [directory]",
	Short: "Create an empty Git repository",
	Long: `Description:
  Init command creates an empty Git repository to the current working directory.

  The init command also creates a .git directory with subdirectories for objects,
  refs/heads, refs/tags and template files.

  By default Gong initializes the repository's default branch as main instead of master.`,
	Run: run,
}

func initFlags() {
	initCmd.Flags().StringVarP(
		&defaultBranch, "default-branch", "d", "",
		"Use specified name for the default branch, when creating a new repository.",
	)
}

func run(cmd *cobra.Command, args []string) {
}

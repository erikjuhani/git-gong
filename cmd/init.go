package cmd

import (
	"os"

	"github.com/erikjuhani/git-gong/gong"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
	initFlags()
}

var (
	defaultBranch string
	bare          bool
)

var initCmd = &cobra.Command{
	Use:   "init [directory]",
	Short: "Create an empty Git repository",
	Long: `Description:
  Init command creates an empty Git repository to the current working directory.

  The init command also creates a .git directory with subdirectories for objects,
  refs/heads, refs/tags and template files.

  By default Gong initializes the repository's default branch as main instead of master.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path, err := os.Getwd()
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		if len(args) == 1 {
			path = args[0]
		}

		err = initRepository(path)
		if err != nil {
			cmd.PrintErr(err)
			return
		}
	},
}

func initFlags() {
	initCmd.Flags().StringVarP(
		&defaultBranch, "default-branch", "d", gong.DefaultReference,
		"Use specified name for the default branch, when creating a new repository.",
	)
}

func initRepository(path string) error {
	repo, err := gong.Init(path, bare, defaultBranch)
	defer gong.Free(repo)

	return err
}

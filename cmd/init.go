package cmd

import (
	"fmt"
	"os"

	git "github.com/libgit2/git2go/v30"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
	initFlags()
}

var (
	defaultBranch string = "main"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := os.Getwd()
		if err != nil {
			return err
		}

		if len(args) == 1 {
			path = args[0]
		}

		return initRepository(path)
	},
}

func initFlags() {
	initCmd.Flags().StringVarP(
		&defaultBranch, "default-branch", "d", defaultBranch,
		"Use specified name for the default branch, when creating a new repository.",
	)
}

func initRepository(path string) error {
	repo, err := git.InitRepository(path, bare)
	if err != nil {
		return err
	}

	idx, err := repo.Index()
	if err != nil {
		return err
	}

	treeID, err := idx.WriteTree()
	if err != nil {
		return err
	}

	initRef := fmt.Sprintf("refs/heads/%s", defaultBranch)

	ref, err := repo.References.Create(initRef, treeID, false, "Repository initialized")
	if err != nil {
		return err
	}
	defer ref.Free()

	return repo.SetHead(initRef)
}

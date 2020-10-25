package cmd

import (
	"fmt"
	"gong/gong"

	git "github.com/libgit2/git2go/v30"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(commitCmd)
	commitFlags()
}

var (
	repo      *git.Repository
	stageOnly bool
)

var commitCmd = &cobra.Command{
	Use:   "commit [pathspec]",
	Short: "Record changes to index and repository",
	Long: `Create a new commit containing the contents of the index.
  
  To only stage file changes apply a flag --stage. The files won't be recorded
  until the next call for commit.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		err := commit(args)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func commitFlags() {
	commitCmd.Flags().BoolVarP(
		&stageOnly, "stage", "s", false,
		"Use to stage file changes instead of a commit",
	)

}

func commit(paths []string) (err error) {
	repo, err := gong.Open()
	if err != nil {
		return
	}

	treeID, err := repo.AddToIndex(paths)
	if err != nil {
		return
	}

	_, err = repo.Commit(treeID)
	if err != nil {
		return
	}

	return
}

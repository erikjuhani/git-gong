package cmd

import (
	"github.com/erikjuhani/git-gong/gong"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(commitCmd)
	commitFlags()
}

var (
	stageOnly bool
	commitMsg string
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
			cmd.PrintErr(err)
			return
		}
	},
}

func commitFlags() {
	commitCmd.Flags().BoolVarP(
		&stageOnly, "stage", "s", false,
		"Use to stage file changes instead of a commit",
	)
	commitCmd.Flags().StringVarP(
		&commitMsg, "message", "m", "",
		"Set commit message",
	)
}

func commit(paths []string) (err error) {
	repo, err := gong.Open()
	if err != nil {
		return
	}

	defer repo.Free()

	treeID, err := repo.AddToIndex(paths)
	if err != nil {
		return
	}

	if stageOnly {
		return
	}

	_, err = repo.Commit(treeID, commitMsg)

	return
}

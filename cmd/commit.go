package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(commitCmd)
}

var commitCmd = &cobra.Command{
	Use:   "commit [pathspec]",
	Short: "Record changes to index and repository",
	Long: `Create a new commit containing the contents of the index.
  
  To only stage file changes apply a flag --stage. The files won't be recorded
  until the next call for commit.
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return commit()
	},
}

func commitFlags() {
	commitCmd.Flags().BoolP(
		"stage", "s", false,
		"Use to stage file changes instead of a commit",
	)

}

func commit() error {
	return nil
}

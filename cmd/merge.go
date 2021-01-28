package cmd

import (
	"github.com/erikjuhani/git-gong/gong"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(mergeCmd)
}

var mergeCmd = &cobra.Command{
	Use:   "merge [branchname]",
	Short: "Merges the given branch to current branch",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := gong.Open()
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		defer gong.Free(repo)

		if err := repo.Merge(args[0]); err != nil {
			cmd.PrintErr(err)
		}
	},
}

package cmd

import (
	"github.com/erikjuhani/git-gong/gong"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(infoCmd)
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show information of the current HEAD.",
	Long: `Display difference between the index file and the current HEAD in short format.
  Example status info output
  Branch <branchname>
  Commit <commitHash>

  Changes (<number of file changes):
  A <filename>
  M <filename>
  D <filename>`,
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := gong.Open()
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		defer gong.Free(repo)

		info, err := repo.Info()
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		cmd.Println(info)
	},
}

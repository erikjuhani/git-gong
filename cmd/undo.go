package cmd

import (
	"github.com/erikjuhani/git-gong/gong"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(undoCmd)
}

var undoCmd = &cobra.Command{
	Use:   "undo",
	Short: "Undo undoes the last command of the user.",
	Long: `If user for example has made a mistake commit gong commit -m "mistake"
		the undo, undoes the commit command and sets the repository to a prior state.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := gong.Open()
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		defer gong.Free(repo)
	},
}

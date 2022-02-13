package cmd

import (
	"github.com/erikjuhani/git-gong/gong"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(gitCmd)
}

var gitCmd = &cobra.Command{
	Use:   "git [args]",
	Short: "Run Git commands by providing any Git compatible arguments",
	Long: `Run Git commands by providing any Git compatible arguments

  Needs Git installed in the system.
	`,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		if err := gong.RunGitCommand(args); err != nil {
			cmd.PrintErr(err)
		}
	},
}

package cmd

import (
	"github.com/erikjuhani/git-gong/gong"
	"github.com/spf13/cobra"
)

var (
	rename        bool
	changeMessage string
)

func init() {
	rootCmd.AddCommand(changeCmd)
	changeCmd.AddCommand(changeBranchCmd, changeCommitCmd, changeFileCmd)
	initChangeFlags()
}

func initChangeFlags() {
	changeBranchCmd.Flags().BoolVarP(
		&rename, "rename", "r", false,
		"Change branch's name",
	)

	changeFileCmd.Flags().BoolVarP(
		&rename, "rename", "r", false,
		"Change file's name",
	)

	changeCommitCmd.Flags().StringVarP(
		&changeMessage, "message", "m", "",
		"Change last commit message",
	)
}

var changeCmd = &cobra.Command{
	Use:   "change",
	Short: "",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
}

var changeBranchCmd = &cobra.Command{
	Use:   "branch [branchname]",
	Short: "Change branch attributes",
	Long:  ``,
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := gong.Open()
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		defer gong.Free(repo)
	},
}

var changeCommitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Change commit attributes",
	Long:  ``,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := gong.Open()
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		defer gong.Free(repo)
	},
}

var changeFileCmd = &cobra.Command{
	Use:   "file",
	Short: "Change file attributes",
	Long:  ``,
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := gong.Open()
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		defer gong.Free(repo)
	},
}

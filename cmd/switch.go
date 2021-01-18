package cmd

import (
	"github.com/erikjuhani/git-gong/gong"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(switchCmd)

	switchCmd.AddCommand(
		switchBranchCmd,
		switchCommitCmd,
		switchTagCmd,
		switchReleaseCmd,
	)
}

var switchCmd = &cobra.Command{
	Use:   "switch [subcommand]",
	Short: "Switch to branches, commits, tags or releases.",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var switchBranchCmd = &cobra.Command{
	Use:   "branch [branchname]",
	Short: "Switch to branch with branchname.",
	Long: `If branchname does not exist create branch with branchname.
		if there are any unsaved changes stash them to @<previousbranchname>.
		When switching to branch check if there exists a stash and pop the stash.
	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := gong.Open()
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		defer gong.Free(repo)

		branch, err := repo.CheckoutBranch(args[0])
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		defer gong.Free(branch)

		cmd.Printf("checkout to branch %s\n", args[0])
	},
}

var switchCommitCmd = &cobra.Command{
	Use:   "commit [commithash]",
	Short: "Switch to commit.",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := gong.Open()
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		defer gong.Free(repo)

		commit, err := repo.CheckoutCommit(args[0])
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		defer gong.Free(commit)

		cmd.Printf("switched to commit %s\n", args[0])
	},
}

var switchTagCmd = &cobra.Command{
	Use:   "tag [tag]",
	Short: "Switch to tag.",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := gong.Open()
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		defer gong.Free(repo)

		tag, err := repo.CheckoutTag(args[0])
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		defer gong.Free(tag)

		cmd.Printf("checkout to tag %s\n", tag.Name)
	},
}

var switchReleaseCmd = &cobra.Command{
	Use:   "release [release]",
	Short: "Switch to release.",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := gong.Open()
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		defer gong.Free(repo)

		tag, err := repo.CheckoutTag(args[0])
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		defer gong.Free(tag)

		cmd.Printf("checkout to release %s\n", tag.Name)
	},
}

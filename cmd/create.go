package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(
		createCmd,
		createBranchCmd,
		createFileCmd,
		createReleaseCmd,
		createTagCmd,
	)
}

// TODO: write long descriptions
var createCmd = &cobra.Command{
	Use:   "create [subcommand]",
	Short: "Create branches, files, releases and tags.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var createBranchCmd = &cobra.Command{
	Use:   "branch [branchname]",
	Short: "Creates a new branch.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var createFileCmd = &cobra.Command{
	Use:   "file [filename]",
	Short: "Creates a regular file.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var createDirectoryCmd = &cobra.Command{
	Use:   "directory [dirname]",
	Short: "Creates a directory.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var createReleaseCmd = &cobra.Command{
	Use:   "release [releasename]",
	Short: "Creates a release / tag.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var createTagCmd = &cobra.Command{
	Use:   "tag [tagname]",
	Short: "Creates a tag / release.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

package cmd

import (
	"os"
	"path/filepath"

	"github.com/erikjuhani/git-gong/gong"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.AddCommand(
		createBranchCmd,
		createFileCmd,
		createDirectoryCmd,
		createReleaseCmd,
		createTagCmd,
	)
}

// TODO: write long descriptions
var createCmd = &cobra.Command{
	Use:   "create [subcommand]",
	Short: "Create branches, files, releases and tags.",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
}

var createBranchCmd = &cobra.Command{
	Use:   "branch [branchname]",
	Short: "Creates a new branch.",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := gong.Open()
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		defer repo.Free()

		branch, err := repo.CreateLocalBranch(args[0])
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		name, err := branch.Name()
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		cmd.Printf("created a new branch %s\n", name)
	},
}

var createFileCmd = &cobra.Command{
	Use:   "file [filename]",
	Short: "Creates a regular file",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]
		if err := os.MkdirAll(filepath.Dir(file), os.ModePerm); err != nil {
			cmd.PrintErr(err)
			return
		}

		_, err := os.Create(file)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		cmd.Printf("created a new file %s\n", file)
	},
}

var createDirectoryCmd = &cobra.Command{
	Use:   "directory [dirname]",
	Short: "Creates a directory",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		directory := args[0]
		if err := os.MkdirAll(directory, os.ModePerm); err != nil {
			cmd.PrintErr(err)
			return
		}

		cmd.Printf("created a new directory %s\n", directory)
	},
}

var createReleaseCmd = &cobra.Command{
	Use:   "release [releasename]",
	Short: "Creates a release / tag",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := gong.Open()
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		defer repo.Free()

		message := ""

		if len(args) > 1 {
			message = args[1]
		}

		tagID, err := repo.CreateTag(args[0], message)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		tag, err := repo.Core.LookupTag(tagID)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		cmd.Printf("created a new release %s\n", tag.Name())
	},
}

var createTagCmd = &cobra.Command{
	Use:   "tag [tagname]",
	Short: "Creates a tag / release",
	Long:  ``,
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := gong.Open()
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		defer repo.Free()

		message := ""

		if len(args) > 1 {
			message = args[1]
		}

		tagID, err := repo.CreateTag(args[0], message)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		tag, err := repo.Core.LookupTag(tagID)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		cmd.Printf("created a new tag %s\n", tag.Name())
	},
}

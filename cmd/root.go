package cmd

import (
	"fmt"
	"os"

	git "github.com/libgit2/git2go/v30"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "gong",
	Run: func(cmd *cobra.Command, args []string) {
		path, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		repo, err := git.OpenRepository(path)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		statusList, err := repo.StatusList(&git.StatusOptions{
			Show:  git.StatusShowIndexAndWorkdir,
			Flags: git.StatusOptIncludeUntracked,
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		defer statusList.Free()

		count, err := statusList.EntryCount()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for i := 0; i < count; i++ {
			entry, err := statusList.ByIndex(i)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
				return
			}
			if entry.Status <= 0 {
				continue
			}
			fmt.Printf("modified: %s\n", entry.IndexToWorkdir.NewFile.Path)
		}

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

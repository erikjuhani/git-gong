package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/erikjuhani/git-gong/gong"
	git "github.com/libgit2/git2go/v31"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cloneCmd)
}

var cloneCmd = &cobra.Command{
	Use:   "clone [repository] [directory]",
	Short: "Clone a repository",
	Long: `Description:
  Clone command clones a repository into a newly created directory and checks out
  to initial branch of the cloned repository.

  Clone command takes a valid GIT URL, where the repository is located as an argument.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := os.Getwd()
		if err != nil {
			cmd.PrintErr(err)
			return err
		}

		if len(args) > 1 {
			path = fmt.Sprintf("%s/%s", path, args[1])
		} else {
			parts := strings.Split(args[0], "/")
			path = fmt.Sprintf("%s/%s", parts[len(parts)-3], strings.TrimSuffix(parts[len(parts)-1], ".git"))
		}

		return cloneRepository(args[0], path)
	},
}

// TODO: Create a clone function to gong package.
func cloneRepository(url string, path string) error {
	opts := git.CloneOptions{}
	repo, err := git.Clone(url, path, &opts)
	defer gong.Free(repo)

	return err
}

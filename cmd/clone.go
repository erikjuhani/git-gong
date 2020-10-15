package cmd

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	git "github.com/libgit2/git2go/v30"
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
			return err
		}

		gitURL, err := url.Parse(args[0])
		if err != nil {
			return err
		}

		if len(args) > 1 {
			path = fmt.Sprintf("%s/%s", path, args[1])
		} else {
			parts := strings.Split(gitURL.String(), "/")
			path = fmt.Sprintf("%s/%s", path, strings.TrimSuffix(parts[len(parts)-1], ".git"))
		}

		return cloneRepository(gitURL.String(), path)
	},
}

func cloneRepository(url string, path string) error {
	opts := git.CloneOptions{}
	_, err := git.Clone(url, path, &opts)
	return err
}

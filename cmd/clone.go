package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cloneCmd)
}

var cloneCmd = &cobra.Command{
	Use:   "clone [remote]",
	Short: "Clone a repository",
	Long: `Description:
  Clone command clones a repository into a newly created directory and checks out
  to initial branch of the cloned repository.

  Clone command takes a valid GIT URL, where the repository is located as an argument.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cloneRepository(args[0])
	},
}

func cloneRepository(url string) error {
	return nil
}

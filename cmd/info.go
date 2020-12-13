package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(infoCmd)
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show information of the current HEAD.",
	Long:  `Display difference between the index file and the current HEAD in short format.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
	},
}

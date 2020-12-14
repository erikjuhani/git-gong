package cmd

import (
	"fmt"
	"os"

	"github.com/erikjuhani/git-gong/doc"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "gong",
	Short:   "",
	Long:    ``,
	Version: doc.Version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

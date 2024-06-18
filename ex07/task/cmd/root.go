package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "task",
		Short: "task is a CLI for managing your TODOs.",
	}
)

func Execute() error {
	return rootCmd.Execute()
}

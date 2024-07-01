package cmd

import (
	"github.com/spf13/cobra"

	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

var encodingKey string

var rootCmd = &cobra.Command{
	Use:   "secret",
	Short: "Secret is and API Key and other secrets manager.",
}

func Execute() error {
	rootCmd.PersistentFlags().StringVarP(&encodingKey, "key", "k", "", "The key to use when encoding and decoding secrets")

	return rootCmd.Execute()
}

func secretsPath() string {
	home, _ := homedir.Dir()
	return filepath.Join(home, ".secrets")
}

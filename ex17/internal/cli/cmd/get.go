package cmd

import (
	"fmt"
	"main/internal/secret"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets a secret in your secret vault.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Insufficient number of arguments")
			return
		}

		v := secret.NewFileVault(encodingKey, secretsPath())
		value, err := v.Get(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(value)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}

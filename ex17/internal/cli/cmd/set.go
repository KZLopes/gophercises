package cmd

import (
	"fmt"
	"main/internal/secret"

	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Sets a secret in your secret vault.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("Insufficient number of arguments")
			return
		}
		v := secret.NewFileVault(encodingKey, secretsPath())
		err := v.Set(args[0], args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Secret set for %s\n", args[0])
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}

package cmd

import (
	"fmt"
	"os"
	"strings"
	"taskm/db"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: " Add a new task to your TODO list",
	Run: func(cmd *cobra.Command, args []string) {
		task := strings.Join(args, " ")
		_, err := db.CreateTask(task)
		if err != nil {
			fmt.Println("Something went wrong:", err.Error())
			os.Exit(1)
		}
		fmt.Printf("Added \"%s\" to your task list.\n", task)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

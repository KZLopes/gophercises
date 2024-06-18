package cmd

import (
	"fmt"
	"taskm/db"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all of your incomplete tasks",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := db.AllTasks()
		if err != nil {
			fmt.Println(err.Error())
		}

		if len(tasks) == 0 {
			fmt.Println("You have no more tasks!")
			return
		}
		fmt.Println("You have the following tasks:")
		for i, task := range tasks {
			fmt.Printf("%d. %s\n", i+1, task.Value)
		}

	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

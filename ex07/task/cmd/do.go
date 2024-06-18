package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"taskm/db"

	"github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Mark a task on your TODO list as complete",
	Run: func(cmd *cobra.Command, args []string) {
		var ids []int
		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Println("Not a Valid Integer:", arg)
				continue
			}
			ids = append(ids, id)
		}

		tasks, err := db.AllTasks()
		if err != nil {
			fmt.Println("Something Went Wrong", err)
			os.Exit(1)
		}

		for _, i := range ids {
			if i < 0 || i > len(tasks) {
				fmt.Println("invalid Task Number:", i)
				continue
			}

			t := tasks[i-1]
			err := db.DeleteTask(t.Key)
			if err != nil {
				fmt.Printf("Failed to Complete Task \"%d\"\nError: %s", t.Key, err)
				continue
			}
			fmt.Printf("You have completed the \"%s\" task.", t.Value)
		}
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}

func doTask(slice []db.Task, idx int) ([]db.Task, error) {
	if idx < 0 || idx >= len(slice) {
		return slice, errors.New("Invalid Position")
	}
	return append(slice[:idx], slice[idx+1:]...), nil
}

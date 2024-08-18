package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"taskman/task"
	"time"
)

var newCmd = &cobra.Command{
	Use:   "new \"task name\" \"due date (YYYY-MM-DD H:MM AM/PM)\"",
	Short: "Create a new task",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		date, err := time.Parse("2006-01-02 11:59 PM", args[1])
		if err != nil {
			fmt.Println("Invalid date format. Please use YYYY-MM-DD H:MM AM/PM.")
			return
		}
		task.AddTask(DB, name, date)
		fmt.Printf("Task (%s) added\n", name)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}

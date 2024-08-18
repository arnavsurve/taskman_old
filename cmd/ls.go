package cmd

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
	"taskman/task"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all tasks",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := task.ListTasks(DB)
		if err != nil {
			log.Fatal("Error querying tasks:", err)
			return
		}

		if len(tasks) < 1 {
			fmt.Println("No tasks found.")
			return
		}

		// Create a new tablewriter table
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Due", "Status"})

		for _, task := range tasks {
			// Convert each task to a slice of strings for the table
			row := []string{
				strconv.Itoa(task.ID),
				task.Name,
				task.DueDate.Format("2006-01-02 11:59 PM"),
				string(task.Status),
			}
			table.Append(row)
		}

		// Set table options for better formatting
		table.SetBorder(true)
		table.SetRowLine(true)
		table.Render() // Render the table to stdout
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
	"taskman/task"
	"time"

	"github.com/manifoldco/promptui"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Interactively edit a task's details",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := task.ListTasks(DB)
		if err != nil {
			fmt.Println("Error fetching tasks:", err)
			return
		}

		if len(tasks) == 0 {
			fmt.Println("No tasks available to edit.")
			return
		}

		// Select a task to edit
		taskID, err := selectTask(tasks)
		if err != nil {
			fmt.Println("Selection canceled.")
			return
		}

		// Find the task by ID
		var selectedTask task.Task
		for _, t := range tasks {
			if t.ID == taskID {
				selectedTask = t
				break
			}
		}

		// Edit the task
		newName, err := promptForInput("Name", selectedTask.Name)
		if err != nil {
			fmt.Println("Edit canceled.")
			return
		}

		newDueDateStr, err := promptForInput("Due Date (YYYY-MM-DD H:MM AM/PM)", selectedTask.DueDate.Format("2006-01-02 11:59 PM"))
		if err != nil {
			fmt.Println("Edit canceled.")
			return
		}
		newDueDate, err := time.Parse("2006-01-02 11:59 PM", newDueDateStr)
		if err != nil {
			fmt.Println("Invalid date format.")
			return
		}

		newStatusStr, err := promptForSelect("Status", []string{"TODO", "IN_PROGRESS", "DONE"})
		if err != nil {
			fmt.Println("Edit canceled.")
			return
		}

		newStatus := task.TaskStatus(newStatusStr)

		task.UpdateTask(DB, taskID, newName, newDueDate, newStatus)
		fmt.Println("Task updated successfully.")
	},
}

// selectTask prompts the user to select a task from a list of tasks
func selectTask(tasks []task.Task) (int, error) {
	taskOptions := make([]string, len(tasks))
	for i, task := range tasks {
		taskOptions[i] = fmt.Sprintf("%d: %s [%s] Due: %s", task.ID, task.Name, task.Status, task.DueDate.Format("2006-01-02 11:59 PM"))
	}

	prompt := promptui.Select{
		Label: "Select Task",
		Items: taskOptions,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return 0, err
	}

	// Extract the task ID from the selected string because this library fucking sucks
	// and forces me to parse a string slice instead of the actual task object (????)
	idStr := strings.Split(result, ":")[0]
	taskID, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("Failed to parse task ID: %v", err)
	}

	return taskID, err
}

// promptForInput prompts the user for input with a default value
func promptForInput(label, defaultValue string) (string, error) {
	prompt := promptui.Prompt{
		Label:   label,
		Default: defaultValue,
	}
	return prompt.Run()
}

// promptForSelect prompts the user to select an option from a list
func promptForSelect(label string, items []string) (string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
		Size:  len(items),
	}

	_, result, err := prompt.Run()
	return result, err
}

func init() {
	rootCmd.AddCommand(editCmd)
}

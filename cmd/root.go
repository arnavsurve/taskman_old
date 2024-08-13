/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"taskman/task"
	"time"

	"github.com/manifoldco/promptui"
	_ "github.com/mattn/go-sqlite3"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	DB       *sql.DB
	dataFile = "tasks.db"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "taskman",
	Short: "taskman is a streamlined CLI task management solution",
	// 	Long: `A longer description that spans multiple lines and likely contains
	// examples and usage of using your application. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.taskman.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	var err error
	DB, err = sql.Open("sqlite3", dataFile)
	if err != nil {
		log.Fatal("Error opening database: ", err)
		os.Exit(1)
	}

	if err := task.InitializeDatabase(DB); err != nil {
		log.Fatal("Error initializing database:", err)
		os.Exit(1)
	}

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(editCmd)
}

var newCmd = &cobra.Command{
	Use:   "new \"task name\" \"due date (YYYY-MM-DD H:MM AM/PM)\"",
	Short: "Create a new task",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		date, err := time.Parse("2006-01-02 3:4 PM", args[1])
		if err != nil {
			fmt.Println("Invalid date format. Please use YYYY-MM-DD H:MM AM/PM.")
			return
		}
		task.AddTask(DB, name, date)
		fmt.Printf("Task (%s) added\n", name)
	},
}

var listCmd = &cobra.Command{
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
				task.DueDate.Format("2006-01-02 3:4 PM"),
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
		fmt.Println(taskID)

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

		newDueDateStr, err := promptForInput("Due Date (YYYY-MM-DD H:MM AM/PM)", selectedTask.DueDate.Format("2006-01-02 3:4 PM"))
		if err != nil {
			fmt.Println("Edit canceled.")
			return
		}
		newDueDate, err := time.Parse("2006-01-02 3:4 PM", newDueDateStr)
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
		taskOptions[i] = fmt.Sprintf("%d: %s [%s] Due: %s", task.ID, task.Name, task.Status, task.DueDate.Format("2006-01-02 3:4 PM"))
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

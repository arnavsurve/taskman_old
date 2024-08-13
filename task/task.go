package task

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

// TaskStatus represents the status of a task
type TaskStatus string

const (
	TODO        TaskStatus = "TODO"
	IN_PROGRESS            = "IN_PROGRESS"
	DONE                   = "DONE"
)

// Task represents a task with an ID, name, description, and status.
type Task struct {
	ID      int
	Name    string
	DueDate time.Time
	Status  TaskStatus
}

// InitializeDatabase initializes the tasks table if it does not already exist
func InitializeDatabase(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS tasks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT,
        due_date DATETIME,
        status TEXT
    );`
	_, err := db.Exec(query)
	return err
}

// AddTask adds a new task to the database
func AddTask(db *sql.DB, name string, date time.Time) error {
	query := "INSERT INTO tasks (name, due_date, status) VALUES (?, ?, ?)"
	_, err := db.Exec(query, name, date, TODO)
	return err
}

// ListTasks lists all tasks in the database
func ListTasks(db *sql.DB) ([]Task, error) {
	query := "SELECT id, name, due_date, status FROM tasks"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Name, &task.DueDate, &task.Status)
		if err != nil {
			log.Fatal(err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// UpdateTask updates the state of a task in the database
func UpdateTask(db *sql.DB, id int, name string, date time.Time, status TaskStatus) error {
	query := "UPDATE tasks SET name = ?, status = ?, due_date = ? WHERE id = ?"
	_, err := db.Exec(query, name, status, date, id)
	return err
}

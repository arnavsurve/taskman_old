package task

import "fmt"

// TaskStatus holds the completion status of a task
type TaskStatus string

const (
	TODO        TaskStatus = "TODO"
	IN_PROGRESS TaskStatus = "IN_PROGRESS"
	DONE        TaskStatus = "DONE"
)

// Task represents a task with respective properties
type Task struct {
	ID          int
	Name        string
	Description string
	Status      TaskStatus
}

// TaskManager managers a list of tasks
type TaskManager struct {
	Tasks  []Task
	NextID int
}

// NewTaskManager initializes a new TaskManager instance
func NewTaskManager() *TaskManager {
	return &TaskManager{
		Tasks:  []Task{},
		NextID: 1,
	}
}

// AddTask adds a new task to TaskManager
func (tm *TaskManager) AddTask(name string, desc string) {
	task := Task{
		ID:          tm.NextID,
		Name:        name,
		Description: desc,
		Status:      TODO,
	}
	tm.Tasks = append(tm.Tasks, task)
	tm.NextID++
}

// ListTasks lists all tasks in a TaskManager instance
func (tm *TaskManager) ListTasks() []Task {
	return tm.Tasks
}

// UpdateTaskStatus updates the completion status of a task
func (tm *TaskManager) UpdateTaskStatus(id int, status TaskStatus) bool {
	for i, task := range tm.Tasks {
		if len(tm.Tasks) == 0 {
			fmt.Println("No tasks found.")
		}
		if task.ID == id {
			tm.Tasks[i].Status = status
			return true
		}
	}
	return false
}

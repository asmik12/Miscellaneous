package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Task struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Done   bool   `json:"done"`
}

const taskFile = "tasks.json"

func loadTasks() ([]Task, error) {
	var tasks []Task
	file, err := os.ReadFile(taskFile)
	if err != nil {
		if os.IsNotExist(err) {
			return tasks, nil // return empty list
		}
		return nil, err
	}
	err = json.Unmarshal(file, &tasks)
	return tasks, err
}

func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(taskFile, data, 0644)
}

func addTask(title string) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}
	id := len(tasks) + 1
	tasks = append(tasks, Task{ID: id, Title: title})
	return saveTasks(tasks)
}

func listTasks() error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return nil
	}
	for _, task := range tasks {
		status := " "
		if task.Done {
			status = "✓"
		}
		fmt.Printf("[%s] %d: %s\n", status, task.ID, task.Title)
	}
	return nil
}

func completeTask(id int) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Done = true
			return saveTasks(tasks)
		}
	}
	return fmt.Errorf("task %d not found", id)
}

func deleteTask(id int) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}
	newTasks := []Task{}
	for _, task := range tasks {
		if task.ID != id {
			newTasks = append(newTasks, task)
		}
	}
	return saveTasks(newTasks)
}

package main

import (
	"fmt"
	"os"
	"strconv"
)

func printHelp() {
	fmt.Println("Todo CLI App")
	fmt.Println("Usage:")
	fmt.Println("  add <task>         Add a new task")
	fmt.Println("  list               List all tasks")
	fmt.Println("  done <id>          Mark task as done")
	fmt.Println("  delete <id>        Delete a task")
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	switch os.Args[1] {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a task.")
			return
		}
		title := os.Args[2]
		if err := addTask(title); err != nil {
			fmt.Println("Error:", err)
		}
	case "list":
		if err := listTasks(); err != nil {
			fmt.Println("Error:", err)
		}
	case "done":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a task ID.")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid ID.")
			return
		}
		if err := completeTask(id); err != nil {
			fmt.Println("Error:", err)
		}
	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a task ID.")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid ID.")
			return
		}
		if err := deleteTask(id); err != nil {
			fmt.Println("Error:", err)
		}
	default:
		printHelp()
	}
}

package main

import (
	"log"

	"github.com/matteoc91/simple-task-runner/simpletask"
	"github.com/matteoc91/simple-task-runner/taskmanager"
)

func main() {

	// Create task
	var simpletask simpletask.Task
	simpletask.Name = "task-001"
	simpletask.Title = "Task 001"

	// Store task
	_, err := taskmanager.Create(&simpletask)

	// Print error
	if err != nil {
		log.Fatal(err)
	}

	// Print simple task
	log.Println(simpletask)
}

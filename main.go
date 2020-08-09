package main

import (
	"flag"
	"log"

	"github.com/matteoc91/simple-task-runner/simpletask"
	"github.com/matteoc91/simple-task-runner/taskmanager"
)

func main() {

	// Create task
	var simpletask simpletask.Task

	// Read command line arguments
	bucket := flag.String("b", "default", "The bucket in which the task will be created")
	flag.StringVar(&(simpletask).Name, "n", "", "The name of the task, no white space should be set")
	flag.StringVar(&(simpletask).Title, "t", "", "The title of the task")
	flag.StringVar(&(simpletask).Description, "d", "", "the description of the task")
	deadline := flag.String("dl", "", "The deadline of the task, accepted format: yyyy-MM-dd")
	operation := flag.String("o", "read", "The operation to be performed, available: create, read, update, delete")
	flag.Parse()

	// Get reference to simpletask
	simpletaskReference := &simpletask

	// Validate command line options
	err := taskmanager.Validate(simpletaskReference, *deadline)

	// Evaluate errors
	if err != nil {
		log.Fatal(err)
	}

	// Evaluate CRUD operation
	if taskmanager.IsCreate(*operation) { // Create task
		err = taskmanager.Create(simpletaskReference, *bucket)
	} else if taskmanager.IsRead(*operation) { // Read task
		simpletaskReference, err = taskmanager.Read(simpletask.Name, *bucket)
	} else if taskmanager.IsUpdate(*operation) { // Update task
		log.Println("To be implemented")
	} else if taskmanager.IsDelete(*operation) { // Delete task
		log.Println("To be implemented")
	} else { // Operation unknown
		log.Fatal("operation: unknown")
	}

	// Print error
	if err != nil {
		log.Fatal(err)
	}

	// Print simple task
	log.Println(*simpletaskReference)
}

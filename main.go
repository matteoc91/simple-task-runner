package main

import (
	"errors"
	"flag"
	"log"

	"github.com/matteoc91/simple-task-runner/simpletask"
	"github.com/matteoc91/simple-task-runner/taskmanager"
)

func main() {

	// Create task
	var requestedTask simpletask.Task

	// Read command line arguments
	bucket := flag.String("b", "default", "The bucket in which the task will be created")
	flag.StringVar(&(requestedTask).Name, "n", "", "The name of the task, no white space should be set")
	flag.StringVar(&(requestedTask).Title, "t", "", "The title of the task")
	flag.StringVar(&(requestedTask).Description, "d", "", "the description of the task")
	deadline := flag.String("dl", "", "The deadline of the task, accepted format: yyyy-MM-dd")
	operation := flag.String("o", "read", "The operation to be performed, available: create, read, update, delete")
	flag.Parse()

	// Validate command line options
	err := taskmanager.Validate(&requestedTask, *deadline)

	// Evaluate errors
	if err != nil {
		log.Fatal(err)
	}

	// Evaluate CRUD operation
	if taskmanager.IsCreate(*operation) { // Create task

		// Create
		err = taskmanager.Create(&requestedTask, *bucket)
		// Print the created task
		log.Println("### Created task:")
		log.Println(requestedTask)

	} else if taskmanager.IsRead(*operation) { // Read task

		// Create a task slice
		var readedTasks []simpletask.Task
		// Read task(s)
		readedTasks, err = taskmanager.Read(requestedTask.Name, *bucket)
		if err == nil {
			// Print task(s)
			log.Println("### Readed task(s):")
			for _, v := range readedTasks {
				log.Println(v)
			}
		}

	} else if taskmanager.IsUpdate(*operation) { // Update task

		log.Println("To be implemented")

	} else if taskmanager.IsDelete(*operation) { // Delete task

		err = taskmanager.Delete(requestedTask.Name, *bucket)

	} else { // Operation unknown

		err = errors.New("operation: unknown")

	}

	// Print error
	if err != nil {
		log.Fatal(err)
	}
}

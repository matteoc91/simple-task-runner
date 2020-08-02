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
	flag.Parse()

	// Validate command line options
	err := taskmanager.Validate(&simpletask, *deadline)

	// Evaluate errors
	if err != nil {
		log.Fatal(err)
	}

	// Store task
	err = taskmanager.Create(&simpletask, *bucket)

	// Print error
	if err != nil {
		log.Fatal(err)
	}

	// Print simple task
	log.Println(simpletask)
}

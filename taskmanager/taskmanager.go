package taskmanager

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/matteoc91/simple-task-runner/simpletask"
)

const dbname string = "simple-task-runner-workarea.db"

// Validate validates the command line arguments
func Validate(simpleTask *simpletask.Task, deadline string) error {

	/*
	***************************
	****** Validate name ******
	***************************
	 */

	// Mandatory
	if simpleTask.Name != "" {

		// No whitespaces
		re := regexp.MustCompile(`(\s){1}`)
		if re.MatchString(simpleTask.Name) {
			return errors.New("name: whitespaces not allowed")
		}
	}

	/*
	***************************
	**** Validate deadline ****
	***************************
	 */

	var err error
	if deadline != "" {
		var dl time.Time
		dl, err = time.Parse("2006-01-02", deadline)
		simpleTask.Deadline = dl
	}

	return err
}

// Create creates a simple task if not present
func Create(task *simpletask.Task, bucket string) error {

	// Name should not be empty
	if task.Name == "" {
		return errors.New("name: should not be empty")
	}

	// Open DB
	db, err := bolt.Open(dbname, 0600, nil)

	// Defer close
	defer db.Close()

	// Return error
	if err != nil {
		return err
	}

	// Create into bucket
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return fmt.Errorf("Create bucket: %s", err)
		}
		task.Created = time.Now()
		task.Updated = time.Now()
		v, _ := json.Marshal(&task)
		b.Put([]byte(task.Name), v)
		return nil
	})
}

// Read reads from bucket the given task
func Read(name string, bucket string) ([]simpletask.Task, error) {

	// Open DB
	db, err := bolt.Open(dbname, 0600, nil)

	// Defer close
	defer db.Close()

	// Return error
	if err != nil {
		return nil, err
	}

	// Define variables
	var simpleTasks []simpletask.Task
	var simpleTask simpletask.Task

	// Read from bucket
	err = db.View(func(tx *bolt.Tx) error {

		// Get the bucket
		b := tx.Bucket([]byte(bucket))
		if err != nil {
			return fmt.Errorf("Read bucket: %s", err)
		}

		if name != "" {

			// A name has been supplied, look for it
			v := b.Get([]byte(name))

			// No task found
			if v == nil {
				return errors.New("No task found")
			}

			// No error to be raised
			json.Unmarshal(v, &simpleTask)
			simpleTasks = append(simpleTasks, simpleTask)
			return nil
		}

		// No name has been supplied, look for each task
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			json.Unmarshal(v, &simpleTask)
			simpleTasks = append(simpleTasks, simpleTask)
		}
		return nil
	})

	return simpleTasks, err
}

// Delete deletes a task by its name
func Delete(name string, bucket string) error {

	// Name should not be empty
	if name == "" {
		return errors.New("name: should not be empty")
	}

	// Open DB
	db, err := bolt.Open(dbname, 0600, nil)

	// Defer close
	defer db.Close()

	// Return error
	if err != nil {
		return err
	}

	// Try to read bucket and delete record
	return db.Update(func(tx *bolt.Tx) error {

		// Get the bucket
		b := tx.Bucket([]byte(bucket))
		if err != nil {
			return fmt.Errorf("Read bucket: %s", err)
		}

		// Try to delete
		return b.Delete([]byte(name))
	})
}

// Update updates a task
func Update(task *simpletask.Task, bucket string) (*simpletask.Task, error) {

	// Name should not be empty
	if task.Name == "" {
		return nil, errors.New("name: should not be empty")
	}

	// Open DB
	db, err := bolt.Open(dbname, 0600, nil)

	// Defer close
	defer db.Close()

	// Return error
	if err != nil {
		return nil, err
	}

	// Define task
	var simpleTask simpletask.Task

	// Try to update
	err = db.Update(func(tx *bolt.Tx) error {

		// Get the bucket
		b := tx.Bucket([]byte(bucket))
		if err != nil {
			return fmt.Errorf("Read bucket: %s", err)
		}

		// Get task
		v := b.Get([]byte(task.Name))
		if v == nil {
			return fmt.Errorf("No task found")
		}
		json.Unmarshal(v, &simpleTask)

		// Update title
		if task.Title != "" {
			simpleTask.Title = task.Title
		}

		// Update description
		if task.Description != "" {
			simpleTask.Description = task.Description
		}

		// Update deadline
		if !task.Deadline.IsZero() {
			simpleTask.Deadline = task.Deadline
		}

		// Update updated
		simpleTask.Updated = time.Now()

		// Try to update
		v, _ = json.Marshal(&simpleTask)
		return b.Put([]byte(simpleTask.Name), v)
	})

	return &simpleTask, err
}

// IsCreate check if operation is create
func IsCreate(operation string) bool {
	return isRequestedOperation(operation, "create")
}

// IsRead check if operation is read
func IsRead(operation string) bool {
	return isRequestedOperation(operation, "read")
}

// IsUpdate check if operation is update
func IsUpdate(operation string) bool {
	return isRequestedOperation(operation, "update")
}

// IsDelete check if operation is delete
func IsDelete(operation string) bool {
	return isRequestedOperation(operation, "delete")
}

// isRequestedOperation check if operation is requested operation
func isRequestedOperation(operation string, requestedOperation string) bool {
	return strings.EqualFold(operation, requestedOperation)
}

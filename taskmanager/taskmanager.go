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

// dbname: database name
const dbname string = "simple-task-runner-workarea.db"

// txOperation: transactional operation
type txOperation func(b *bolt.Bucket) error

// AddComment adds a comment to the list
func AddComment(comments []simpletask.Comment, commentText string) []simpletask.Comment {
	if commentText != "" {
		var comment simpletask.Comment
		comment.Text = commentText
		comment.Created = time.Now()
		comments = append([]simpletask.Comment{comment}, comments...)
	}
	return comments
}

// Validate validates the command line arguments
func Validate(simpleTask *simpletask.Task, deadline string) error {

	// Validate name
	if simpleTask.Name != "" {

		// No whitespaces
		re := regexp.MustCompile(`(\s){1}`)
		if re.MatchString(simpleTask.Name) {
			return errors.New("name: whitespaces not allowed")
		}
	}

	// Validate deadline
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

	// Define create operation
	var create txOperation = func(b *bolt.Bucket) error {
		task.Created = time.Now()
		task.Updated = time.Now()
		v, _ := json.Marshal(&task)
		b.Put([]byte(task.Name), v)
		return nil
	}

	// Open in write mode
	return openInWriteMode(create, bucket)
}

// Read reads from bucket the given task
func Read(name string, bucket string) ([]simpletask.Task, error) {

	// Define variables
	var simpleTasks []simpletask.Task
	var simpleTask simpletask.Task

	// Define read operation
	var read txOperation = func(b *bolt.Bucket) error {

		// If name is supplied, check for single task
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
	}

	// Open in read mode
	err := openInReadMode(read, bucket)

	return simpleTasks, err
}

// Delete deletes a task by its name
func Delete(name string, bucket string) error {

	// Name should not be empty
	if name == "" {
		return errors.New("name: should not be empty")
	}

	// Define delete operation
	var delete txOperation = func(b *bolt.Bucket) error {
		return b.Delete([]byte(name))
	}

	// Open in write mode
	return openInWriteMode(delete, bucket)
}

// Update updates a task
func Update(task *simpletask.Task, bucket string) (*simpletask.Task, error) {

	// Name should not be empty
	if task.Name == "" {
		return nil, errors.New("name: should not be empty")
	}

	// Define task
	var simpleTask simpletask.Task

	// Define update operation
	var update txOperation = func(b *bolt.Bucket) error {

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

		// Add a comment
		if len(task.Comments) > 0 {
			simpleTask.Comments = AddComment(simpleTask.Comments, task.Comments[0].Text)
		}

		// Update updated
		simpleTask.Updated = time.Now()

		// Try to update
		v, _ = json.Marshal(&simpleTask)
		return b.Put([]byte(simpleTask.Name), v)
	}

	// Open in write mode
	err := openInWriteMode(update, bucket)

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

// openInWriteMode opens the connection in write mode
func openInWriteMode(txOp txOperation, bucket string) error {

	// Open DB
	db, err := bolt.Open(dbname, 0600, nil)

	// Defer close
	defer db.Close()

	// Return error
	if err != nil {
		return err
	}

	// Update operation
	return db.Update(func(tx *bolt.Tx) error {

		// Create bucket if not exist
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return fmt.Errorf("Create bucket: %s", err)
		}

		// Perform txOp
		return txOp(b)
	})
}

// openInReadMode opens the connection in read mode
func openInReadMode(txOp txOperation, bucket string) error {

	// Open DB
	db, err := bolt.Open(dbname, 0600, nil)

	// Defer close
	defer db.Close()

	// Return error
	if err != nil {
		return err
	}

	// Read from bucket
	return db.View(func(tx *bolt.Tx) error {

		// Get the bucket
		b := tx.Bucket([]byte(bucket))
		if err != nil {
			return fmt.Errorf("Read bucket: %s", err)
		}

		// Perform txOp
		return txOp(b)
	})
}

// isRequestedOperation check if operation is requested operation
func isRequestedOperation(operation string, requestedOperation string) bool {
	return strings.EqualFold(operation, requestedOperation)
}

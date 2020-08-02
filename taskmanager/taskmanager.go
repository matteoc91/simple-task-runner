package taskmanager

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/boltdb/bolt"
	"github.com/matteoc91/simple-task-runner/simpletask"
)

const dbname string = "simple-task-runner-workarea.db"

// Validate the command line arguments
func Validate(simpleTask *simpletask.Task, deadline string) error {

	/*
	***************************
	****** Validate name ******
	***************************
	 */

	// Mandatory
	if simpleTask.Name == "" {
		return errors.New("name: should not be empty")
	}

	// No whitespaces
	re := regexp.MustCompile(`(\s){1}`)
	if re.MatchString(simpleTask.Name) {
		return errors.New("name: whitespaces not allowed")
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

// Create creates a simple task if not present.
func Create(task *simpletask.Task, bucket string) error {

	// Open DB
	db, err := bolt.Open(dbname, 0600, nil)

	// Defer close
	defer db.Close()

	// Print error
	if err != nil {
		return err
	}

	// Insert into bucket
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return fmt.Errorf("Create bucket: %s", err)
		}
		task.Created = time.Now()
		v, _ := json.Marshal(&task)
		b.Put([]byte(task.Name), v)
		return nil
	})

	// Return data
	return err
}

package taskmanager

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	"github.com/matteoc91/simple-task-runner/simpletask"
)

const dbname string = "simple-task-runner.db"
const workarea string = "personal-tasks"

// Create creates a simple task if not present.
func Create(task *simpletask.Task) (*simpletask.Task, error) {

	// Open DB
	db, err := bolt.Open(dbname, 0600, nil)

	// Defer close
	defer db.Close()

	// Print error
	if err != nil {
		return task, err
	}

	// Insert into bucket
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(workarea))
		if err != nil {
			return fmt.Errorf("Create bucket: %s", err)
		}
		task.Created = time.Now()
		v, _ := json.Marshal(&task)
		b.Put([]byte(task.Name), v)
		return nil
	})

	// Return data
	return task, err
}

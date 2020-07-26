package simpletask

import "time"

// Task is a simple task definition.
// It is defined to be shared between packages.
type Task struct {
	Name        string
	Title       string
	Description string
	Created     time.Time
	Deadline    time.Time
	Comments    []Comment
}

// Comment is a simple comment definition for Tasks.
// It is defined to be shared between Tasks.
type Comment struct {
	Text    string
	Created time.Time
}

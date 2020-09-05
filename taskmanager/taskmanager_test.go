package taskmanager

import (
	"testing"
	"time"

	"github.com/matteoc91/simple-task-runner/simpletask"
)

// TestValidate validates a task
func TestValidate(t *testing.T) {

	// Create a task
	var task simpletask.Task
	task.Name = "MyName"
	task.Description = "MyDescription"

	// Validate without deadline
	err := Validate(&task, "")
	if err != nil {
		t.Errorf("Validate %v, expects ok", err.Error())
	}

	// Validate with deadline
	err = Validate(&task, "2020-01-01")
	if err != nil {
		t.Errorf("Validate %v, expects ok", err.Error())
	}

	// Create a task with spaced name
	var taskWithSpacedName simpletask.Task
	taskWithSpacedName.Name = "My Name"
	taskWithSpacedName.Description = "MyDescription"

	// Validate white spaced
	err = Validate(&taskWithSpacedName, "")
	if err == nil {
		t.Errorf("Validate ok, expects error 'name with whitespaces'")
	}

	// Validate a wrong deadline
	err = Validate(&task, "02-09-2020")
	if err == nil {
		t.Errorf("Validate ok, expects error 'wrong deadline date format'")
	}
}

// TestCreates creates a task
func TestCreate(t *testing.T) {

	// Create a task
	var task simpletask.Task
	task.Name = "TestTask"
	task.Description = "Test Task"

	// Validate the task
	err := Validate(&task, "")
	if err != nil {
		t.Errorf("Create %v, expects ok", err.Error())
	}

	// Create the task
	err = Create(&task, "default")
	if err != nil {
		t.Errorf("Create %v, expects ok", err.Error())
	}

	// Create a task without name
	var taskWithoutName simpletask.Task
	taskWithoutName.Description = "Test Task"

	// Validate the task
	err = Validate(&taskWithoutName, "")
	if err != nil {
		t.Errorf("Create %v, expects ok", err.Error())
	}

	// Attempt to create the task
	err = Create(&taskWithoutName, "default")
	if err == nil {
		t.Errorf("Create ok, expects error 'name should not be empty'")
	}
}

// TestRead reads tasks
func TestRead(t *testing.T) {

	// Create a task
	var task simpletask.Task
	task.Name = "TestTask"
	task.Description = "Test Task"

	// Validate the task
	err := Validate(&task, "")
	if err != nil {
		t.Errorf("Read %v, expects ok", err.Error())
	}

	// Create the task
	err = Create(&task, "default")
	if err != nil {
		t.Errorf("Read %v, expects ok", err.Error())
	}

	// Read a task
	tasks, readError := Read("TestTask", "default")
	if readError != nil {
		t.Errorf("Read %v, expects ok", err.Error())
	}
	if len(tasks) != 1 {
		t.Errorf("Read n. %d, expects 1", len(tasks))
	}

	// Read tasks
	tasks, readError = Read("", "default")
	if readError != nil {
		t.Errorf("Read %v, expects ok", err.Error())
	}
	if len(tasks) < 1 {
		t.Errorf("Read n. %d, expects at least 1", len(tasks))
	}

	// Read a no-existing task
	tasks, readError = Read("____TestTask", "default")
	if readError == nil {
		t.Errorf("Read ok, expects 'no task found'")
	}
}

// TestDelete delets a task
func TestDelete(t *testing.T) {

	// Create a task
	var task simpletask.Task
	task.Name = "TestTask"
	task.Description = "Test Task"

	// Validate the task
	err := Validate(&task, "")
	if err != nil {
		t.Errorf("Delete %v, expects ok", err.Error())
	}

	// Create the task
	err = Create(&task, "default")
	if err != nil {
		t.Errorf("Delete %v, expects ok", err.Error())
	}

	// Try to delete
	err = Delete(task.Name, "default")
	if err != nil {
		t.Errorf("Delete %v, expects ok", err.Error())
	}

	// Delete with error
	err = Delete("", "default")
	if err == nil {
		t.Errorf("Delete ok, expects ko")
	}
}

// TestUpdate updates a task
func TestUpdate(t *testing.T) {

	// Create a task
	var task simpletask.Task
	task.Name = "TestTask"
	task.Description = "Test Task"

	// Validate the task
	err := Validate(&task, "")
	if err != nil {
		t.Errorf("Update %v, expects ok", err.Error())
	}

	// Create the task
	err = Create(&task, "default")
	if err != nil {
		t.Errorf("Update %v, expects ok", err.Error())
	}

	// Update task
	task.Title = "New Title"
	task.Description = "New Description"
	task.Deadline = time.Now()
	var updatedTask *simpletask.Task
	updatedTask, err = Update(&task, "default")
	if err != nil {
		t.Errorf("Update %v, expects ok", err.Error())
	}
	if updatedTask.Title != task.Title {
		t.Errorf("Update Title not updated, expects to match")
	}
	if updatedTask.Description != task.Description {
		t.Errorf("Update Description not updated, expects to match")
	}
	if updatedTask.Deadline != task.Deadline {
		t.Errorf("Update Deadline not updated, expects to match")
	}

	// Create a task without name
	var taskWithoutName simpletask.Task
	taskWithoutName.Description = "Test Task"

	// Validate the task
	err = Validate(&taskWithoutName, "")
	if err != nil {
		t.Errorf("Update %v, expects ok", err.Error())
	}

	// Attempt to update the task
	_, err = Update(&taskWithoutName, "default")
	if err == nil {
		t.Errorf("Update ok, expects error 'name should not be empty'")
	}
}

// TestIsCreate validates a create operation
func TestIsCreate(t *testing.T) {
	if !IsCreate("create") {
		t.Errorf("IsCreate false, expects true")
	}
	if IsCreate("operation") {
		t.Errorf("IsCreate true, expects false")
	}
}

// TestIsRead validates a read operation
func TestIsRead(t *testing.T) {
	if !IsRead("read") {
		t.Errorf("IsRead false, expects true")
	}
	if IsRead("operation") {
		t.Errorf("IsRead true, expects false")
	}
}

// TestIsUpdate validates an update operation
func TestIsUpdate(t *testing.T) {
	if !IsUpdate("update") {
		t.Errorf("IsUpdate false, expects true")
	}
	if IsUpdate("operation") {
		t.Errorf("IsUpdate true, expects false")
	}
}

// TestIsDelete validates a delete operation
func TestIsDelete(t *testing.T) {
	if !IsDelete("delete") {
		t.Errorf("IsDelete false, expects true")
	}
	if IsDelete("operation") {
		t.Errorf("IsDelete true, expects false")
	}
}

package client

import (
	"github.com/stehrn/hpc-poc/internal/utils"
)

// Task a task represents a unit of work, it has an identiy and some data
type Task interface {
	Data() []byte
	ObjectPath() *ObjectPath
	AddError(err error)
	Errors() []error
}

// LocalTask a client side (local) task
type LocalTask struct {
	job        Job
	ID         string
	data       []byte
	taskErrors []error
}

// newTask create new task with given payload, id will be automatically generated
func newTask(job Job, data []byte) *LocalTask {
	return &LocalTask{
		job:  job,
		ID:   utils.GenerateID(),
		data: data,
	}
}

// ObjectPath location for given task
func (t *LocalTask) ObjectPath() *ObjectPath {
	return ObjectPathForTask(t.job.ObjectPath(), t.ID)
}

// Data binary data
func (t *LocalTask) Data() []byte {
	return t.data
}

// AddError add error to task
func (t *LocalTask) AddError(err error) {
	t.taskErrors = append(t.taskErrors, err)
}

func (t *LocalTask) hasError() bool {
	return len(t.taskErrors) != 0
}

// Errors get errors associated with task
func (t *LocalTask) Errors() []error {
	return t.taskErrors
}

package client

import (
	"log"

	"github.com/stehrn/hpc-poc/internal/utils"
)

// Job represents a unit of work, made up of one or more tasks
type Job interface {
	CreateTask(data []byte) Task
	AddTask(task Task)
	Tasks() []Task
	TaskIterator() TaskDataSourceIterator
	ObjectPath() *ObjectPath
	SetState(state State)
	Errors() map[Task][]error
	Close()
}

// LocalJob client side job
type LocalJob struct {
	Session *LocalSession
	Name    string
	ID      string
	State   State
	tasks   []Task
}

// Errors errors associated with job (and tasks)
type Errors struct {
	jobError   error
	taskErrors map[*Task][]error
}

// State state of Job
type State string

const (
	// Initial job just been created
	Initial State = "Job Created"
	// TaskDataUploading task data uploading
	TaskDataUploading State = "Task Data Uploading"
	// TaskDataUploaded task data uploaded
	TaskDataUploaded State = "Task Data Uploaded"
	// TaskDataUploadError error uploading taks data
	TaskDataUploadError State = "Task Data Upload Error"
	// JobMessagePublishing punblishing job message
	JobMessagePublishing State = "Job Message Publishing"
	// JobMessagePublished job message published
	JobMessagePublished State = "Job Message Published"
	// JobMessagePublishError error publishing job message
	JobMessagePublishError State = "Job Message Publishing Error"
)

// NewJob create a new Job with given name, id will be automatically generated
func NewJob(name string, session *LocalSession) *LocalJob {
	job := &LocalJob{
		Session: session,
		Name:    name,
		State:   Initial,
		ID:      utils.GenerateID(),
	}
	session.AddJob(job)
	return job
}

// CreateTask create a new task for given data and adds it to this job
func (j *LocalJob) CreateTask(data []byte) Task {
	task := newTask(j, data)
	j.tasks = append(j.tasks, task)
	return task
}

// AddTask add a task
func (j *LocalJob) AddTask(task Task) {
	j.tasks = append(j.tasks, task)
}

// Tasks get tasks associated with job
func (j *LocalJob) Tasks() []Task {
	return j.tasks
}

// TaskIterator return a task iterator
func (j *LocalJob) TaskIterator() TaskDataSourceIterator {
	return TaskDataSourceIterator(j.Tasks())
}

// ObjectPath location for given job (parent of tasks)
func (j *LocalJob) ObjectPath() *ObjectPath {
	return ObjectPathForJob(j.Session.Business, j.Session.Name, j.Name)
}

// SetState set state of job
// scope to use this as entry point for tracing
func (j *LocalJob) SetState(state State) {
	j.State = state
}

// Errors get back any errors associates with the Job
func (j *LocalJob) Errors() map[Task][]error {
	errors := make(map[Task][]error)
	for _, task := range j.tasks {
		taskErrors := task.Errors()
		if taskErrors != nil {
			errors[task] = taskErrors
		}
	}
	return errors
}

// HasErrors do we have any errors?
func (j *LocalJob) HasErrors() bool {
	for _, task := range j.tasks {
		taskErrors := task.Errors()
		if taskErrors != nil && len(taskErrors) != 0 {
			return true
		}
	}
	return false
}

// Close close off job
func (j *LocalJob) Close() {
	log.Printf("Closing job %q\n", j.Name)
}

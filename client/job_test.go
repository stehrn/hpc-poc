package client

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const project = "project1"
const business = "bu1"

func session(name string) *LocalSession {
	return NewSession(name, business)
}

func TestJobCreation(t *testing.T) {

	// create a session and add 2 jobs
	session1 := session("session 1")
	defer session1.Destroy()
	job1 := NewJob("Test Job 1", session1)

	session1jobs := session1.Jobs()
	assert.Equal(t, 1, len(session1jobs), fmt.Sprintf("Expected just 1 job in session 1, got %d", len(session1jobs)))
	assert.Equal(t, job1, session1jobs[0], "Unexpected job")

	job2 := NewJob("Test Job 2", session1)
	session1jobs = session1.Jobs()
	assert.Equal(t, 2, len(session1jobs), fmt.Sprintf("Expected 2 jobs in session 1, got %d", len(session1jobs)))
	assert.Equal(t, job2, session1jobs[1], "Unexpected job")

	// crea a new session, add job
	session2 := session("session 2")
	defer session2.Destroy()
	job3 := NewJob("Test Job 3", session2)
	session1jobs = session1.Jobs()
	assert.Equal(t, 2, len(session1jobs), fmt.Sprintf("Expected 2 jobs in session 1, got %d", len(session1jobs)))
	session2jobs := session2.Jobs()
	assert.Equal(t, 1, len(session2jobs), fmt.Sprintf("Expected 1 job in session 2, got %d", len(session2jobs)))
	assert.Equal(t, job3, session2jobs[0], "Unexpected job")
}

func TestTaskCreation(t *testing.T) {
	session := session("session A")
	defer session.Destroy()
	job := NewJob("Test Job", session)

	task1 := job.CreateTask([]byte("ABC€"))
	task2 := job.CreateTask([]byte("123€"))

	tasks := job.Tasks()

	assert.Equal(t, 2, len(tasks), "Expected 2 tasks")
	assert.Equal(t, task1, tasks[0], "Unexpected task")
	assert.Equal(t, task2, tasks[1], "Unexpected task")
	assert.Equal(t, len(session.Jobs()), 1, "Expected 1 job")
	assert.Equal(t, "bu1/session A/Test Job", job.ObjectPath().String(), "Unexpected job location")
}

func TestErrorHandling(t *testing.T) {
	session := session("session A")
	defer session.Destroy()
	job := NewJob("Error Test Job", session)

	errors := job.Errors()
	assert.Equal(t, 0, len(errors), "Expected zero errors")

	task := job.CreateTask(nil)
	error1 := fmt.Errorf("error 1")
	task.AddError(error1)

	errors = job.Errors()
	assert.Equal(t, 1, len(errors), "Expected 1 task with errors for job")
	assert.Equal(t, 1, len(errors[task]), "Expected 1 error for task")
	assert.Equal(t, error1, errors[task][0], "Unexpected error for task")

	error2 := fmt.Errorf("error 2")
	task.AddError(error2)
	errors = job.Errors()
	assert.Equal(t, 1, len(errors), "Expected 1 task with errors for job")
	assert.Equal(t, 2, len(errors[task]), "Expected 2 errors for task")
	assert.Equal(t, error2, errors[task][1], "Unexpected error for task")

}

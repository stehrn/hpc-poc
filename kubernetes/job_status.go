package kubernetes

import (
	batchv1 "k8s.io/api/batch/v1"
)

// JobState state of a Job
type JobState string

const (
	// RunningWithFailures means the job is still running, and has had some failures
	RunningWithFailures JobState = "Running (with Failures)"
	// RunningNoFailures means the job is still running, and has so far had no failures
	RunningNoFailures JobState = "Running (no Failures)"
	// Complete means the job has succesfully completed its execution.
	Complete JobState = "Complete"
	// Failed means the job has failed its execution.
	Failed JobState = "Failed"
)

// GetJobState get state of the job - a summary of various conditions to make it easier to assertain state of things
func GetJobState(status batchv1.JobStatus) JobState {
	if hasCondition(status.Conditions, batchv1.JobComplete) {
		return Complete
	} else if hasCondition(status.Conditions, batchv1.JobFailed) {
		return Failed
	}
	if status.Failed == 0 {
		return RunningNoFailures
	}
	return RunningWithFailures
}

// HasFailures any failures? regarldess of running or not
func (s JobState) HasFailures() bool {
	return s == RunningWithFailures || s == Failed
}

// IsRunning is job still running
func (s JobState) IsRunning() bool {
	return s == RunningNoFailures || s == RunningWithFailures
}

// SUCCESS return true is job has completed succesfully, if it returns false, job may still be running
func SUCCESS(status batchv1.JobStatus) bool {
	return GetJobState(status) == Complete
}

// FINISHED return true is Job has finished (either succesfully or failed)
func FINISHED(status batchv1.JobStatus) (JobState, bool) {
	state := GetJobState(status)
	finished := !state.IsRunning()
	return state, finished
}

// ANY just always return true, regardless of status
func ANY(status batchv1.JobStatus) bool {
	return true
}

func hasCondition(conditions []batchv1.JobCondition, conditionType batchv1.JobConditionType) bool {
	for _, condition := range conditions {
		if condition.Type == conditionType {
			return true
		}
	}
	return false
}

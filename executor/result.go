package executor

import (
	"fmt"
	"log"
	"sync"

	"github.com/stehrn/hpc-poc/kubernetes"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Result result of execution
type Result struct {
	Error     error
	namespace string
	messageID string
}

// ErrorResult a result in error
func ErrorResult(err error) *Result {
	return &Result{err, "", ""}
}

// Watch watch for result
// TODO: change this to watch for JOB ID
func (r *Result) Watch() error {
	if r.Error != nil {
		return fmt.Errorf("Could not watch for result, already in error: %w", r.Error)
	}

	jobWatcher, err := kubernetes.NewClient(r.namespace)
	if err != nil {
		return err
	}

	log.Printf("Listening to subscription ID %q", r.messageID)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		options := metav1.ListOptions{LabelSelector: fmt.Sprintf("gcp.pubsub.subscription.id=%s", r.messageID)}
		err = jobWatcher.Watch(options, kubernetes.ANY, func(job *batchv1.Job) {
			state, done := kubernetes.FINISHED(job.Status)
			log.Printf("Received update for Job %q, status: %v", job.Name, state)
			if done {
				log.Printf("Job %q finished", job.Name)
				wg.Done()
			}
		})
	}()
	wg.Wait()
	return err
}

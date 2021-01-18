package executor

import (
	"fmt"
	"log"

	"github.com/pkg/errors"

	"github.com/stehrn/hpc-poc/client"
	"github.com/stehrn/hpc-poc/gcp/pubsub"
	"github.com/stehrn/hpc-poc/gcp/storage"
)

var pubsubClient *pubsub.Client

var (
	// ErrorUploadingTaskData error uploading task data to cloud storage
	ErrorUploadingTaskData = errors.New("Error uploading task data, see job for details")
	// ErrorPublishJobStorageLocation error publishing job storage location
	ErrorPublishJobStorageLocation = errors.New("Error publishing job storage location")
	// ErrorDeletingTaskDataObjectNotFound error  deleting task storage data - object not found
	ErrorDeletingTaskDataObjectNotFound = errors.New("Error deleting task data: Storage Object not found")
)

// Executor client
type Executor struct {
	*GcpContext
	Storage storage.ClientInterface
}

// New create new executor client
func New(gcpContext *GcpContext) (*Executor, error) {
	storageClient, err := gcpContext.NewStorageClient()
	if err != nil {
		return nil, err
	}

	pubsubClient, err = pubsub.NewPubClient(gcpContext.Project)
	if err != nil {
		return nil, err
	}

	return &Executor{gcpContext, storageClient}, nil
}

// Execute run a job
func (e Executor) Execute(job client.Job) *Result {
	fmt.Printf("Executing job, name: %q, ID: %s)\n", job.Name(), job.ID())

	// upload task data to cloud storage
	job.SetState(client.TaskDataUploading)
	e.Storage.UploadMany(job.TaskIterator())
	if job.HasErrors() {
		// fail fast here, even though some of task data may have been uploaded.
		// an alternative strategy might continue anyway ...
		job.SetState(client.TaskDataUploadError)
		return ErrorResult(ErrorUploadingTaskData)
	}
	job.SetState(client.TaskDataUploaded)

	// send message with location of cloud storage for job
	job.SetState(client.JobMessagePublishing)
	location := e.Location(job.ObjectPath().String())
	messageID, err := e.publishJobStorageLocation(e.TopicName(), location)
	if err != nil {
		job.SetState(client.JobMessagePublishError)
		return ErrorResult(errors.Wrap(err, ErrorPublishJobStorageLocation.Error()))
	}
	job.SetState(client.JobMessagePublished)

	return &Result{
		Error:     nil,
		namespace: e.Namespace,
		messageID: messageID}
}

// Cancel cancel a job
func (e Executor) Cancel(job client.Job) error {
	return e.deleteTaskData(job)
	// TODO: send message to k8 to cancel job
}

// deleteTaskData delete _all_ data associated with a job, called when job is cancelled
func (e Executor) deleteTaskData(job client.Job) error {
	location := e.Location(asDirectory(job.ObjectPath().String()))
	log.Printf("Deleting data at: '%v'\n", location)
	err := e.Storage.Delete(location)
	if err != nil {
		return fmt.Errorf("%q: %w", location, ErrorDeletingTaskDataObjectNotFound)
	}
	return nil
}

func asDirectory(path string) string {
	return fmt.Sprintf("%s/", path)
}

// publishJobStorageLocation location to topic
func (e Executor) publishJobStorageLocation(topicName string, location storage.Location) (string, error) {
	log.Printf("Publishing location (%v) to topic: %q\n", location, topicName)
	bytes, err := location.ToBytes()
	if err != nil {
		return "", err
	}
	id, err := pubsubClient.Publish(topicName, bytes)
	if err != nil {
		return "", err
	}
	return id, nil
}

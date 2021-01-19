// Integration test for executor
//
// Following env variables required:
//
// export PROJECT_NAME=hpc-poc
// export CLOUD_STORAGE_BUCKET_NAME=hpc-poc-bucket
// export GOOGLE_APPLICATION_CREDENTIALS=${HOME}/integration_test_key.json
//
package executor

import (
	"bytes"
	ctx "context"
	"fmt"
	"os"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/stehrn/hpc-poc/client"
	messaging "github.com/stehrn/hpc-poc/gcp/pubsub"
	"github.com/stehrn/hpc-poc/gcp/storage"
)

const business = "integration-executor-test"

var (
	project       = os.Getenv("PROJECT_NAME")
	bucketName    = storage.BucketNameFromEnv()
	storageClient storage.ClientInterface
	subClient     *messaging.TempPubSub
	topicName     = fmt.Sprintf("%s-integration-executor-test-topic", project)
)

func TestMain(m *testing.M) {
	code := m.Run()
	teardown()
	os.Exit(code)
}

// Test upload of job data to cloud storage
func TestJobExecute(t *testing.T) {

	initClients(t)

	// create a sessiob, and job with 2 tasks
	session := client.NewSession("session-1", business)
	defer session.Destroy()
	job := session.NewJob("test-job-1")
	task1 := job.CreateTask([]byte("123"))
	task2 := job.CreateTask([]byte("abc"))

	gcpContext := &GcpContext{
		Project:    project,
		BucketName: bucketName,
		Business:   business}

	var err error

	// create new executor and execute the job
	exe, err := New(gcpContext)
	if err != nil {
		t.Fatalf("Error creating client, %v", err)
	}
	result := exe.Execute(job)

	// expect 2 storage objects and a message published to topic giving location of data
	err = result.Error
	if err != nil {
		t.Error("Could not execute job", err)
	} else {
		verifyData(t, storage.NewLocation(bucketName, task1.ObjectPath().String()), "123")
		verifyData(t, storage.NewLocation(bucketName, task2.ObjectPath().String()), "abc")
		verifySubscription(t)
	}
}

func initClients(t *testing.T) {
	var err error
	sub, err := messaging.NewSubClient(project, topicName)
	if err != nil {
		t.Fatalf("Could not create sub client for project %s, topic %s, error: %v", project, topicName, err)
	}
	ID := fmt.Sprintf("%s-integration-executor-test", project)
	subClient, err = sub.NewTempPubSub(ID)
	if err != nil {
		t.Fatalf("Could not create NewTempPubSub client for ID %s, error: %v", ID, err)
	}
	// used to verify job data uploaded
	storageClient, err = storage.NewClient()
	if err != nil {
		t.Fatal("Could not create storage client:", err)
	}
}

func verifyData(t *testing.T, location storage.Location, expected string) {
	download, err := storageClient.Download(location)
	if err != nil {
		t.Errorf("Could not download object at location %s, error: %v", location, err)
	} else {
		if !bytes.Equal([]byte(expected), download) {
			t.Errorf("Download looks odd, got: %s, want: %s", string(download), expected)
		}
	}
}

func verifySubscription(t *testing.T) {
	var err error

	expected := storage.NewLocation(bucketName, "integration-executor-test/session-1/test-job-1/")
	count, err := subClient.PullMsgsSync(func(ctx ctx.Context, cancel ctx.CancelFunc, m *pubsub.Message) {
		m.Ack()
		location, err := storage.ToLocation(m.Data)
		if err != nil {
			t.Error("Could not convert message dtaa into location", err)
		} else if location != expected {
			t.Errorf("Location looks odd, got: %v, want: %v", location, expected)

		}
		cancel()
	})

	if err != nil {
		t.Errorf("Errors calling PullMsgsSync, error: %v", err)
	} else {
		if count != 1 {
			t.Errorf("Sub msg count looks odd, got: %d, want: %d", count, 1)
		}
	}
}

func teardown() {
	// delete pub/sub resources
	subClient.Delete()
	// tear down test data
	storageClient.Delete(storage.NewLocation(bucketName, business+"/"))
}

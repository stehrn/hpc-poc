package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/stehrn/hpc-poc/client"
	"github.com/stehrn/hpc-poc/executor"
)

var exe *executor.Executor

// Upload test data to cluster
// Set env variables:
//
// export GOOGLE_APPLICATION_CREDENTIALS=${HOME}/client_key.json
// export KUBE_CONFIG=${HOME}/.kube/config
//
// Example usage when running locally:
//
//     go run client_test_driver.go -business=bu3 --project=hpc-poc --bucket=hpc-poc-bucket --namespace=default --session=session-a --numTasks=2
//
func main() {

	business := flag.String("business", "", "name of business")
	project := flag.String("project", "", "name of GCP project")
	bucket := flag.String("bucket", "", "name of GCP cloud storage bucket")
	namespace := flag.String("namespace", "", "name of kubernetes namespace")
	sessionID := flag.String("session", "", "ID of session")
	numTasks := flag.Int("numTasks", 0, "number of tasks to create per job")
	flag.Parse()

	log.Printf("Following args passed in:\nbusiness: %q\nproject: %q\nbucket: %q\nnamespace: %q\nsessionID: %q\nnumTasks: %d",
		*business, *project, *bucket, *namespace, *sessionID, *numTasks)
	fmt.Println("Press the Enter Key to start")
	fmt.Scanln()

	var err error

	gcpContext := &executor.GcpContext{
		Project:    *project,
		Namespace:  *namespace,
		BucketName: *bucket,
		Business:   *business}
	exe, err = executor.New(gcpContext)
	if err != nil {
		log.Fatalf("Error creating client, %v", err)
	}

	executeJob(*numTasks, *business, *sessionID)
}

func executeJob(numTasks int, business, sessionName string) {

	session := client.NewSession(sessionName, business)
	job := client.NewJob("test-job", session)
	var n int
	for n < numTasks {
		fmt.Printf("Creating task %d", n)
		data := []byte(fmt.Sprintf("payload %d", n))
		job.CreateTask(data)
		n++
	}
	result := exe.Execute(job)
	err := result.Error

	log.Printf("Job state: %q, has errors: %v, error: %v", job.State, job.HasErrors(), err)

	if err == nil {
		log.Printf("Watching for result")
		// err = result.Watch()
	}

	if err != nil {
		log.Printf("Error: %v", err)
		cxlErr := exe.Cancel(job)
		if cxlErr != nil {
			log.Printf("Error cancelling job: %v", cxlErr)
		}
	}

	exe.Close(session)
}

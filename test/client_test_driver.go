package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/stehrn/hpc-poc/client"
	"github.com/stehrn/hpc-poc/executor"
)

var exe *executor.Executor

// Upload test data to cluster
// Set env variables:
//
// export GOOGLE_APPLICATION_CREDENTIALS=${HOME}/key.json
// export KUBE_CONFIG=${HOME}/.kube/config
//
// Example usage when running locally:
//
//     go run client_test_driver.go -business=bu1 --project=hpc-poc --bucket=hpc-poc-bucket --namespace=default --session=session-a --numJobs=3
//     go run client_test_driver.go -business=bu3 --project=hpc-poc --bucket=hpc-poc-bucket --namespace=default --session=session-a --numTasks=2
//
func main() {

	business := flag.String("business", "", "name of business")
	project := flag.String("project", "", "name of GCP project")
	bucket := flag.String("bucket", "", "name of GCP cloud storage bucket")
	namespace := flag.String("namespace", "", "name of kubernetes namespace")
	sessionID := flag.String("session", "", "ID of session")
	numJobs := flag.Int("numJobs", 0, "number of jobs to create")
	numTasks := flag.Int("numTasks", 0, "number of tasks to create per job")
	flag.Parse()

	log.Printf("Following args passed in:\nbusiness: %q\nproject: %q\nbucket: %q\nnamespace: %q\nsessionID: %q\nnumJobs: %d\nnumTasks: %d",
		*business, *project, *bucket, *namespace, *sessionID, *numJobs, *numTasks)
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

	err := result.Watch()

	if err != nil {
		log.Printf("%v", err)
		cxlErr := exe.Cancel(job)
		if cxlErr != nil {
			log.Printf("Error cancelling job: %v", cxlErr)
		}
		os.Exit(1)
	}
}

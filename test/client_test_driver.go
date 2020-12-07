package main

import (
	"flag"
	"fmt"
	"log"
	"sync"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/stehrn/hpc-poc/client"
	"github.com/stehrn/hpc-poc/kubernetes"
)

var wg sync.WaitGroup
var k8Client *kubernetes.Client

// Upload test data to cluster
// Set env variables:
// GOOGLE_APPLICATION_CREDENTIALS
// KUBE_CONFIG
//
// e.g.:
// export GOOGLE_APPLICATION_CREDENTIALS=${HOME}/key.json
// export KUBE_CONFIG=${HOME}/.kube/config
//
// Example usage when running locally:
//
// go run client_test_driver.go -business=bu1 --project=hpc-poc --bucket=hpc-poc-bucket --namespace=default --numJobs=3
//
func main() {

	business := flag.String("business", "", "name of business")
	project := flag.String("project", "", "name of GCP project")
	bucket := flag.String("bucket", "", "name of GCP cloud storage bucket")
	namespace := flag.String("namespace", "", "name of kubernetes namespace")
	numJobs := flag.Int("numJobs", 0, "number of jobs to create")
	flag.Parse()

	log.Printf("Following args passed in:\nbusiness: %q\nproject: %q\nbucket: %q\nnumJobs: %d", *business, *project, *bucket, *numJobs)
	fmt.Println("Press the Enter Key to start")
	fmt.Scanln()

	var err error
	k8Client, err = kubernetes.NewClient(*namespace)
	if err != nil {
		panic(err)
	}

	wg.Add(*numJobs)
	Business := client.Business(*business)
	client, err := client.NewClient(*project, *bucket)
	if err != nil {
		log.Fatalf("Error creating client, %v", err)
	}

	var n int
	for n < *numJobs {
		data := []byte(fmt.Sprintf("payload %d", n))
		location, id, err := client.Handle(Business, data)
		if err != nil {
			log.Fatalf("client.handle() err: %v", err)
		}
		log.Printf("Run %d, payload uploaded to cloud storage location: %s, notification sent with message ID: %s", n, location, id)
		go watch(id)
		n++
	}

	log.Print("Waiting to finish")
	wg.Wait()
	log.Print("Finished")
}

func watch(messageID string) error {
	log.Printf("Listening to subscription %q", messageID)

	options := metav1.ListOptions{LabelSelector: fmt.Sprintf("gcp.pubsub.subscription_id=%s", messageID)}
	err := k8Client.Watch(options, kubernetes.ANY, func(job *batchv1.Job) {
		podStatus, _ := k8Client.LastPodStatus(job.Name)
		log.Printf("Received update for Job %q, status: %v", job.Name, podStatus)
		if kubernetes.SUCCESS(job.Status) {
			log.Printf("Job %q succesfully finished", job.Name)
			wg.Done()
		}
	})
	return err
}

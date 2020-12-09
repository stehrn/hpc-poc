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
var myClient *client.Client

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
// go run client_test_driver.go -business=bu1 --project=hpc-poc --bucket=hpc-poc-bucket --namespace=default --session=session-a --numJobs=3
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
	k8Client, err = kubernetes.NewClient(*namespace)
	if err != nil {
		panic(err)
	}

	myClient, err = client.NewClient(*project, *bucket)
	if err != nil {
		log.Fatalf("Error creating client, %v", err)
	}

	if *numJobs != 0 {
		manyJobs(*numJobs, *business, *project, *bucket)
	} else if *numTasks != 0 {
		manyTasks(*numTasks, *business, *sessionID)
	}

	log.Print("Waiting to finish")
	wg.Wait()
	log.Print("Finished")
}

func manyJobs(numJobs int, business, project, bucket string) {
	wg.Add(numJobs)
	Business := client.Business(business)

	var n int
	for n < numJobs {
		data := []byte(fmt.Sprintf("payload %d", n))
		location, id, err := myClient.Handle(Business, data)
		if err != nil {
			log.Fatalf("client.handle() err: %v", err)
		}
		log.Printf("Run %d, payload uploaded to cloud storage location: %s, notification sent with message ID: %s", n, location, id)
		go watch(id)
		n++
	}
}

func manyTasks(numTasks int, business, session string) {
	wg.Add(1)
	Business := client.Business(fmt.Sprintf("%s/%s", business, session))

	var n int
	for n < numTasks {
		data := []byte(fmt.Sprintf("payload %d", n))
		_, err := myClient.Upload(Business, data)
		if err != nil {
			log.Fatalf("client.Upload() err: %v", err)
		}
		n++
	}

	location := myClient.Storage.LocationForObject(string(Business))
	id, err := myClient.Publish(client.Business(business), location)
	if err != nil {
		log.Fatalf("client.Publish() err: %v", err)
	}
	log.Printf("Payload uploaded to cloud storage location: %s, notification sent with message ID: %s", location, id)
	go watch(id)
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

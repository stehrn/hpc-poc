package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"cloud.google.com/go/pubsub"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/stehrn/hpc-poc/client"
	messaging "github.com/stehrn/hpc-poc/gcp/pubsub"
	"github.com/stehrn/hpc-poc/gcp/storage"
	"github.com/stehrn/hpc-poc/internal/utils"
	k8 "github.com/stehrn/hpc-poc/kubernetes"
)

// Global API clients used across function invocations.
var (
	k8Client  *k8.Client
	subClient *messaging.Client
	business  client.Business
)

func init() {
	var err error
	k8Client, err = k8.NewEnvClient()
	if err != nil {
		log.Fatalf("Could not create k8 client: %v", err)
	}
}

func init() {
	var err error
	business = client.BusinessFromEnv()
	project := utils.Env("PROJECT_NAME")
	subscription := business.SubscriptionName(project)
	subClient, err = messaging.NewSubClient(project, subscription)
	if err != nil {
		log.Fatalf("Could not create gcp sub client: %v", err)
	}
}

func main() {
	log.Print("Starting orchestrator")
	startJobWatcher()
	subscribe()
}

// TODO: this blindly tries to delete stuff already deleted when service 1st run
func startJobWatcher() {
	log.Print("Startng job watcher")
	var err error

	storageClient, err := storage.NewEnvClient()
	if err != nil {
		log.Fatalf("Could not create gcp storage client: %v", err)
	}

	options := metav1.ListOptions{LabelSelector: fmt.Sprintf("business=%s", string(business))}
	err = k8Client.Watch(options, func(job *batchv1.Job) {
		location := storage.Location{
			Bucket: job.Labels["gcp.storage.bucket"],
			Object: job.Labels["gcp.storage.object"],
		}
		log.Printf("Deleting cloud storage data at location: %v", location)
		err = storageClient.Delete(location)
		if err != nil {
			log.Printf("Failed to delete object at location: %v, error: %v", location, err)
		}
	})
	if err != nil {
		log.Fatal("Could not start watching jobs", err)
	}
}

func subscribe() {
	log.Print("Startng subscriber")
	engineImage := utils.Env("ENGINE_IMAGE")
	log.Printf("k8 job will use engine image: '%s'", engineImage)

	err := subClient.Subscribe(func(ctx context.Context, m *pubsub.Message) {
		location, err := storage.ToLocation(m.Data)
		if err != nil {
			log.Printf("Could not get location from message data (%v), error: %v", m.Data, err)
			return
		}

		options := k8.JobOptions{
			Name:     "engine-job-" + m.ID,
			Image:    engineImage,
			Labels:   labels(location, m.ID),
			Location: location}
		log.Printf("Creating Job with options: %v", options)
		_, err = k8Client.CreateJob(options)
		if err != nil {
			log.Printf("Could not create job with options: %v, error: %v", options, err)
			return
		}
		m.Ack()
	})

	if err != nil {
		panic(err)
	}
}

func labels(location storage.Location, messageID string) map[string]string {
	labels := make(map[string]string)
	labels["business"] = string(business)
	labels["k8.namespace"] = k8Client.Namespace
	labels["gcp.storage.bucket"] = location.Bucket
	labels["gcp.storage.object"] = clean(location.Object)
	labels["gcp.pubsub.project"] = subClient.Project
	labels["gcp.pubsub.subscription"] = subClient.SubscriptionID
	labels["gcp.pubsub.subscription_id"] = messageID
	return labels
}

func clean(item string) string {
	return strings.ReplaceAll(item, "/", "_")
}

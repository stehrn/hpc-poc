package main

import (
	"context"
	"log"
	"strings"

	"cloud.google.com/go/pubsub"
	batchv1 "k8s.io/api/batch/v1"

	"github.com/stehrn/hpc-poc/client"
	messaging "github.com/stehrn/hpc-poc/gcp/pubsub"
	"github.com/stehrn/hpc-poc/gcp/storage"
	"github.com/stehrn/hpc-poc/internal/utils"
	k8 "github.com/stehrn/hpc-poc/kubernetes"
)

func main() {
	log.Print("Starting orchestrator")

	k8Client, err := k8.NewEnvClient()
	if err != nil {
		log.Fatalf("Could not create k8 client: %v", err)
	}

	project := utils.Env("PROJECT_NAME")
	business := client.BusinessFromEnv()
	subscription := business.SubscriptionName(project)
	subClient, err := messaging.NewSubClient(project, subscription)
	if err != nil {
		log.Fatalf("Could not create gcp sub client: %v", err)
	}
	storageClient, err := storage.NewEnvClient()
	if err != nil {
		log.Fatalf("Could not create gcp storage client: %v", err)
	}

	engineImage := utils.Env("ENGINE_IMAGE")
	log.Printf("k8 job will use engine image: '%s'", engineImage)

	// TODO: this blindly tries to delete stuff already deleted when service 1st run
	log.Print("Startng job watcher")
	err = k8Client.Watch(func(job *batchv1.Job) {
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

	err = subClient.Subscribe(func(ctx context.Context, m *pubsub.Message) {
		location, err := storage.ToLocation(m.Data)
		if err != nil {
			log.Printf("Could not get location from message data (%v), error: %v", m.Data, err)
			return
		}

		options := k8.JobOptions{
			Name:     "engine-job-" + m.ID,
			Image:    engineImage,
			Labels:   labels(business, k8Client, subClient, location, m.ID),
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

func labels(business client.Business, k8Client *k8.Client, pubsubClient *messaging.Client, location storage.Location, messageID string) map[string]string {
	labels := make(map[string]string)
	labels["business"] = string(business)
	labels["k8.namespace"] = k8Client.Namespace
	labels["gcp.storage.bucket"] = location.Bucket
	labels["gcp.storage.object"] = clean(location.Object)
	labels["gcp.pubsub.project"] = pubsubClient.Project
	labels["gcp.pubsub.subscription"] = pubsubClient.SubscriptionID
	labels["gcp.pubsub.subscription_id"] = messageID
	return labels
}

func clean(item string) string {
	return strings.ReplaceAll(item, "/", "_")
}

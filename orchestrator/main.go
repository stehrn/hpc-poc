package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	gcp "github.com/stehrn/hpc-poc/gcp"
	k8 "github.com/stehrn/hpc-poc/kubernetes"
)

func main() {
	log.Print("Starting orchestrator")

	k8Client := k8Client()
	engineImage := env("ENGINE_IMAGE")
	log.Printf("k8 job will use engine image: %s", engineImage)
	gcpClient := gcpClient()

	err := gcpClient.Subscribe(func(ctx context.Context, m *pubsub.Message) {
		jobName := "engine-job-" + m.ID
		payload := string(m.Data)
		log.Printf("Got message: %s, creating Job: %s", payload, jobName)
		k8Client.CreateJob(k8.JobCreate{Name: jobName, Image: engineImage, PayLoad: payload})
		m.Ack()
	})

	if err != nil {
		panic(err)
	}
}

func k8Client() k8.Client {
	// k8 client
	namespace := env("NAMSPACE")
	log.Printf("Creating k8 jobs client for namespace: %s", namespace)
	return k8.NewClient(namespace)
}

func gcpClient() gcp.Client {
	project := "hpc-poc"
	subscriptionID := env("SUBSCRIPTION_NAME")
	topic := ""
	log.Printf("Creating gcp client for project: %s, subscriptionID: %s, topic: %s", project, subscriptionID, topic)
	return gcp.NewClient(project, subscriptionID, topic)
}

func env(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("No '%s' env variable", key)
	}
	return value
}

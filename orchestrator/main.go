package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	k8 "github.com/stehrn/hpc-poc/kubernetes"
)

func main() {
	log.Print("Starting orchestrator")

	namespace := env("NAMSPACE")
	engineImage := env("ENGINE_IMAGE")

	log.Printf("Creating jobs client for namespace %s (job will use image: %s)", namespace, engineImage)
	jobService := k8.New(namespace)

	subscribe(func(ctx context.Context, m *pubsub.Message) {
		jobName := "engine-job-" + m.ID
		payload := string(m.Data)
		log.Printf("Got message: %s, creating Job: %s", payload, jobName)
		jobService.CreateJob(JobCreate{jobName, engineImage, payload})
		m.Ack()
	})
}

func subscribe(callback func(ctx context.Context, m *pubsub.Message)) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, "hpc-poc")
	if err != nil {
		log.Fatalf("Could not create client: %v", err)
	}

	subName := env("SUBSCRIPTION_NAME")
	log.Printf("Subscribing to %s", subName)
	sub := client.Subscription(subName)
	err = sub.Receive(context.Background(), callback)
	if err != nil {
		log.Fatalf("Could not receive message: %v", err)
	}
}

func env(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("No '%s' env variable", key)
	}
	return value
}

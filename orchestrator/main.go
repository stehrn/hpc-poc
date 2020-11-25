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

	// k8 settings
	namespace := env("NAMSPACE")
	engineImage := env("ENGINE_IMAGE")
	// gcp settings
	project := "hpc-poc"
	subscriptionID := env("SUBSCRIPTION_NAME")

	log.Printf("Creating jobs client for namespace %s (job will use image: %s)", namespace, engineImage)
	client := k8.NewClient(namespace)

	err := gcp.Subscribe(project, subscriptionID, func(ctx context.Context, m *pubsub.Message) {
		jobName := "engine-job-" + m.ID
		payload := string(m.Data)
		log.Printf("Got message: %s, creating Job: %s", payload, jobName)
		client.CreateJob(k8.JobCreate{Name: jobName, Image: engineImage, PayLoad: payload})
		m.Ack()
	})

	if err != nil {
		panic(err)
	}
}

func env(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("No '%s' env variable", key)
	}
	return value
}

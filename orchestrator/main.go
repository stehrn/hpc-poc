package main

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
	k8 "github.com/stehrn/hpc-poc/kubernetes"
)

func main() {
	log.Print("Starting orchestrator")

	k8Client := k8Client()
	engineImage := utils.Env("ENGINE_IMAGE")
	log.Printf("k8 job will use engine image: %s", engineImage)
	gcpClient, err := gcpClient()
	if err != nil {
		log.Fatalf("Could not create gcp client: %v", err)
	}

	err = gcpClient.Subscribe(func(ctx context.Context, m *pubsub.Message) {
		jobName := "engine-job-" + m.ID
		payload := string(m.Data)
		log.Printf("Got message: %s, creating Job: %s", payload, jobName)
		k8Client.CreateJob(k8.JobInfo{Name: jobName, Image: engineImage, PayLoad: payload})
		m.Ack()
	})

	if err != nil {
		panic(err)
	}
}

func k8Client() *k8.Client {
	namespace := utils.Env("NAMSPACE")
	log.Printf("Creating k8 jobs client for namespace: %s", namespace)
	return k8.NewClient(namespace)
}

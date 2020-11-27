package main

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
	gcp "github.com/stehrn/hpc-poc/gcp"
	"github.com/stehrn/hpc-poc/internal/utils"
	k8 "github.com/stehrn/hpc-poc/kubernetes"
)

func main() {
	log.Print("Starting orchestrator")

	k8Client, err := k8.NewClientFromEnvironment()
	if err != nil {
		log.Fatalf("Could not create k8 client: %v", err)
	}
	gcpClient, err := gcp.NewClientFromEnvironment()
	if err != nil {
		log.Fatalf("Could not create gcp client: %v", err)
	}

	engineImage := utils.Env("ENGINE_IMAGE")
	log.Printf("k8 job will use engine image: %s", engineImage)

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

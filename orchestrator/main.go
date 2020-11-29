package main

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
	messaging "github.com/stehrn/hpc-poc/gcp/pubsub"
	"github.com/stehrn/hpc-poc/gcp/storage"
	"github.com/stehrn/hpc-poc/internal/utils"
	k8 "github.com/stehrn/hpc-poc/kubernetes"
)

func main() {
	log.Print("Starting orchestrator")

	k8Client, err := k8.NewClient()
	if err != nil {
		log.Fatalf("Could not create k8 client: %v", err)
	}
	pubsubClient, err := messaging.NewClient()
	if err != nil {
		log.Fatalf("Could not create gcp pubsub client: %v", err)
	}

	engineImage := utils.Env("ENGINE_IMAGE")
	log.Printf("k8 job will use engine image: %s", engineImage)

	err = pubsubClient.Subscribe(func(ctx context.Context, m *pubsub.Message) {
		jobName := "engine-job-" + m.ID
		location, err := storage.ToLocation(m.Data)
		if err != nil {
			log.Fatalf("Could not get location from message data: %v", err)
		}
		log.Printf("Creating Job %s to processes data at %s", jobName, location)
		k8Client.CreateJob(k8.JobInfo{jobName, engineImage, location})
		m.Ack()
	})

	if err != nil {
		panic(err)
	}
}

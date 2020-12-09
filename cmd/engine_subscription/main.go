package main

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
	messaging "github.com/stehrn/hpc-poc/gcp/pubsub"
	"github.com/stehrn/hpc-poc/gcp/storage"
	"github.com/stehrn/hpc-poc/internal/utils"
)

// Global API clients used across function invocations.
var (
	subClient     *messaging.Client
	storageClient *storage.Client
)

// init subcription client, requires: PROJECT_NAME, SUBSCRIPTION_NAME
func init() {
	var err error
	project := utils.Env("PROJECT_NAME")
	subscription := utils.Env("SUBSCRIPTION_NAME")
	subClient, err = messaging.NewSubClient(project, subscription)
	if err != nil {
		log.Fatalf("Could not create gcp sub client: %v", err)
	}
}

// init storage client, requires: CLOUD_STORAGE_BUCKET_NAME
func init() {
	var err error
	storageClient, err = storage.NewEnvClient()
	if err != nil {
		log.Fatalf("Could not create storage client: %v", err)
	}
}

// Simple engine that subscribes to topic and loads locations
func main() {
	log.Print("Starting Engine")

	for {
		count, err := subClient.PullMsgsSync(func(ctx context.Context, m *pubsub.Message) {
			location, err := storage.ToLocation(m.Data)
			if err != nil {
				log.Printf("Could not get location from message data (%v), error: %v", m.Data, err)
				return
			}

			data, err := storageClient.Download(location)
			if err != nil {
				log.Fatalf("Failed to download object, error: %v", err)
			}

			log.Printf("Loaded data: %v", string(data))

			// to simulate engine failure
			if string(data) == "PANIC" {
				panic("engine failed!")
			}

			m.Ack()
		})

		if err != nil {
			panic(err)
		}

		if count == 0 {
			log.Print("No messages left, exiting")
			break
		} else {
			log.Printf("Processed %d messages", count)
		}
	}

	log.Print("Exit")
}

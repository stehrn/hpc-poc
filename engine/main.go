package main

import (
	"log"

	"github.com/stehrn/hpc-poc/gcp/storage"
)

// Simple engine that downloads a cloud storage object based on info in env, prints contents out, and exits
func main() {
	log.Print("Starting Engine")

	client, err := storage.NewClient()
	if err != nil {
		log.Fatalf("Could not create storage client: %v", err)
	}

	location := client.LocationFromEnvironment()
	log.Printf("Loading data from cloud storage location (bucket/object) %v", location)
	data, err := client.Download(location)
	if err != nil {
		log.Fatalf("Failed to download object, error: %v", err)
	}

	log.Printf("Loaded data: %v", string(data))
}

package main

import (
	"fmt"
	"log"

	"github.com/stehrn/hpc-poc/gcp/storage"
)

// Simple engine that download a cloud storage object based on info in env, prints contents out, and exits
func main() {
	log.Print("Starting Engine")

	location := storage.LocationFromEnvironment()
	log.Printf("Loading data from cloud storage %v", location)
	data, err := storage.Download(location)
	if err != nil {
		log.Fatalf("Failed to download object, error: %v", err)
	}

	fmt.Printf("Loaded data: %v", data)

	log.Printf("Deleting cloud storage data at %v", location)
	err = storage.Delete(location)
	if err != nil {
		log.Fatalf("Failed to delete object, error: %v", err)
	}
}

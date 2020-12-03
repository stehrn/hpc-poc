package main

import (
	"fmt"
	"log"
	"os"

	"github.com/stehrn/hpc-poc/client"
	"github.com/stehrn/hpc-poc/gcp/pubsub"
)

// to run, the following env variables need to be set:
// GOOGLE_APPLICATION_CREDENTIALS (see main README and 'Get GCP JSON key...')
// BUCKET_NAME
// TOPIC_NAME
func main() {
	log.Print("Starting test driver")
	bucket := os.Getenv("BUCKET")
	client, err := client.NewClient(bucket, pubsub.ConfigFromEnvironment())
	if err != nil {
		log.Fatalf("Could not create client: %v", err)
	}

	var n int
	for n < 5 {
		data := []byte(fmt.Sprintf("payload %d", n))
		location, id, err := client.Handle(data)
		if err != nil {
			log.Fatalf("client.handle() err: %v", err)
		}

		log.Printf("Payload uploaded to cloud storage location: %s, notification sent with message ID: %s", location, id)
	}
}

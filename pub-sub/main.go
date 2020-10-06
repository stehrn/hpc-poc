package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
)

func main() {
	log.Print("Starting pub/sub")

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, "hpc-poc")
	if err != nil {
		log.Fatalf("Could not create pubsub Client: %v", err)
	}

	subName := os.Getenv("SUBSCRIPTION_NAME")
	if subName == "" {
		log.Fatal("No 'SUBSCRIPTION_NAME' env variable")
	}
	log.Printf("Subscribing to %s", subName)
	sub := client.Subscription(subName)

	err = sub.Receive(context.Background(), func(ctx context.Context, m *pubsub.Message) {
		log.Printf("Got message: %s", m.Data)
		m.Ack()
	})
	if err != nil {
		log.Fatalf("Could not receive message: %v", err)
	}
}

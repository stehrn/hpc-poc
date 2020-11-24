package gcp

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
)

// Subscribe subscribes to givne project/id, passing message into callpack
func Subscribe(project string, subscriptionID string, callback func(ctx context.Context, m *pubsub.Message)) error {
	log.Printf("Subscribing to project: %s, subscriptionID: %s", project, subscriptionID)

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, project)
	if err != nil {
		return errors.Wrapf(err, "Failed to create pubsub client (project %s)", project)
	}
	

	// topic, err := client.
	// res := topic.Publish(ctx, &pubsub.Message{Data: []byte("payload")})

	sub := client.Subscription(subscriptionID)
	err = sub.Receive(context.Background(), callback)
	if err != nil {
		return errors.Wrap(err, "Could not receive message")
	}

}

func Publish() {

	

}

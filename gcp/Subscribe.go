package gcp

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
)

// Subscribe subscribe to given project/id, passing message into callpack
func (c Client) Subscribe(callback func(ctx context.Context, m *pubsub.Message)) error {
	if c.info.Subscription == "" {
		return errors.New("Subscription required")
	}
	log.Printf("Subscribing to project: %s, subscriptionID: %s", c.info.Project, c.info.Subscription)
	sub := c.client.Subscription(c.info.Subscription)
	err := sub.Receive(context.Background(), callback)
	if err != nil {
		return errors.Wrap(err, "Could not receive message")
	}
	return nil
}

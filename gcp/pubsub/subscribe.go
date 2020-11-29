package pubsub

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
)

// Subscribe subscribe to given project/id, passing message into callpack
func (c Client) Subscribe(callback func(ctx context.Context, m *pubsub.Message)) error {
	if c.Subscription == "" {
		return errors.New("Subscription required")
	}
	log.Printf("Subscribing to project: %s, subscriptionID: %s", c.Project, c.Subscription)
	sub := c.client.Subscription(c.Subscription)
	err := sub.Receive(context.Background(), callback)
	if err != nil {
		return errors.Wrap(err, "Could not receive message")
	}
	return nil
}

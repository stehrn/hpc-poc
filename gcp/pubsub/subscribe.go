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
	log.Printf("Subscribing to project: '%s', subscriptionID: '%s'", c.Project, c.Subscription)
	sub := c.client.Subscription(c.Subscription)
	ctx := context.Background()
	ok, err := sub.Exists(ctx)
	if err != nil {
		return errors.Wrapf(err, "Failed to find out if subscription '%s' exists", c.Subscription)
	}
	if !ok {
		return errors.Errorf("Subscription '%s' does not exist", c.Subscription)
	}

	err = sub.Receive(context.Background(), callback)
	if err == context.Canceled {
		log.Print("Subscription cancelled")
		return nil

	}
	return errors.Wrap(err, "Error recieving message")
}

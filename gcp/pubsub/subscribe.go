package pubsub

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
)

// Subscribe subscribe to given project/id, passing message into callpack
func (c Client) Subscribe(callback func(ctx context.Context, m *pubsub.Message)) error {
	if c.SubscriptionID == "" {
		return errors.New("Subscription required")
	}
	log.Printf("Subscribing to project: '%s', subscriptionID: '%s'", c.Project, c.SubscriptionID)
	sub, err := c.subscription()
	if err != nil {
		return err
	}

	err = sub.Receive(context.Background(), callback)
	if err == context.Canceled {
		log.Print("Subscription cancelled")
		return nil

	}
	return errors.Wrap(err, "Error recieving message")
}

func (c Client) subscription() (*pubsub.Subscription, error) {
	sub := c.Subscription(c.SubscriptionID)
	ctx := context.Background()
	ok, err := sub.Exists(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to find out if subscription '%s' exists", c.SubscriptionID)
	}
	if !ok {
		return nil, errors.Errorf("Subscription '%s' does not exist", c.SubscriptionID)
	}
	return sub, nil
}

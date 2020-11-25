package gcp

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	ps "cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
)

// Client client
type Client struct {
	project        string
	subscriptionID string
	topic          string
	client         *ps.Client
}

// NewClient create Client
func NewClient(project, subscriptionID, topic string) (*Client, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, project)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create pubsub client (project %s)", project)
	}
	return &Client{project, subscriptionID, topic, client}, nil
}

// Subscribe subscribes to given project/id, passing message into callpack
func (c Client) Subscribe(callback func(ctx context.Context, m *pubsub.Message)) error {
	log.Printf("Subscribing to project: %s, subscriptionID: %s", c.project, c.subscriptionID)
	sub := c.client.Subscription(c.subscriptionID)
	err := sub.Receive(context.Background(), callback)
	if err != nil {
		return errors.Wrap(err, "Could not receive message")
	}
	return nil
}

// Publish publish to topic
func (c Client) Publish(payload []byte) error {
	topic := c.client.Topic(c.topic)
	ctx := context.Background()
	res := topic.Publish(ctx, &pubsub.Message{Data: payload})
	id, err := res.Get(ctx)
	if err != nil {
		return errors.Wrap(err, "Could not publish message")
	}
	fmt.Printf("Published message with a message ID: %s", id)
	return nil
}

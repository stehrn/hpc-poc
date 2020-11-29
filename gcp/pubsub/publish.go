package pubsub

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
)

// Publish publish to topic
func (c Client) Publish(payload []byte) (string, error) {
	if c.Topic == "" {
		return "", errors.New("Topic required")
	}
	topic := c.client.Topic(c.Topic)
	ctx := context.Background()
	res := topic.Publish(ctx, &pubsub.Message{Data: payload})
	id, err := res.Get(ctx)
	if err != nil {
		return "", errors.Wrap(err, "Could not publish message")
	}
	fmt.Printf("Published message with a message ID: %s", id)
	return id, nil
}

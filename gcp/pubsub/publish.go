package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
)

// Publish publish to topic
func (c Client) Publish(payload []byte) (string, error) {
	if c.TopicName == "" {
		return "", errors.New("Topic required")
	}
	ctx := context.Background()
	topic := c.Topic(c.TopicName)
	ok, err := topic.Exists(ctx)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to find out if topic %s exists", c.Topic)
	}
	if !ok {
		return "", errors.Errorf("Topic %s does not exist", c.Topic)
	}

	defer topic.Stop()
	res := topic.Publish(ctx, &pubsub.Message{Data: payload})
	id, err := res.Get(ctx)
	if err != nil {
		return "", errors.Wrap(err, "Could not publish message")
	}
	return id, nil
}

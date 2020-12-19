package pubsub

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
)

// Publish publish to topic
func (c Client) Publish(topicName string, payload []byte) (string, error) {
	topic, err := c.topic(topicName)
	if err != nil {
		return "", err
	}

	ctx := context.Background()

	defer topic.Stop()
	res := topic.Publish(ctx, &pubsub.Message{Data: payload})
	id, err := res.Get(ctx)
	if err != nil {
		return "", errors.Wrap(err, "Could not publish message")
	}
	return id, nil
}

// PublishMany publish many payloads
func (c Client) PublishMany(topicName string, payloads [][]byte) error {
	topic, err := c.topic(topicName)
	if err != nil {
		return err
	}

	ctx := context.Background()

	var wg sync.WaitGroup
	var totalErrors uint64

	for _, payload := range payloads {
		res := topic.Publish(ctx, &pubsub.Message{Data: payload})
		wg.Add(1)
		go func(res *pubsub.PublishResult) {
			defer wg.Done()
			// Get() blocks until a server-generated ID or an error is returned for the published message.
			_, err := res.Get(ctx)
			if err != nil {
				fmt.Printf("Failed to publish: %s\n", err)
				atomic.AddUint64(&totalErrors, 1)
				return
			}
		}(res)
	}

	wg.Wait()

	if totalErrors > 0 {
		return fmt.Errorf("%d of %d messages did not publish successfully", totalErrors, len(payloads))
	}

	return nil
}

func (c Client) topic(topicName string) (*pubsub.Topic, error) {
	topic := c.Topic(topicName)
	ctx := context.Background()
	ok, err := topic.Exists(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to find out if topic %q exists", topicName)
	}
	if !ok {
		return nil, errors.Errorf("Topic %q does not exist", topicName)
	}
	return topic, nil
}

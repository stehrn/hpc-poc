package pubsub

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	if err != nil && status.Code(err) != codes.Canceled {
		return fmt.Errorf("Receive: %v", err)
	}
	return nil
}

// PullMsgsSync
func (c Client) PullMsgsSync(callback func(ctx context.Context, m *pubsub.Message)) (int32, error) {
	if c.SubscriptionID == "" {
		return 0, errors.New("Subscription required")
	}
	log.Printf("Subscribing to project: '%s', subscriptionID: '%s'", c.Project, c.SubscriptionID)
	sub, err := c.subscription()
	if err != nil {
		return 0, err
	}
	sub.ReceiveSettings.Synchronous = true
	sub.ReceiveSettings.MaxOutstandingMessages = 1

	ctx := context.Background()

	// Receive messages for 10 seconds.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var counter int32
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		callback(ctx, msg)
		atomic.AddInt32(&counter, 1)
	})
	if err != nil && status.Code(err) != codes.Canceled {
		return counter, fmt.Errorf("Receive: %v", err)
	}
	return counter, nil
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

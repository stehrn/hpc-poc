package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
)

// ClientConfg pubsub client configuration
type ClientConfg struct {
	Project        string
	SubscriptionID string
	TopicName      string
}

// Client pubsub client
type Client struct {
	*ClientConfg
	*pubsub.Client
}

// NewPubClient create new client for subscriptions
func NewPubClient(project, topic string) (*Client, error) {
	if project == "" {
		return nil, errors.New("Missing project")
	}
	if topic == "" {
		return nil, errors.New("Missing topic")
	}
	return create(&ClientConfg{
		Project:   project,
		TopicName: topic})
}

// NewSubClient create new client for subscriptions
func NewSubClient(project, subscription string) (*Client, error) {
	if project == "" {
		return nil, errors.New("Missing project")
	}
	if subscription == "" {
		return nil, errors.New("Missing subscription")
	}
	return create(&ClientConfg{
		Project:        project,
		SubscriptionID: subscription})
}

func create(config *ClientConfg) (*Client, error) {
	ctx := context.Background()
	pubsubClient, err := pubsub.NewClient(ctx, config.Project)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create pubsub client (project %s)", config.Project)
	}
	return &Client{config, pubsubClient}, nil
}

package pubsub

import (
	"context"
	"os"

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

// NewClientFromEnvironment create Client
func NewClientFromEnvironment() (*Client, error) {
	return create(ConfigFromEnvironment())
}

// NewPubClient create new client for subscriptions
func NewPubClient(config *ClientConfg) (*Client, error) {
	if config.Project == "" {
		return nil, errors.New("Missing project (PROJECT_NAME env variable)")
	}
	if config.TopicName == "" {
		return nil, errors.New("Missing topics (TOPIC_NAME env variable)")
	}
	return create(config)
}

// NewSubClient create new client for subscriptions
func NewSubClient(config *ClientConfg) (*Client, error) {
	if config.TopicName == "" {
		return nil, errors.New("Missing topics (TOPIC_NAME env variable)")
	}
	return create(config)
}

func create(config *ClientConfg) (*Client, error) {
	ctx := context.Background()
	pubsubClient, err := pubsub.NewClient(ctx, config.Project)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create pubsub client (project %s)", config.Project)
	}
	return &Client{config, pubsubClient}, nil
}

// ConfigFromEnvironment create info from environment variables: PROJECT_NAME, SUBSCRIPTION_NAME, TOPIC_NAME
func ConfigFromEnvironment() *ClientConfg {
	return &ClientConfg{
		Project:        os.Getenv("PROJECT_NAME"),
		SubscriptionID: os.Getenv("SUBSCRIPTION_NAME"),
		TopicName:      os.Getenv("TOPIC_NAME")}
}

package pubsub

import (
	"context"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
)

// Client client
type Client struct {
	Project      string
	Subscription string
	Topic        string
	client       *pubsub.Client
}

// NewClient create Client
func NewClient() (*Client, error) {
	info := clientFromEnvironment()
	if info.Project == "" {
		return nil, errors.New("Project required")
	}
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, info.Project)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create pubsub client (project %s)", info.Project)
	}
	info.client = client
	return &info, nil
}

// create info from environment variables: PROJECT_NAME, SUBSCRIPTION_NAME, TOPIC_NAME
func clientFromEnvironment() Client {
	return Client{
		Project:      os.Getenv("PROJECT_NAME"),
		Subscription: os.Getenv("SUBSCRIPTION_NAME"),
		Topic:        os.Getenv("TOPIC_NAME")}
}

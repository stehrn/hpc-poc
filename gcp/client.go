package gcp

import (
	"context"
	"os"

	"cloud.google.com/go/pubsub"
	ps "cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
)

// Client client
type Client struct {
	info   Info
	client *ps.Client
}

// Info information about pub/sub env
type Info struct {
	Project      string
	Subscription string
	Topic        string
}

// NewClientFromEnvironment create Client
func NewClientFromEnvironment() (*Client, error) {
	return NewClient(InfoFromEnvironment())
}

// NewClient create Client
func NewClient(info Info) (*Client, error) {
	if info.Project == "" {
		return nil, errors.New("Project required")
	}
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, info.Project)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create pubsub client (project %s)", info.Project)
	}
	return &Client{info, client}, nil
}

// InfoFromEnvironment create GcpInfoFromEnvironment from environment variables: PROJECT_NAME, SUBSCRIPTION_NAME, TOPIC_NAME
func InfoFromEnvironment() Info {
	return Info{
		Project:      os.Getenv("PROJECT_NAME"),
		Subscription: os.Getenv("SUBSCRIPTION_NAME"),
		Topic:        os.Getenv("TOPIC_NAME")}
}

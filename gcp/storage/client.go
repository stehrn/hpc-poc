package storage

import (
	"context"
	"fmt"

	gcp "cloud.google.com/go/storage"
)

// Client storage client
type Client struct {
	*gcp.Client
}

// NewClient create new storage client
func NewClient() (*Client, error) {
	ctx := context.Background()
	client, err := gcp.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	return &Client{client}, nil
}

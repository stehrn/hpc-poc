package storage

import (
	"context"
	"errors"
	"fmt"
	"os"

	gcp "cloud.google.com/go/storage"
)

// Client storage client
type Client struct {
	BucketName string
	*gcp.Client
}

// NewEnvClient create new storage client from env
func NewEnvClient() (*Client, error) {
	bucketName := os.Getenv("CLOUD_STORAGE_BUCKET_NAME")
	if bucketName == "" {
		return &Client{}, errors.New("storage.NewEnvClient: env CLOUD_STORAGE_BUCKET_NAME blank")
	}
	return NewClient(bucketName)
}

// NewClient create new storage client
func NewClient(bucket string) (*Client, error) {
	if bucket == "" {
		return nil, errors.New("storage.NewClient: bucket cannot be blank")
	}
	ctx := context.Background()
	gcpClient, err := gcp.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	return &Client{bucket, gcpClient}, nil
}

package storage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/storage"
	gcp "cloud.google.com/go/storage"
)

// ClientInterface defines API client methods for storage
type ClientInterface interface {
	BucketName() string
	ListObjects(prefix string) ([]Object, error)
	ForEachObject(prefix string, consumer func(attrs *storage.ObjectAttrs)) error
	Upload(location Location, content []byte) error
	Download(location Location) ([]byte, error)
	Delete(location Location) error
	LocationClient
}

// LocationClient location specific
type LocationClient interface {
	Location(business string) Location
	LocationForObject(object string) Location
	ToLocationByteSlice(objects []Object) ([][]byte, error)
}

// Object compact representation of a storage object
type Object struct {
	Object  string
	Size    int64
	Created time.Time
}

// Client storage client implements ClientInterface
type Client struct {
	bucketName string
	*gcp.Client
}

// BucketName name of bucket
func (c *Client) BucketName() string {
	return c.bucketName
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

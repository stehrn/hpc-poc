package storage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/storage"
	gcp "cloud.google.com/go/storage"
	"github.com/stehrn/hpc-poc/client"
)

// ClientInterface defines API client methods for storage
type ClientInterface interface {
	// BucketName name of bucket client interacts with
	BucketName() string
	// ListObjects list objects at given prefix
	ListObjects(prefix string) ([]Object, error)
	// ForEachObject iteratore over objects at given prefix, passing result(s) to lamda function
	ForEachObject(prefix string, consumer func(attrs *storage.ObjectAttrs) error) error
	// Upload upload data to given location
	Upload(location Location, content []byte) error
	// UploadMany data upload if many items to upload
	UploadMany(items client.DataSourceIterator) uint64
	// Download download object at given location
	Download(location Location) ([]byte, error)
	// Delete delete object at given location
	Delete(location Location) error
	LocationClient
}

// LocationClient location specific
type LocationClient interface {
	Location(business string) Location
	LocationForObject(object string) Location
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
func NewClient(bucketName string) (*Client, error) {
	if bucketName == "" {
		return nil, errors.New("storage.NewClient: bucket cannot be blank")
	}
	ctx := context.Background()
	gcpClient, err := gcp.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	bucket := gcpClient.Bucket(bucketName)
	_, err = bucket.Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient, bucket %s does not exist: %v", bucketName, err)
	}
	return &Client{bucketName, gcpClient}, nil
}

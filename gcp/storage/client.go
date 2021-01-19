package storage

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/storage"
	gcp "cloud.google.com/go/storage"
	"github.com/stehrn/hpc-poc/client"
)

// ClientInterface defines API client methods for storage
type ClientInterface interface {
	// ListObjects list objects at given prefix, just returns a lighweight representation of a storage object
	ListObjects(location Location) ([]Object, error)
	// ForEachObject iteratore over objects at given prefix, passing result(s) to lamda function
	// Use if you need access to all possible object atttributes
	ForEachObject(location Location, consumer func(attrs *storage.ObjectAttrs) error) error
	// Upload upload data to given location
	Upload(location Location, content []byte) error
	// UploadMany data upload if many items to upload
	UploadMany(bucketName string, items client.DataSourceIterator) uint64
	// Download download object at given location
	Download(location Location) ([]byte, error)
	// Delete delete object at given location
	Delete(location Location) error
	// BucketExists check if bucket exists
	BucketExists(bucketName string) (bool, error)
}

// Object compact representation of a storage object
type Object struct {
	Object  string
	Size    int64
	Created time.Time
}

// Client storage client implements ClientInterface
type Client struct {
	*gcp.Client
}

// NewClient create new storage client
func NewClient() (*Client, error) {
	ctx := context.Background()
	gcpClient, err := gcp.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	return &Client{gcpClient}, nil
}

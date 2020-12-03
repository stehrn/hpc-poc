package client

import (
	"log"

	"github.com/rs/xid"
	"github.com/stehrn/hpc-poc/gcp/pubsub"
	"github.com/stehrn/hpc-poc/gcp/storage"
)

// Client client
type Client struct {
	Bucket  string
	Pubsub  *pubsub.Client
	Storage *storage.Client
}

// NewClient create Client
func NewClient(bucket string, config *pubsub.ClientConfg) (*Client, error) {
	pubsubClient, err := pubsub.NewPubClient(config)
	if err != nil {
		log.Fatalf("Could not create pubsub client: %v", err)
	}
	storageClient, err := storage.NewClient()
	if err != nil {
		log.Fatalf("Could not create storage client: %v", err)
	}
	return &Client{bucket, pubsubClient, storageClient}, nil
}

func (c Client) location() storage.Location {
	return storage.Location{
		Bucket: c.Bucket,
		Object: xid.New().String()}
}

// Handle handle data
func (c Client) Handle(data []byte) (storage.Location, string, error) {
	// upload payload to cloud storage
	location := c.location()
	log.Printf("Uploading data to %v", location)
	err := c.Storage.Upload(location, data)
	if err != nil {
		return wrap(err)
	}

	// publish object location
	log.Printf("Publishing location (%v) to topic %s", location, c.Pubsub.TopicName)
	bytes, err := storage.ToBytes(location)
	if err != nil {
		return wrap(err)
	}
	id, err := c.Pubsub.Publish(bytes)
	if err != nil {
		return wrap(err)
	}
	return location, id, nil
}

func wrap(err error) (storage.Location, string, error) {
	return storage.Location{}, "", err
}

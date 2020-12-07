package client

import (
	"log"

	"github.com/pkg/errors"

	"github.com/stehrn/hpc-poc/gcp/pubsub"
	"github.com/stehrn/hpc-poc/gcp/storage"
	"github.com/stehrn/hpc-poc/internal/utils"
)

var pubsubClients = make(map[Business]*pubsub.Client)

// Client client
type Client struct {
	Project string
	Storage *storage.Client
}

// NewEnvClientOrFatal create new client from env
func NewEnvClientOrFatal() *Client {
	project := utils.Env("PROJECT_NAME")
	bucket := utils.Env("CLOUD_STORAGE_BUCKET_NAME")
	client, err := NewClient(project, bucket)
	if err != nil {
		log.Fatalf("Could not create client: %v", err)
	}
	return client
}

// NewClient create Client
func NewClient(project, bucket string) (*Client, error) {
	if project == "" {
		return nil, errors.New("Missing project")
	}
	if bucket == "" {
		return nil, errors.New("Missing bucket")
	}
	storageClient, err := storage.NewClient(bucket)
	if err != nil {
		return nil, err
	}
	return &Client{project, storageClient}, nil
}

// Handle handle data for given business:
//  1) upload payload to cloud storage
//  2) publish object location
func (c Client) Handle(business Business, data []byte) (storage.Location, string, error) {
	location, err := c.upload(business, data)
	if err != nil {
		return wrap(err)
	}
	id, err := c.publish(business, location)
	if err != nil {
		log.Printf("Publish failed, deleting data at '%v'", location)
		c.delete(location)
		return wrap(err)
	}
	return location, id, nil
}

// Topic get name of topic for given business
func (c Client) Topic(business Business) string {
	return business.TopicName(c.Project)
}

// upload upload data to cloud storage
func (c Client) upload(business Business, data []byte) (storage.Location, error) {
	location := c.Storage.Location(string(business))
	log.Printf("Uploading data to: '%v'", location)
	return location, c.Storage.Upload(location, data)
}

func (c Client) delete(location storage.Location) {
	err := c.Storage.Delete(location)
	if err != nil {
		log.Printf("Failed to delete object at: '%v', error: %v", location, err)
	}
}

// publish location to topic derived off business
func (c Client) publish(business Business, location storage.Location) (string, error) {
	pubsubClient, err := c.pubsubClient(business)
	if err != nil {
		return "", err
	}
	log.Printf("Publishing location (%v) to topic: '%s'", location, pubsubClient.TopicName)
	bytes, err := storage.ToBytes(location)
	if err != nil {
		return "", err
	}
	id, err := pubsubClient.Publish(bytes)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (c Client) pubsubClient(business Business) (*pubsub.Client, error) {
	if pubsubClient, found := pubsubClients[business]; found {
		return pubsubClient, nil
	}
	topic := business.TopicName(c.Project)
	pubsubClient, err := pubsub.NewPubClient(c.Project, topic)
	if err != nil {
		return &pubsub.Client{}, err
	}
	log.Printf("Created new pubsubClient for business '%s' with topic '%s'", business, pubsubClient.TopicName)
	pubsubClients[business] = pubsubClient
	return pubsubClient, nil
}

func wrap(err error) (storage.Location, string, error) {
	return storage.Location{}, "", err
}

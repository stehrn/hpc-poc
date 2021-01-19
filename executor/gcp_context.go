package executor

import (
	"fmt"

	"github.com/stehrn/hpc-poc/gcp/storage"
	"github.com/stehrn/hpc-poc/internal/utils"
)

// GcpContext context for GCP interaction
type GcpContext struct {
	Project    string
	Namespace  string
	BucketName string
	Business   string
}

// NewGcpContextFromEnv create GcpContext from environment
func NewGcpContextFromEnv() *GcpContext {
	return &GcpContext{
		Project:    utils.Env("PROJECT_NAME"),
		Namespace:  utils.Env("NAMESPACE"),
		BucketName: utils.Env("CLOUD_STORAGE_BUCKET_NAME"),
		Business:   utils.Env("BUSINESS_NAME"),
	}
}

// NewStorageClient create a new storage client
func (gcp GcpContext) NewStorageClient() (storage.ClientInterface, error) {
	return storage.NewClient()
}

// TopicName derive name of topic
func (gcp GcpContext) TopicName() string {
	return fmt.Sprintf("%s-%s-topic", gcp.Project, gcp.Business)
}

// SubscriptionName derive name of subscription
func (gcp GcpContext) SubscriptionName() string {
	return fmt.Sprintf("%s-%s-subscription", gcp.Project, gcp.Business)
}

// Location location for object
func (gcp GcpContext) Location(object string) storage.Location {
	return storage.Location{
		Bucket: gcp.BucketName,
		Object: object}
}

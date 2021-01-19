package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/stehrn/hpc-poc/internal/utils"
)

// Location location of payload
type Location struct {
	Bucket string `json:"bucket"`
	Object string `json:"object"`
}

// NewLocation generate storage bucket location for existing object
func NewLocation(bucket, object string) Location {
	return Location{
		Bucket: bucket,
		Object: object}
}

// NewLocationWithAutoGeneratedName generate storage bucket location
func NewLocationWithAutoGeneratedName(bucket, path string) Location {
	return Location{
		Bucket: bucket,
		Object: objectName(path, utils.GenerateID())}
}

// NewLocationFromFromEnv derive location from environment
func NewLocationFromFromEnv() Location {
	bucket := BucketNameFromEnv()
	object := utils.Env("CLOUD_STORAGE_OBJECT_NAME")
	return Location{bucket, object}
}

// IsDirectory does location represent a directory
func (l *Location) IsDirectory() bool {
	return strings.HasSuffix(l.Object, "/") || strings.HasSuffix(l.Object, "/*")
}

// BucketNameFromEnv derive name of bucket from env variable CLOUD_STORAGE_BUCKET_NAME
func BucketNameFromEnv() string {
	return utils.Env("CLOUD_STORAGE_BUCKET_NAME")
}

// BucketExists check whether bucket exists
func (c *Client) BucketExists(bucketName string) (bool, error) {
	bucket := c.Bucket(bucketName)
	ctx := context.Background()
	_, err := bucket.Attrs(ctx)
	if err != nil {
		return false, fmt.Errorf("storage.NewClient, bucket %s does not exist: %v", bucketName, err)
	}
	return true, nil
}

func objectName(path, id string) string {
	return fmt.Sprintf("%s/%s", path, id)
}

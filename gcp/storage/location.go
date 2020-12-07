package storage

import (
	"fmt"

	"github.com/rs/xid"
	"github.com/stehrn/hpc-poc/internal/utils"
)

// Location location of payload
type Location struct {
	Bucket string `json:"bucket"`
	Object string `json:"object"`
}

// Location generate storage bucket location
func (c Client) Location(business string) Location {
	return Location{
		Bucket: c.BucketName,
		Object: objectName(business, uniqueID())}
}

// LocationForObject generate storage bucket location for existing object
func (c Client) LocationForObject(object string) Location {
	return Location{
		Bucket: c.BucketName,
		Object: object}
}

// LocationFromEnv derive location from environment
func (c Client) LocationFromEnv() Location {
	bucket := utils.Env("CLOUD_STORAGE_BUCKET_NAME")
	object := utils.Env("CLOUD_STORAGE_OBJECT_NAME")
	return Location{bucket, object}
}

func objectName(business, id string) string {
	return fmt.Sprintf("%s/%s", business, id)
}

func uniqueID() string {
	return xid.New().String()
}

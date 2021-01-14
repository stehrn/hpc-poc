package storage

import (
	"fmt"
	"strings"

	"github.com/stehrn/hpc-poc/internal/utils"
)

// Location location of payload
type Location struct {
	Bucket string `json:"bucket"`
	Object string `json:"object"`
}

// IsDirectory does location represent a directory
func (l *Location) IsDirectory() bool {
	return strings.HasSuffix(l.Object, "/") || strings.HasSuffix(l.Object, "/*")
}

// Location generate storage bucket location
func (c *Client) Location(path string) Location {
	return Location{
		Bucket: c.BucketName(),
		Object: objectName(path, utils.GenerateID())}
}

// LocationForObject generate storage bucket location for existing object
func (c *Client) LocationForObject(object string) Location {
	return Location{
		Bucket: c.BucketName(),
		Object: object}
}

// LocationFromEnv derive location from environment
func (c *Client) LocationFromEnv() Location {
	bucket := utils.Env("CLOUD_STORAGE_BUCKET_NAME")
	object := utils.Env("CLOUD_STORAGE_OBJECT_NAME")
	return Location{bucket, object}
}

func objectName(path, id string) string {
	return fmt.Sprintf("%s/%s", path, id)
}

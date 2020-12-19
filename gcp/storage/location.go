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
func (c *Client) Location(business string) Location {
	return Location{
		Bucket: c.BucketName(),
		Object: objectName(business, utils.GenerateID())}
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

// ToLocationByteSlice concert object slice into byte slice
func (c *Client) ToLocationByteSlice(objects []Object) ([][]byte, error) {
	var results [][]byte
	for _, object := range objects {

		location := c.LocationForObject(object.Object)
		bytes, err := location.ToBytes()
		if err != nil {
			return nil, err
		}
		results = append(results, bytes)
	}
	return results, nil
}

func objectName(business, id string) string {
	return fmt.Sprintf("%s/%s", business, id)
}

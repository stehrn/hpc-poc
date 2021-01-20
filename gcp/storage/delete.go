package storage

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
)

// Delete delete object at given location
func (c *Client) Delete(location Location) error {
	if location.IsDirectory() {
		return c.deleteDirectory(location)
	}
	return c.deleteObject(location)
}

func (c *Client) deleteObject(location Location) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.Bucket(location.Bucket).Object(location.Object).Delete(ctx)
}

func (c *Client) deleteDirectory(location Location) error {
	ctx := context.Background()
	bucket := c.Bucket(location.Bucket)
	return c.ForEachObject(location, func(attrs *storage.ObjectAttrs) error {
		if err := bucket.Object(attrs.Name).Delete(ctx); err != nil {
			return fmt.Errorf("Could not delete directory: %v", err)
		}
		return nil
	})
}

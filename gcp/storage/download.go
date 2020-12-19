package storage

import (
	"context"
	"fmt"
	"io/ioutil"
)

// Download download an object
func (c *Client) Download(location Location) ([]byte, error) {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	rc, err := c.Bucket(location.Bucket).Object(location.Object).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("Object(%q).NewReader: %v", location.Object, err)
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll: %v", err)
	}
	return data, nil
}

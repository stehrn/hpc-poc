package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/stehrn/hpc-poc/internal/utils"
)

// Location location of payload
type Location struct {
	Bucket string `json:"bucket"`
	Object string `json:"object"`
}

// LocationFromEnvironment  derive location from environment
func (c Client) LocationFromEnvironment() Location {
	bucket := utils.Env("BUCKET_NAME")
	object := utils.Env("OBJECT_NAME")
	return Location{bucket, object}
}

// Upload upload object to Cloud Storage bucket
func (c Client) Upload(location Location, content []byte) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := c.Bucket(location.Bucket).Object(location.Object).NewWriter(ctx)
	if _, err := io.Copy(wc, bytes.NewReader(content)); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	return nil
}

// Download download an object
func (c Client) Download(location Location) ([]byte, error) {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
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

// Delete delete object at given location
func (c Client) Delete(location Location) error {
	ctx := context.Background()

	defer c.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	return c.Bucket(location.Bucket).Object(location.Object).Delete(ctx)
}

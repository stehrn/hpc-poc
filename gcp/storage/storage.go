package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

const timeout = time.Second * 50

// ForEachObject iterate over each storage object, whic his passed to consumer
func (c *Client) ForEachObject(prefix string, consumer func(attrs *storage.ObjectAttrs)) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var query *storage.Query
	if prefix != "" {
		query = &storage.Query{Prefix: prefix}
	}

	it := c.Bucket(c.BucketName()).Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("ListStorageObjects for prefix: '%s', error: %v", prefix, err)
		}
		consumer(attrs)
	}
	return nil
}

// ListObjects list storage objects
func (c *Client) ListObjects(prefix string) ([]Object, error) {
	var objects []Object
	err := c.ForEachObject(prefix, func(attrs *storage.ObjectAttrs) {
		object := Object{
			Object:  attrs.Name,
			Size:    attrs.Size,
			Created: attrs.Created}
		objects = append(objects, object)
	})
	return objects, err
}

// Upload upload object to Cloud Storage bucket
func (c *Client) Upload(location Location, content []byte) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, timeout)
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

// Delete delete object at given location
func (c *Client) Delete(location Location) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.Bucket(location.Bucket).Object(location.Object).Delete(ctx)
}

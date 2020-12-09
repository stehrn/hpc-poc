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

// Object compact representation of a storage object
type Object struct {
	Object  string
	Size    int64
	Created time.Time
}

// List list contents of bucket
// equivallent to: gsutil ls -r gs://c.BucketName/prefix**
func (c *Client) List(prefix string) *storage.ObjectIterator {

	ctx := context.Background()

	// ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	// defer cancel()

	var query *storage.Query
	if prefix != "" {
		query = &storage.Query{Prefix: prefix}
	}
	return c.Bucket(c.BucketName).Objects(ctx, query)
}

// ListStorageObjects list storage objects
func (c *Client) ListStorageObjects(prefix string) ([]Object, error) {
	var objects []Object
	it := c.List(prefix)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("ListStorageObjects for prefix: '%s', error: %v", prefix, err)
		}
		object := Object{
			Object:  attrs.Name,
			Size:    attrs.Size,
			Created: attrs.Created}
		objects = append(objects, object)
	}
	return objects, nil
}

// Upload upload object to Cloud Storage bucket
func (c *Client) Upload(location Location, content []byte) error {
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
func (c *Client) Download(location Location) ([]byte, error) {
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

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	return c.Bucket(location.Bucket).Object(location.Object).Delete(ctx)
}

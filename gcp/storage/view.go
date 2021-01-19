package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

const timeout = time.Second * 50

// ListObjects list storage objects
func (c *Client) ListObjects(location Location) ([]Object, error) {
	prefix := strings.Trim(location.Object, "/")
	var objects []Object
	err := c.ForEachObject(NewLocation(location.Bucket, prefix), func(attrs *storage.ObjectAttrs) error {
		object := Object{
			Object:  attrs.Name,
			Size:    attrs.Size,
			Created: attrs.Created}
		objects = append(objects, object)
		return nil
	})
	return objects, err
}

// ForEachObject iterate over each storage object, which is passed to consumer
func (c *Client) ForEachObject(location Location, consumer func(attrs *storage.ObjectAttrs) error) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var errors error
	var query *storage.Query
	if location.Object != "" {
		query = &storage.Query{Prefix: location.Object}
	}

	bucket := c.Bucket(location.Bucket)
	it := bucket.Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("ListStorageObjects for prefix: '%s', error: %w", location.Object, err)
		}
		err = consumer(attrs)
		if err != nil {
			multierror.Append(errors, err)
		}
	}
	return errors
}

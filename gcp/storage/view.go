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
func (c *Client) ListObjects(prefix string) ([]Object, error) {
	prefix = strings.Trim(prefix, "/")
	var objects []Object
	err := c.ForEachObject(prefix, func(attrs *storage.ObjectAttrs) error {
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
func (c *Client) ForEachObject(prefix string, consumer func(attrs *storage.ObjectAttrs) error) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var errors error
	var query *storage.Query
	if prefix != "" {
		query = &storage.Query{Prefix: prefix}
	}

	bucket := c.Bucket(c.BucketName())
	it := bucket.Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("ListStorageObjects for prefix: '%s', error: %w", prefix, err)
		}
		err = consumer(attrs)
		if err != nil {
			multierror.Append(errors, err)
		}
	}
	return errors
}

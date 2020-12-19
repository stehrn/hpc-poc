package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/stehrn/hpc-poc/client"
)

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

// UploadMany upload many data items
// scope to look into batch API here
// no limits on number of threads
func (c *Client) UploadMany(items client.DataSourceIterator) {
	fmt.Printf("Uploading %d items\n", items.Size())
	var wg sync.WaitGroup
	wg.Add(items.Size())
	items.Each(func(item client.DataSource) {
		go func(item client.DataSource) {
			defer wg.Done()
			
			location := Location{
				Bucket: c.BucketName(),
				Object: item.ObjectPath()}

			log.Printf("Uploading data to: '%v'\n", location)
			err := c.Upload(location, item.Data())
			if err != nil {
				log.Printf("Error uploading data to: '%v': %v\n", location, err)
				item.AddError(err)
			}
		}(item)
	})
	wg.Wait()
}

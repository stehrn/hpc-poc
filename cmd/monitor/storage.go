package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/stehrn/hpc-poc/gcp/storage"
	"google.golang.org/api/iterator"
)

var storageClient *storage.Client

type storageTemplate struct {
	Business string
	Bucket   string
	Objects  []storageObject
}

type storageObject struct {
	Object  string
	Size    int64
	Created time.Time
}

func init() {
	var err error
	storageClient, err = storage.NewEnvClient()
	if err != nil {
		log.Fatalf("Could not create storage client: %v", err)
	}
}

// uri pattern: /bucket/<business>
func (ctx *handlerContext) BucketHandler(w http.ResponseWriter, r *http.Request) error {
	split := strings.Split(r.URL.Path, "/")
	business := split[2]
	if business == "" {
		return fmt.Errorf("BucketHandler() expected /bucket/<business>, got: %v", r.URL.Path)
	}

	log.Printf("Listing objects for bucket: %s, business: %s", storageClient.BucketName, business)
	var objects []storageObject
	it := storageClient.List(business)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("Bucket(%q).Objects, for business: '%s', error: %v", storageClient.BucketName, business, err)
		}
		object := storageObject{
			Object:  attrs.Name,
			Size:    attrs.Size,
			Created: attrs.Created}
		objects = append(objects, object)
	}
	return ctx.bucketTemplate.Execute(w, storageTemplate{
		Business: business,
		Bucket:   storageClient.BucketName,
		Objects:  objects})
}

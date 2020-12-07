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

// uri pattern one of:
//    /bucket/business/<business> to list objects
//    /bucket/object/<object> t0 view object
func (ctx *handlerContext) BucketHandler(w http.ResponseWriter, r *http.Request) error {
	split := strings.Split(r.URL.Path, "/")
	if len(split) < 3 {
		return fmt.Errorf("BucketHandler() bad request, expected /bucket/<business|object>/<value>, got: %v", r.URL.Path)
	}
	action := split[2]
	if action == "business" {
		business := split[3]
		return ctx.objects(business, w)
	} else if action == "object" {
		object := strings.Join(split[3:], "/")
		return ctx.download(object, w)
	} else {
		return fmt.Errorf("BucketHandler() expected /bucket/<business|object>, got: %v", r.URL.Path)
	}
}

func (ctx *handlerContext) objects(business string, w http.ResponseWriter) error {
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

func (ctx *handlerContext) download(object string, w http.ResponseWriter) error {
	location := storageClient.LocationForObject(object)
	log.Printf("Downloading object from location: %v", location)
	data, err := storageClient.Download(location)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/stehrn/hpc-poc/gcp/storage"
)

var storageClient *storage.Client

type storageTemplate struct {
	Business string
	Bucket   string
	Objects  []storage.Object
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
	objects, err := storageClient.ListStorageObjects(business)
	if err != nil {
		return err
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

package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/stehrn/hpc-poc/gcp/storage"
)

var storageClient storage.ClientInterface

type storageTemplate struct {
	Business            string
	Bucket              string
	Objects             []storage.Object
	BusinessNameOptions []BusinessNameOptions
	Page                string
}

func init() {
	var err error
	storageClient, err = storage.NewEnvClient()
	if err != nil {
		log.Fatalf("Could not create storage client: %v", err)
	}
}

// uri pattern one of:
//    /storage/business/<business> to list objects
//    /storage/object/<object> t0 view object
func (ctx *handlerContext) StorageHandler(w http.ResponseWriter, r *http.Request) error {

	log.Printf("ok: " + r.Method + ", uri:" + r.URL.Path)

	split := strings.Split(r.URL.Path, "/")
	action := split[2]
	if action == "" {
		business := "bu3"
		return ctx.objects(business, w)
	} else if action == "business" {
		business := split[3]
		return ctx.objects(business, w)
	} else if action == "object" {
		object := strings.Join(split[3:], "/")
		return ctx.download(object, w)
	} else {
		return fmt.Errorf("BucketHandler() expected /storage/<business|object>, got: %v", r.URL.Path)
	}

	// business := "bu3"
	// summary := storageTemplate{
	// 	BusinessNameOptions: NewBusinessNameOptions(business),
	// 	Page:                "storage"}
	// return ctx.storageTemplate.Execute(w, summary)
	// return ctx.objects(business, w)

	// switch r.Method {
	// case "GET":
	// 	summary := storageTemplate{
	// 		BusinessNameOptions: NewBusinessNameOptions(""),
	// 		Page:                "storage"}
	// 	return ctx.storageTemplate.Execute(w, summary)
	// case "POST":
	// 	business := r.FormValue("business")
	// 	return ctx.objects(business, w)

	// }
	// return nil

	// split := strings.Split(r.URL.Path, "/")
	// if len(split) < 3 {
	// 	return fmt.Errorf("BucketHandler() bad request, expected /bucket/<business|object>/<value>, got: %v", r.URL.Path)
	// }
	// action := split[2]

}

func (ctx *handlerContext) objects(business string, w http.ResponseWriter) error {
	log.Printf("Listing objects for bucket: %s, business: %s", storageClient.BucketName, business)
	objects, err := storageClient.ListObjects(business)
	if err != nil {
		return err
	}
	return ctx.storageTemplate.Execute(w, storageTemplate{
		Business:            business,
		Bucket:              storageClient.BucketName(),
		Objects:             objects,
		BusinessNameOptions: NewBusinessNameOptions(business),
		Page:                "storage",
	})
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

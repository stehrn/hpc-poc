package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/stehrn/hpc-poc/internal/utils"

	"github.com/rs/xid"
	"github.com/stehrn/hpc-poc/gcp/pubsub"
	"github.com/stehrn/hpc-poc/gcp/storage"
)

// PageInfo info to render into template
type PageInfo struct {
	*pubsub.Client
	Bucket  string
	Message string
}

type handlerContext struct {
	client   *pubsub.Client
	Bucket   string
	template *template.Template
}

func (ctx *handlerContext) handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		ctx.template.Execute(w, &PageInfo{ctx.client, ctx.Bucket, ""})
	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		data := []byte(r.FormValue("payload"))

		// upload payload to cloud storage
		location := storage.Location{
			Bucket: ctx.Bucket,
			Object: xid.New().String()}
		log.Printf("Uploading data to %v", location)
		err := storage.Upload(location, data)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// publish object location
		log.Printf("Publishing location (%v) to topic %s", location, ctx.client.Topic)
		bytes, err := storage.ToBytes(location)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		id, err := ctx.client.Publish(bytes)
		if err != nil {
			http.Error(w, err.Error(), 500)
		} else {
			ctx.template.Execute(w, &PageInfo{ctx.client, ctx.Bucket, fmt.Sprintf("Payload uploaded to cloud storage location: %s, notification send with ID: %s", location, id)})
		}
	}
}

func main() {
	log.Print("Starting client")

	templatePath := os.Getenv("TEMPLATE_PATH")
	clientTemplate := filepath.Join(templatePath, "./index.tmpl")
	log.Printf("Loading template from: %s", clientTemplate)

	pubsubClient, err := pubsub.NewClient()
	if err != nil {
		log.Fatalf("Could not create pubsub client: %v", err)
	}

	bucket := utils.Env("BUCKET_NAME")
	ctx := &handlerContext{
		client:   pubsubClient,
		Bucket:   bucket,
		template: template.Must(template.ParseFiles(clientTemplate))}

	http.HandleFunc("/client", ctx.handle)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Client created for project: %s, topic: %s, bucket: %s; listening on port %s",
		pubsubClient.Project, pubsubClient.Topic, bucket, port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

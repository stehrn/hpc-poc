package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/stehrn/hpc-poc/internal/utils"

	"github.com/stehrn/hpc-poc/gcp/pubsub"
)

// templateData data to render into template
type templateData struct {
	*Client
	Message string
}

type handlerContext struct {
	client   *Client
	template *template.Template
}

func (ctx handlerContext) templateData(message string) *templateData {
	return &templateData{ctx.client, message}
}

func errorHandler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (ctx *handlerContext) handle(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		ctx.template.Execute(w, ctx.templateData(""))
	case "POST":
		if err := r.ParseForm(); err != nil {
			return fmt.Errorf("ParseForm() err: %v", err)
		}

		data := []byte(r.FormValue("payload"))
		location, id, err := ctx.client.handle(data)
		if err != nil {
			return fmt.Errorf("client.handle() err: %v", err)
		}

		message := fmt.Sprintf("Payload uploaded to cloud storage location: %s, notification sent with message ID: %s", location, id)
		log.Print(message)
		ctx.template.Execute(w, ctx.templateData(message))
	}
	return nil
}

func main() {
	log.Print("Starting client")

	templatePath := os.Getenv("TEMPLATE_PATH")
	clientTemplate := filepath.Join(templatePath, "./index.tmpl")
	log.Printf("Loading template from: %s", clientTemplate)

	bucket := utils.Env("BUCKET_NAME")
	client, err := NewClient(bucket, pubsub.ConfigFromEnvironment())
	if err != nil {
		log.Fatalf("Could not create client: %v", err)
	}

	ctx := &handlerContext{
		client:   client,
		template: template.Must(template.ParseFiles(clientTemplate))}

	http.HandleFunc("/client", errorHandler(ctx.handle))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Client created for project: %s, topic: %s, bucket: %s; listening on port %s",
		client.Pubsub.Project, client.Pubsub.TopicName, bucket, port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

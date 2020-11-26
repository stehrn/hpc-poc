package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	gcp "github.com/stehrn/hpc-poc/gcp"
)

// PageInfo info to render into template
type PageInfo struct {
	Gcp     gcp.Info
	Message string
}

type handlerContext struct {
	gcpInfo  gcp.Info
	client   *gcp.Client
	template *template.Template
}

func (ctx *handlerContext) handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		ctx.template.Execute(w, &PageInfo{ctx.gcpInfo, ""})
	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		payload := []byte(r.FormValue("payload"))
		id, err := ctx.client.Publish(payload)
		if err != nil {
			http.Error(w, err.Error(), 500)
		} else {
			ctx.template.Execute(w, &PageInfo{ctx.gcpInfo, "Published payload, message ID " + id})
		}
	}
}

func main() {
	log.Print("Starting client")

	templatePath := os.Getenv("TEMPLATE_PATH")
	clientTemplate := filepath.Join(templatePath, "./index.tmpl")
	log.Printf("Loading template from: %s", clientTemplate)

	gcpInfo := gcp.InfoFromEnvironment()
	log.Printf("Creating gcp client for project: %s, subscriptionID: %s, topic: %s", gcpInfo.Project, gcpInfo.Subscription, gcpInfo.Topic)
	gcpClient, err := gcp.NewClient(gcpInfo)
	if err != nil {
		log.Fatalf("Could not create gcp client: %v", err)
	}
	ctx := &handlerContext{
		gcpInfo:  gcpInfo,
		client:   gcpClient,
		template: template.Must(template.ParseFiles(clientTemplate))}

	http.HandleFunc("/client", ctx.handle)

	port := os.Getenv("CLIENT_PORT")
	if port == "" {
		port = "8082"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Service Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

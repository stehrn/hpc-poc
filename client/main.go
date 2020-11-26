package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/stehrn/hpc-poc/gcp"
)

type GcpInfo struct {
	Project      string
	Subscription string
	Topic        string
}

type PageInfo struct {
	Gcp     GcpInfo
	Message string
}

type handlerContext struct {
	gcpInfo  GcpInfo
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
		err := ctx.client.Publish(payload)
		if err != nil {
			http.Error(w, err.Error(), 500)
		} else {
			ctx.template.Execute(w, &PageInfo{ctx.gcpInfo, "Published payload"})
		}
	}
}

func main() {
	log.Print("Starting client")

	cwd, _ := os.Getwd()
	clientTemplate := filepath.Join(cwd, "./index.tmpl")
	log.Printf("Loading template from: %s", clientTemplate)

	gcpInfo := GcpInfoFromEnvironment()
	gcpClient, err := gcpClient(gcpInfo)
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

func gcpClient(gcpInfo GcpInfo) (*gcp.Client, error) {
	log.Printf("Creating gcp client for project: %s, subscriptionID: %s, topic: %s", gcpInfo.Project, gcpInfo.Subscription, gcpInfo.Topic)
	return gcp.NewClient(gcpInfo.Project, gcpInfo.Subscription, gcpInfo.Topic)
}

// GcpInfoFromEnvironment create GcpInfoFromEnvironment from environment variables: PROJECT_NAME, SUBSCRIPTION_NAME, TOPIC_NAME
func GcpInfoFromEnvironment() GcpInfo {
	return GcpInfo{Project: "hpc-poc",
		Subscription: env("SUBSCRIPTION_NAME"),
		Topic:        env("TOPIC_NAME")}
}

func env(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("No '%s' env variable", key)
	}
	return value
}

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/stehrn/hpc-poc/client"
	http_common "github.com/stehrn/hpc-poc/internal/http"
	"github.com/stehrn/hpc-poc/internal/utils"
)

var businessNames []string

func init() {
	names := utils.Env("BUSINESS_NAMES")
	businessNames = strings.Split(names, ",")
}

// templateData data to render into template
type templateData struct {
	*client.Client
	BusinessNames []string
	Message       string
}

type handlerContext struct {
	client   *client.Client
	template *template.Template
}

func (ctx handlerContext) templateData(message string) *templateData {
	return &templateData{ctx.client, businessNames, message}
}

func (ctx *handlerContext) handle(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		ctx.template.Execute(w, ctx.templateData(""))
	case "POST":
		if err := r.ParseForm(); err != nil {
			return fmt.Errorf("ParseForm() err: %v", err)
		}

		business := client.Business(r.FormValue("business"))
		data := []byte(r.FormValue("payload"))
		location, id, err := ctx.client.Handle(business, data)
		if err != nil {
			return fmt.Errorf("client.handle() err: %v", err)
		}

		topic := business.TopicName(ctx.client.Project)
		message := fmt.Sprintf("Data uploaded to cloud storage location: %s, notification sent to topic: '%s', message ID: %s", location, topic, id)
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

	client := client.NewEnvClientOrFatal()
	ctx := &handlerContext{
		client:   client,
		template: template.Must(template.ParseFiles(clientTemplate))}

	handler := http_common.ErrorHandler(ctx.handle)
	http.HandleFunc("/", handler)
	http.HandleFunc("/client", handler)

	port := http_common.Port()
	log.Printf("Client created for project: %s, business names: %v, bucket: %s; listening on port %s",
		client.Project, businessNames, client.Storage.BucketName, port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

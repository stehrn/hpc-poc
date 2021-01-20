package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/stehrn/hpc-poc/client"
	"github.com/stehrn/hpc-poc/executor"
	"github.com/stehrn/hpc-poc/gcp/storage"
	http_common "github.com/stehrn/hpc-poc/internal/http"
	"github.com/stehrn/hpc-poc/internal/utils"
)

var businessNames []string
var jobCounter uint64

func init() {
	businessNames = strings.Split(utils.Env("BUSINESS_NAMES"), ",")
}

// templateData data to render into template
type templateData struct {
	*executor.GcpContext
	BusinessNames []string
	Message       string
}

type handlerContext struct {
	gcpContext *executor.GcpContext
	template   *template.Template
}

func (ctx handlerContext) templateData(message string) *templateData {
	return &templateData{ctx.gcpContext, businessNames, message}
}

func (ctx *handlerContext) handle(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		ctx.template.Execute(w, ctx.templateData(""))
	case "POST":
		if err := r.ParseForm(); err != nil {
			return fmt.Errorf("ParseForm() err: %v", err)
		}

		business := r.FormValue("business")

		session := client.NewSession("web-client-session", business)
		atomic.AddUint64(&jobCounter, 1)
		job := client.NewJob(fmt.Sprintf("web-client-job-%d", jobCounter), session)
		data := []byte(r.FormValue("payload"))
		job.CreateTask(data)

		ctx.gcpContext.Business = business
		exe, err := executor.New(ctx.gcpContext)
		if err != nil {
			log.Fatalf("Error creating client, %v", err)
		}

		result := exe.Execute(job)
		if result.Error != nil {
			return fmt.Errorf("Error executing job: %v", result.Error)
		}

		message := fmt.Sprintf("Data uploaded to cloud storage location: %s, notification sent to topic: '%s', message ID: %s", job.ObjectPath().String(), exe.TopicName(), result.MessageID)
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

	gcpContext := &executor.GcpContext{
		Project:    "hpc-poc",
		BucketName: storage.BucketNameFromEnv(),
	}

	ctx := &handlerContext{
		gcpContext: gcpContext,
		template:   template.Must(template.ParseFiles(clientTemplate))}

	handler := http_common.ErrorHandler(ctx.handle)
	http.HandleFunc("/", handler)
	http.HandleFunc("/client", handler)

	port := http_common.Port()
	log.Printf("Client created for project: %s, business names: %v, bucket: %s; listening on port %s",
		gcpContext.Project, businessNames, gcpContext.BucketName, port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

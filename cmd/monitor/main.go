package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/stehrn/hpc-poc/client"
	http_common "github.com/stehrn/hpc-poc/internal/http"
	k8 "github.com/stehrn/hpc-poc/kubernetes"
)

var businessNames []string

func init() {
	businessNames = client.BusinessNamesFromEnv()
}

type handlerContext struct {
	client          *k8.Client
	summaryTemplate *template.Template
	jobTemplate     *template.Template
	bucketTemplate  *template.Template
}

var templatePath string

func init() {
	templatePath = os.Getenv("TEMPLATE_PATH")
	if templatePath != "" {
		log.Printf("Loading template from path: %s", templatePath)
	}
}

func main() {
	log.Print("Starting monitor")

	client, err := k8.NewEnvClient()
	if err != nil {
		log.Fatal(err)
	}
	ctx := &handlerContext{
		client:          client,
		summaryTemplate: loadTemplate("./summary.tmpl"),
		jobTemplate:     loadTemplate("./job.tmpl"),
		bucketTemplate:  loadTemplate("./storage.tmpl")}

	summaryHandler := http_common.ErrorHandler(ctx.SummaryHandler)
	http.HandleFunc("/", summaryHandler)
	http.HandleFunc("/summary", summaryHandler)
	http.HandleFunc("/job/", http_common.ErrorHandler(ctx.JobHandler))
	http.HandleFunc("/logs/", http_common.ErrorHandler(ctx.LogsHandler))
	http.HandleFunc("/bucket/", http_common.ErrorHandler(ctx.BucketHandler))

	port := http_common.Port()
	log.Printf("Monitor service Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func loadTemplate(name string) *template.Template {
	myTemplate := filepath.Join(templatePath, name)
	return template.Must(template.ParseFiles(myTemplate))
}

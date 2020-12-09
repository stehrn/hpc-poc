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

// BusinessNameOptions
type BusinessNameOptions struct {
	Name     string
	Selected bool
}

// NewBusinessNameOptions
func NewBusinessNameOptions(selected string) []BusinessNameOptions {
	options := make([]BusinessNameOptions, len(businessNames))
	for i, name := range businessNames {
		var isOptSelected bool
		if name == selected {
			isOptSelected = true
		}
		options[i] = BusinessNameOptions{name, isOptSelected}
	}
	return options
}

func init() {
	businessNames = client.BusinessNamesFromEnv()
}

type handlerContext struct {
	client          *k8.Client
	summaryTemplate *template.Template
	jobTemplate     *template.Template
	storageTemplate *template.Template
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
		summaryTemplate: loadTemplate("./summary.tmpl", "./option.tmpl", "./navbar.tmpl"),
		jobTemplate:     loadTemplate("./job.tmpl", "./navbar.tmpl"),
		storageTemplate: loadTemplate("./storage.tmpl", "./option.tmpl", "./navbar.tmpl")}

	summaryHandler := http_common.ErrorHandler(ctx.SummaryHandler)
	http.HandleFunc("/", summaryHandler)
	http.HandleFunc("/summary", summaryHandler)
	http.HandleFunc("/job/", http_common.ErrorHandler(ctx.JobHandler))
	http.HandleFunc("/logs/", http_common.ErrorHandler(ctx.LogsHandler))
	http.HandleFunc("/storage/", http_common.ErrorHandler(ctx.StorageHandler))

	port := http_common.Port()
	log.Printf("Monitor service Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func loadTemplate(names ...string) *template.Template {
	var templates []string
	for _, name := range names {
		fullPath := filepath.Join(templatePath, name)
		templates = append(templates, fullPath)
	}
	return template.Must(template.ParseFiles(templates...))
}

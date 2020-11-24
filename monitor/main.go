package main

import (
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	k8 "github.com/stehrn/hpc-poc/kubernetes"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type job struct {
	Name           string
	Status         string
	StartTime      string
	CompletionTime string
	Duration       string
	Logs           string
}

// JobList a list of jobs
type JobList struct {
	Namespace    string
	Subscription string
	Jobs         []job
}

type jobHandlerContext struct {
	client   *k8.Client
	template *template.Template
}

func (h *jobHandlerContext) handler(w http.ResponseWriter, r *http.Request) {
	jobs, err := jobs(h.client)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	h.template.Execute(w, jobs)
}

func main() {
	log.Print("Starting monitor")

	namespace := "default"
	ctx := &jobHandlerContext{
		client:   k8.NewClient(namespace),
		template: template.Must(template.ParseFiles("jobs.tmpl"))}

	http.HandleFunc("/", ctx.handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Service Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// https://gowalker.org/k8s.io/api/batch/v1#JobList
func jobs(client *k8.Client) (JobList, error) {
	var jobs []job
	jobList, err := client.ListJobs()
	if err != nil {
		return JobList{}, err
	}
	for _, item := range jobList.Items {
		item := job{
			Name:           item.Name,
			Status:         status(item.Status),
			StartTime:      toString(item.Status.StartTime),
			CompletionTime: toString(item.Status.CompletionTime),
			Duration:       duration(item.Status.StartTime, item.Status.CompletionTime),
			Logs:           "link to logs"}
		jobs = append(jobs, item)
	}
	return JobList{
		Namespace:    "namesapce",
		Subscription: "sub",
		Jobs:         jobs}, nil
}

func duration(start, end *v1.Time) string {
	if start.IsZero() || end.IsZero() {
		return ""
	}
	startStr, err := start.MarshalQueryParameter()
	if err != nil {
		return ""
	}
	endStr, err := end.MarshalQueryParameter()
	if err != nil {
		return ""
	}
	startTime, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		return ""
	}
	endTime, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		return ""
	}
	return endTime.Sub(startTime).String()
}

func toString(time *v1.Time) string {
	if time.IsZero() {
		return ""
	}
	return time.String()
}

// ssumes we only have 1 job
func status(status batchv1.JobStatus) string {
	if status.Active > 0 {
		return "Job is still running"
	} else if status.Succeeded > 0 {
		return "Job Successful"
	} else if status.Failed > 0 {
		return "Job Failed"
	}
	return "Job has no status"
}

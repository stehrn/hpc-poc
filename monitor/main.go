package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	apiv1 "k8s.io/api/core/v1"

	k8 "github.com/stehrn/hpc-poc/kubernetes"
)

type podSummary struct {
	Name      string
	Status    string
	LastState apiv1.PodCondition
}

type jobSummary struct {
	Name           string
	Status         string
	StartTime      string
	CompletionTime string
	Duration       string
	Pod            podSummary
}

// JobList a list of jobs
type JobList struct {
	Namespace string
	Jobs      []jobSummary
}

type handlerContext struct {
	client   *k8.Client
	template *template.Template
}

func (ctx *handlerContext) jobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := jobs(ctx.client)
	if err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		ctx.template.Execute(w, jobs)
	}
}

func (ctx *handlerContext) logs(w http.ResponseWriter, r *http.Request) {
	job := strings.TrimPrefix(r.URL.Path, "/job/log/")
	if job == "" {
		http.Error(w, "No job sepcified!", 400)
		return
	}

	log.Printf("Loading log for job %s", job)
	logs, err := ctx.client.LogsForJob(job)
	if err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		w.Write([]byte(logs))
	}
}

func main() {
	log.Print("Starting monitor")

	templatePath := os.Getenv("TEMPLATE_PATH")
	jobsTemplate := filepath.Join(templatePath, "./jobs.tmpl")
	log.Printf("Loading template from path: %s", jobsTemplate)

	client, err := k8.NewClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	ctx := &handlerContext{
		client:   client,
		template: template.Must(template.ParseFiles(jobsTemplate))}

	http.HandleFunc("/jobs", ctx.jobs)
	http.HandleFunc("/job/log/", ctx.logs)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Monitor service Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// https://gowalker.org/k8s.io/api/batch/v1#JobList
func jobs(client *k8.Client) (JobList, error) {
	var jobs []jobSummary
	jobList, err := client.ListJobs()
	if err != nil {
		return JobList{}, err
	}
	for _, item := range jobList.Items {
		item := jobSummary{
			Name:           item.Name,
			Status:         k8.Status(item.Status),
			StartTime:      k8.ToString(item.Status.StartTime),
			CompletionTime: k8.ToString(item.Status.CompletionTime),
			Duration:       k8.Duration(item.Status.StartTime, item.Status.CompletionTime),
			Pod:            getPodSummary(client, item.Name)}
		jobs = append(jobs, item)
	}
	return JobList{
		Namespace: client.Namespace,
		Jobs:      jobs}, nil
}

// Pod Status
// Last Pod State (type/reason/message)
func getPodSummary(client *k8.Client, jobName string) podSummary {
	pod, _ := client.Pod(jobName)
	podStatus := pod.Status
	conditions := podStatus.Conditions
	var lastState apiv1.PodCondition
	lastState = apiv1.PodCondition{}
	if len(conditions) != 0 {
		lastState = conditions[0]
	}
	return podSummary{
		Name:      pod.Name,
		Status:    string(podStatus.Phase),
		LastState: lastState}
}

package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	batchv1 "k8s.io/api/batch/v1"
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

// JobTemplate
type JobTemplate struct {
	Job  *batchv1.Job
	Pods []apiv1.Pod
}

// Summary summary information aobut jobs
type Summary struct {
	Namespace string
	Jobs      []jobSummary
}

type handlerContext struct {
	client          *k8.Client
	summaryTemplate *template.Template
	jobTemplate     *template.Template
}

var templatePath string

func init() {
	templatePath := os.Getenv("TEMPLATE_PATH")
	log.Printf("Loading template from path: %s", templatePath)
}

func (ctx *handlerContext) summary(w http.ResponseWriter, r *http.Request) {
	summary, err := summary(ctx.client)
	if err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		ctx.summaryTemplate.Execute(w, summary)
	}
}

func (ctx *handlerContext) job(w http.ResponseWriter, r *http.Request) {
	jobName := strings.TrimPrefix(r.URL.Path, "/job/")
	if jobName == "" {
		http.Error(w, "No job sepcified", 400)
		return
	}

	job, _ := ctx.client.Job(jobName)
	job, err := ctx.client.Job(jobName)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	pods, err := ctx.client.Pods(jobName)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	ctx.jobTemplate.Execute(w, JobTemplate{job, pods})
}

func (ctx *handlerContext) logs(w http.ResponseWriter, r *http.Request) {
	split := strings.Split(r.URL.Path, "/")
	objectType := split[2]
	name := split[3]

	var logs string
	var err error
	if objectType == "job" {
		logs, err = ctx.client.LogsForJob(name)
	} else if objectType == "pod" {
		logs, err = ctx.client.LogsForPod(name)
	} else {
		http.Error(w, "Unkown object type: "+objectType, 400)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		w.Write([]byte(logs))
	}
}

func main() {
	log.Print("Starting monitor")

	client, err := k8.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	ctx := &handlerContext{
		client:          client,
		summaryTemplate: loadTemplate("./summary.tmpl"),
		jobTemplate:     loadTemplate("./job.tmpl")}

	http.HandleFunc("/summary", ctx.summary)
	http.HandleFunc("/job/", ctx.job)
	http.HandleFunc("/logs/", ctx.logs)

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

func loadTemplate(name string) *template.Template {
	myTemplate := filepath.Join(templatePath, name)
	return template.Must(template.ParseFiles(myTemplate))
}

// https://gowalker.org/k8s.io/api/batch/v1#JobList
func summary(client *k8.Client) (Summary, error) {
	var jobs []jobSummary
	jobList, err := client.ListJobs()
	if err != nil {
		return Summary{}, err
	}
	for _, item := range jobList.Items {
		job := jobSummary{
			Name:           item.Name,
			Status:         k8.Status(item.Status),
			StartTime:      k8.ToString(item.Status.StartTime),
			CompletionTime: k8.ToString(item.Status.CompletionTime),
			Duration:       k8.Duration(item.Status.StartTime, item.Status.CompletionTime),
			Pod:            getPodSummary(client, item.Name)}
		jobs = append(jobs, job)
	}
	return Summary{
		Namespace: client.Namespace,
		Jobs:      jobs}, nil
}

// Pod Status
// Last Pod State (type/reason/message)
func getPodSummary(client *k8.Client, jobName string) podSummary {
	pod, _ := client.LatestPod(jobName)
	podStatus := pod.Status
	conditions := podStatus.Conditions
	var lastState = apiv1.PodCondition{}
	if len(conditions) != 0 {
		lastState = conditions[0]
	}
	return podSummary{
		Name:      pod.Name,
		Status:    string(podStatus.Phase),
		LastState: lastState}
}

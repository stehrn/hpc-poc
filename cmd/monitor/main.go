package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"

	http_common "github.com/stehrn/hpc-poc/internal/http"
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

// so we can add our own methods
type myJob struct {
	*batchv1.Job
}

// jobsTemplate data to render into jobs template
type jobsTemplate struct {
	Job  myJob
	Pods []apiv1.Pod
}

type lastPod struct {
	Condition apiv1.PodCondition
	IsError   bool
}

func (j jobsTemplate) LastPod() lastPod {
	if len(j.Pods) != 0 {
		conditions := j.Pods[0].Status.Conditions
		if len(conditions) != 0 {
			condition := conditions[0]
			var jobError bool
			if condition.Reason == "Unschedulable" {
				jobError = true
			}
			return lastPod{condition, jobError}
		}
	}
	return lastPod{apiv1.PodCondition{}, false}
}

func (j myJob) ContainerEnv() []apiv1.EnvVar {
	return j.Spec.Template.Spec.Containers[0].Env
}

// summaryTemplate data to render into summary template
type summaryTemplate struct {
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
	templatePath = os.Getenv("TEMPLATE_PATH")
	if templatePath != "" {
		log.Printf("Loading template from path: %s", templatePath)
	}
}

func (ctx *handlerContext) summary(w http.ResponseWriter, r *http.Request) error {
	summary, err := summary(ctx.client)
	if err != nil {
		return fmt.Errorf("Error creating summary page: %v", err)
	}
	ctx.summaryTemplate.Execute(w, summary)
	return nil
}

func (ctx *handlerContext) job(w http.ResponseWriter, r *http.Request) error {
	jobName := strings.TrimPrefix(r.URL.Path, "/job/")
	if jobName == "" {
		return errors.New("No job sepcified")
	}

	job, _ := ctx.client.Job(jobName)
	job, err := ctx.client.Job(jobName)
	if err != nil {
		return err
	}

	pods, err := ctx.client.Pods(jobName)
	if err != nil {
		return err
	}

	ctx.jobTemplate.Execute(w, jobsTemplate{myJob{job}, pods})
	return nil
}

func (ctx *handlerContext) logs(w http.ResponseWriter, r *http.Request) error {
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
		return fmt.Errorf("Unkown object type: %v", objectType)
	}

	if err != nil {
		return err
	}

	w.Write([]byte(logs))
	return nil
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
		jobTemplate:     loadTemplate("./job.tmpl")}

	summary := http_common.ErrorHandler(ctx.summary)
	http.HandleFunc("/", summary)
	http.HandleFunc("/summary", summary)
	http.HandleFunc("/job/", http_common.ErrorHandler(ctx.job))
	http.HandleFunc("/logs/", http_common.ErrorHandler(ctx.logs))

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

// https://gowalker.org/k8s.io/api/batch/v1#JobList
func summary(client *k8.Client) (summaryTemplate, error) {
	var jobs []jobSummary
	jobList, err := client.ListJobs()
	if err != nil {
		return summaryTemplate{}, err
	}
	for _, item := range jobList.Items {
		job := jobSummary{
			Name:           item.Name,
			Status:         k8.Status(item.Status),
			StartTime:      k8.ToString(item.Status.StartTime),
			CompletionTime: k8.ToString(item.Status.CompletionTime),
			Duration:       k8.Duration(item.Status.StartTime, item.Status.CompletionTime),
			Pod:            summaryForPod(client, item.Name)}
		jobs = append(jobs, job)
	}
	return summaryTemplate{
		Namespace: client.Namespace,
		Jobs:      jobs}, nil
}

// Pod Status
// Last Pod State (type/reason/message)
func summaryForPod(client *k8.Client, jobName string) podSummary {
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

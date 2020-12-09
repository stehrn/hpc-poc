package main

import (
	"fmt"
	"net/http"

	k8 "github.com/stehrn/hpc-poc/kubernetes"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// summaryTemplate data to render into summary template
type summaryTemplate struct {
	Namespace           string
	BusinessNameOptions []BusinessNameOptions
	Jobs                []jobSummary
	Page                string
}

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

func (ctx *handlerContext) SummaryHandler(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		summary := summaryTemplate{
			Namespace:           ctx.client.Namespace,
			BusinessNameOptions: NewBusinessNameOptions(""),
			Page:                "summary"}
		ctx.summaryTemplate.Execute(w, summary)
	case "POST":
		business := r.FormValue("business")
		summary, err := summary(business, ctx.client)
		if err != nil {
			return fmt.Errorf("Error creating summary page: %v", err)
		}
		ctx.summaryTemplate.Execute(w, summary)
	}
	return nil
}

func summary(business string, client *k8.Client) (summaryTemplate, error) {
	var jobs []jobSummary
	options := metav1.ListOptions{LabelSelector: fmt.Sprintf("business=%s", business)}
	jobList, err := client.ListJobs(options)
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
		Namespace:           client.Namespace,
		BusinessNameOptions: NewBusinessNameOptions(business),
		Jobs:                jobs}, nil
}

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

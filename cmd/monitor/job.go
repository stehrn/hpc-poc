package main

import (
	"errors"
	"net/http"
	"strings"

	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
)

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

func (ctx *handlerContext) JobHandler(w http.ResponseWriter, r *http.Request) error {
	jobName := strings.TrimPrefix(r.URL.Path, "/job/")
	if jobName == "" {
		return errors.New("No job sepcified")
	}

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

// called from job.tmpl
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

// called from job.tmpl
func (j myJob) ContainerEnv() []apiv1.EnvVar {
	return j.Spec.Template.Spec.Containers[0].Env
}

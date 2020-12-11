package main

import (
	"errors"
	"net/http"
	"strings"

	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"

	k8 "github.com/stehrn/hpc-poc/kubernetes"
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

func (ctx *handlerContext) JobHandler(w http.ResponseWriter, r *http.Request) error {
	jobName := strings.TrimPrefix(r.URL.Path, "/job/")
	if jobName == "" {
		return errors.New("No job sepcified")
	}

	job, err := ctx.client.FindJob(jobName)
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
func (j myJob) JobState() k8.JobState {
	return k8.GetJobState(j.Status)
}

// called from job.tmpl
func (j myJob) ContainerEnv() []apiv1.EnvVar {
	return j.Spec.Template.Spec.Containers[0].Env
}

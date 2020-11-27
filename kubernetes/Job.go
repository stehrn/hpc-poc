package kubernetes

import (
	"log"

	"github.com/pkg/errors"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	// To support connecting to GKE from outside of cluster (if KUBE_CONFIG used)
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

// ListJobs list all jobs
// calls List(opts metav1.ListOptions) (*v1.JobList, error)
func (c Client) ListJobs() (*batchv1.JobList, error) {
	result, err := c.jobsClient().List(metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "falied to list jobs")
	}
	return result, nil
}

// CreateJob create a k8 job
func (c Client) CreateJob(info JobInfo) (*batchv1.Job, error) {
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      info.Name,
			Namespace: c.Namespace,
		},
		Spec: batchv1.JobSpec{
			Template: apiv1.PodTemplateSpec{
				Spec: apiv1.PodSpec{
					RestartPolicy: "OnFailure",
					Containers: []apiv1.Container{
						{
							Name:  "engine",
							Image: info.Image,
							Env: []apiv1.EnvVar{
								apiv1.EnvVar{
									Name:  "PAYLOAD",
									Value: info.PayLoad,
								},
							},
						},
					},
				},
			},
		},
	}

	result, err := c.jobsClient().Create(job)
	if err != nil {
		return nil, errors.Wrapf(err, "falied to create job with %#v", info)
	}

	log.Printf("Created job %q.\n", result.Name)
	return result, nil
}

// Job load job from job name
func (c Client) Job(jobName string) (*batchv1.Job, error) {
	result, err := c.jobsClient().Get(jobName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get job: '%s'", jobName)
	}
	return result, nil
}

// Status status of job - assumes we only have 1 job
func Status(status batchv1.JobStatus) string {
	log.Printf("Status is: %v", status)
	if status.Active > 0 {
		return "Running"
	} else if status.Succeeded > 0 {
		return "Successful"
	} else if status.Failed > 0 {
		return "Failed"
	}
	return "Unkonwn"
}

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

// JobWatcher watch a job
type JobWatcher interface {
	Watch(filter metav1.ListOptions, predicate func(status batchv1.JobStatus) bool, callback func(job *batchv1.Job)) error
}

// JobInterface has methods to work with Jobs
type JobInterface interface {
	FindJob(jobName string) (*batchv1.Job, error)
	ListJobs(options metav1.ListOptions) (*batchv1.JobList, error)
	CreateJob(options JobOptions) (*batchv1.Job, error)
	JobWatcher
}

// JobOptions details of job to create
type JobOptions struct {
	Name        string
	Image       string
	Parallelism int32
	Labels      map[string]string
	Env         []apiv1.EnvVar
}

// FindJob load Job for given job name
func (c Client) FindJob(jobName string) (*batchv1.Job, error) {
	result, err := c.jobsClient().Get(jobName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get job: '%s'", jobName)
	}
	return result, nil
}

// ListJobs list all jobs
func (c Client) ListJobs(options metav1.ListOptions) (*batchv1.JobList, error) {
	result, err := c.jobsClient().List(options)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to list jobs")
	}
	return result, nil
}

// CreateJob create a new Job
func (c Client) CreateJob(options JobOptions) (*batchv1.Job, error) {
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      options.Name,
			Namespace: c.Namespace(),
			Labels:    options.Labels,
		},
		Spec: batchv1.JobSpec{
			Parallelism: &options.Parallelism,
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: options.Labels,
				},
				Spec: apiv1.PodSpec{
					RestartPolicy: "Never",
					Volumes: []apiv1.Volume{
						{
							Name: "google-cloud-key",
							VolumeSource: apiv1.VolumeSource{
								Secret: &apiv1.SecretVolumeSource{
									SecretName: "pubsub-acc-key",
								},
							},
						},
					},
					Containers: []apiv1.Container{
						{
							Name:  "engine",
							Image: options.Image,
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      "google-cloud-key",
									MountPath: "/var/secrets/google",
									ReadOnly:  true,
								},
							},
							Env: options.Env,
						},
					},
				},
			},
		},
	}

	result, err := c.jobsClient().Create(job)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create job with options %#v", options)
	}
	return result, nil
}

// Watch watch jobs
func (c Client) Watch(filter metav1.ListOptions, predicate func(status batchv1.JobStatus) bool, callback func(job *batchv1.Job)) error {
	watch, err := c.jobsClient().Watch(filter)
	if err != nil {
		return errors.Wrapf(err, "Failed to watch jobs with filter: '%v'", filter)
	}
	go func() {
		for event := range watch.ResultChan() {
			job, ok := event.Object.(*batchv1.Job)
			if !ok {
				log.Panicf("Unexpected type: '%v'", event.Type)
			}
			if predicate(job.Status) {
				callback(job)
			}
		}
	}()
	return nil
}

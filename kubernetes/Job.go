package kubernetes

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/stehrn/hpc-poc/gcp/storage"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	// To support connecting to GKE from outside of cluster (if KUBE_CONFIG used)

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

// JobOptions details of job to create
type JobOptions struct {
	Name   string
	Image  string
	Labels map[string]string
	storage.Location
}

// ListJobs list all jobs
func (c Client) ListJobs(options metav1.ListOptions) (*batchv1.JobList, error) {
	result, err := c.jobsClient().List(options)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to list jobs")
	}
	return result, nil
}

// CreateJob create a k8 job
func (c Client) CreateJob(options JobOptions) (*batchv1.Job, error) {
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      options.Name,
			Namespace: c.Namespace,
			Labels:    options.Labels,
		},
		Spec: batchv1.JobSpec{
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
							Env: []apiv1.EnvVar{
								{
									Name:  "GOOGLE_APPLICATION_CREDENTIALS",
									Value: "/var/secrets/google/key.json",
								},
								{
									Name:  "CLOUD_STORAGE_BUCKET_NAME",
									Value: options.Bucket,
								},
								{
									Name:  "CLOUD_STORAGE_OBJECT_NAME",
									Value: options.Object,
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
		return nil, errors.Wrapf(err, "Failed to create job with options %#v", options)
	}
	return result, nil
}

// Watch watch jobs
func (c Client) Watch(jobSucceeded func(job *batchv1.Job)) error {

	label := ""
	watch, err := c.jobsClient().Watch(metav1.ListOptions{})
	// TODO: add specific label
	// metav1.ListOptions{
	// 	LabelSelector: label,
	// }
	if err != nil {
		return errors.Wrapf(err, "Failed to watch jobs with label: '%s'", label)
	}
	go func() {
		for event := range watch.ResultChan() {
			job, ok := event.Object.(*batchv1.Job)
			if !ok {
				log.Panicf("Unexpected type: '%v'", event.Type)
			}
			if job.Status.Succeeded == 1 {
				jobSucceeded(job)
			}
		}
	}()
	return nil
}

// Job load job from job name
func (c Client) Job(jobName string) (*batchv1.Job, error) {
	result, err := c.jobsClient().Get(jobName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get job: '%s'", jobName)
	}
	return result, nil
}

// Status status of job - assumes we only have 1 job
func Status(status batchv1.JobStatus) string {
	if status.Active > 0 {
		return "Running"
	} else if status.Succeeded > 0 {
		return "Successful"
	} else if status.Failed > 0 {
		if len(status.Conditions) > 0 {
			return fmt.Sprintf("Failed (%s)", status.Conditions[0].Reason)
		}
		return "Failed"
	}
	return "Unkonwn"
}

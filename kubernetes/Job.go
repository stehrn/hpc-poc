package kubernetes

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/pkg/errors"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/batch/v1"

	// To support connecting to GKE from outside of cluster (if KUBE_CONFIG used)
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
)

// Client can be used to create and list kubernete Jobs
type Client struct {
	namespace string
	clientSet *kubernetes.Clientset
}

// JobCreate details of job to create
type JobCreate struct {
	jobName string
	image   string
	payLoad string
}

// NewClient create Client
func NewClient(namespace string) *Client {
	return &Client{namespace, clientset()}
}

// create Job Batch client
func (c Client) jobsClient() v1.JobInterface {
	return c.clientSet.BatchV1().Jobs(c.namespace)
}

func clientset() *kubernetes.Clientset {
	kubeConfig := os.Getenv("KUBE_CONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		log.Fatalf("Could not create k8 client: %v", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Could not create clientset: %v", err)
	}
	return clientset
}

// ListJobs list all jobs
// calls List(opts metav1.ListOptions) (*v1.JobList, error)
func (c Client) ListJobs() (*batchv1.JobList, error) {
	result, err := c.jobsClient().List(metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "falied to list jobs")
	}
	return result, nil
}

// CreateJob create a kubernetes job
func (c Client) CreateJob(info JobCreate) (*batchv1.Job, error) {
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      info.jobName,
			Namespace: c.namespace,
		},
		Spec: batchv1.JobSpec{
			Template: apiv1.PodTemplateSpec{
				Spec: apiv1.PodSpec{
					RestartPolicy: "OnFailure",
					Containers: []apiv1.Container{
						{
							Name:  "engine",
							Image: info.image,
							Env: []apiv1.EnvVar{
								apiv1.EnvVar{
									Name:  "PAYLOAD",
									Value: info.payLoad,
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

// Pod load Pod for given job name
// For now, we just expect 1 pod per job
func (c Client) Pod(jobName string) (apiv1.Pod, error) {
	listOptions := metav1.ListOptions{
		LabelSelector: "job-name=" + jobName,
	}
	pods, err := c.clientSet.CoreV1().Pods(c.namespace).List(listOptions)
	if err != nil {
		return apiv1.Pod{}, errors.Wrapf(err, "failed to get pod from job: '%s'", jobName)
	}

	if len(pods.Items) != 0 {
		return apiv1.Pod{}, fmt.Errorf("Expected 1 pod, got %d", len(pods.Items))
	}
	return pods.Items[0], nil
}

// Logs get logs for pod
func (c Client) Logs(pod apiv1.Pod) (string, error) {
	req := c.clientSet.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &apiv1.PodLogOptions{})
	podLogs, err := req.Stream()
	if err != nil {
		return "", errors.Wrap(err, "falied to get job logs: error opening stream")
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "", errors.Wrap(err, "falied to get job logs: error in copy information from podLogs to buf")
	}
	return buf.String(), nil
}

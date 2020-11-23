package kubernetes

import (
	"log"
	"os"

	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/batch/v1"
	"k8s.io/client-go/tools/clientcmd"
)

// JobService can be used to create and list kubernete Jobs
type JobService struct {
	namspace     string
	jobInterface v1.JobInterface
}

// JobCreate details of job to create
type JobCreate struct {
	jobName string
	image   string
	payLoad string
}

// New create JobService
func New(namspace string) *JobService {
	return &JobService{namspace, client(namspace)}
}

// Client creates a Job Batch client
func client(namspace string) v1.JobInterface {
	return clientset().BatchV1().Jobs(namspace)
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
func (j JobService) ListJobs() *batchv1.JobList {
	result, err := j.jobInterface.List(metav1.ListOptions{})
	if err != nil {
		// TODO: handle this better, we clealy dont want to exit here
		// throw exception intead, or panic?
		log.Fatalf("Could not create job: %v", err)
	}
	return result
}

// CreateJob create a kubernetes job
func (j JobService) CreateJob(info JobCreate) {
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      info.jobName,
			Namespace: j.namspace,
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

	result, err := j.jobInterface.Create(job)
	if err != nil {
		// TODO: handle this better, we clealy dont want to exit here
		// throw exception intead, or panic?
		log.Fatalf("Could not create job: %v", err)
	}
	log.Printf("Created job %q.\n", result.Name)
}

package kubernetes

import (
	"log"

	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/batch/v1"
	"k8s.io/client-go/tools/clientcmd"
)

// JobService can be used to create and list kunernete Jobs
type JobService struct {
	namspace     string
	jobInterface v1.JobInterface
}

// New create JobService
func New(namspace string) *JobService {
	return &JobService{namspace, client(namspace)}
}

// Client creates a Job Batch client
func client(namspace string) v1.JobInterface {
	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		log.Fatalf("Could not create k8 client: %v", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Could not create clientset: %v", err)
	}
	return clientset.BatchV1().Jobs(namspace)
}

// CreateJob create a kubernetes job
func (j JobService) CreateJob(jobName string, image string, payLoad string) {
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: j.namspace,
		},
		Spec: batchv1.JobSpec{
			Template: apiv1.PodTemplateSpec{
				Spec: apiv1.PodSpec{
					RestartPolicy: "OnFailure",
					Containers: []apiv1.Container{
						{
							Name:  "engine",
							Image: image,
							Env: []apiv1.EnvVar{
								apiv1.EnvVar{
									Name:  "PAYLOAD",
									Value: payLoad,
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
		log.Fatalf("Could not create job: %v", err)
	}
	log.Printf("Created job %q.\n", result.Name)
}

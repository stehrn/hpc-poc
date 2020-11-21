package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/batch/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	log.Print("Starting orchestrator")

	namespace := env("NAMSPACE")
	engineImage := env("ENGINE_IMAGE")

	log.Printf("Creating jobs client for namespace %s (job will use image: %s)", namespace, engineImage)
	jobsClient := createJobClient(namespace)

	subscribe(func(ctx context.Context, m *pubsub.Message) {
		jobName := "engine-job-" + m.ID
		payload := string(m.Data)
		log.Printf("Got message: %s, creating Job: %s", payload, jobName)
		createJob(jobsClient, namespace, jobName, engineImage, payload)
		m.Ack()
	})
}

func subscribe(callback func(ctx context.Context, m *pubsub.Message)) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, "hpc-poc")
	if err != nil {
		log.Fatalf("Could not create pubsub Client: %v", err)
	}

	subName := env("SUBSCRIPTION_NAME")
	log.Printf("Subscribing to %s", subName)
	sub := client.Subscription(subName)
	err = sub.Receive(context.Background(), callback)
	if err != nil {
		log.Fatalf("Could not receive message: %v", err)
	}
}

func createJobClient(namspace string) v1.JobInterface {
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

func createJob(jobsClient v1.JobInterface, namspace string, jobName string, image string, payLoad string) {
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: namspace,
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

	result, err := jobsClient.Create(job)
	if err != nil {
		log.Fatalf("Could not create job: %v", err)
	}
	log.Printf("Created job %q.\n", result.Name)
}

func env(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("No '%s' env variable", key)
	}
	return value
}

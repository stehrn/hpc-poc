//
package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"cloud.google.com/go/pubsub"
	cloudStorage "cloud.google.com/go/storage"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/stehrn/hpc-poc/client"
	"github.com/stehrn/hpc-poc/executor"
	messaging "github.com/stehrn/hpc-poc/gcp/pubsub"
	"github.com/stehrn/hpc-poc/gcp/storage"
	"github.com/stehrn/hpc-poc/internal/utils"
	"github.com/stehrn/hpc-poc/kubernetes"
	k8 "github.com/stehrn/hpc-poc/kubernetes"
)

// Global API clients used across function invocations.
var (
	gcpContext     *executor.GcpContext
	k8Client       k8.ClientInterface
	subClient      *messaging.Client
	storageClient  storage.ClientInterface
	business       string
	taskLoadFactor float64
	maxPodsPerJob  int64
)

func init() {
	gcpContext = executor.NewGcpContextFromEnv()
}

// init sub client
func init() {
	var err error
	subClient, err = gcpContext.NewSubClient()
	if err != nil {
		log.Fatalf("Could not create gcp sub client: %v", err)
	}
}

// init storage client
func init() {
	var err error
	storageClient, err = gcpContext.NewStorageClient()
	if err != nil {
		log.Fatalf("Could not create gcp storage client: %v", err)
	}
}

// init k8 client
func init() {
	var err error
	k8Client, err = k8.NewClient(gcpContext.Namespace)
	if err != nil {
		log.Fatalf("Could not create k8 client: %v", err)
	}
}

func main() {
	log.Print("Starting orchestrator")
	startJobWatcher()
	subscribe()
}

func subscribe() {
	log.Print("Startng subscriber")
	imageRegistry := utils.Env("IMAGE_REGISTRY")

	err := subClient.Subscribe(func(ctx context.Context, m *pubsub.Message) {
		jobLocation, err := storage.ToLocation(m.Data)
		if err != nil {
			log.Printf("Could not get location from message data (%v), error: %v", m.Data, err)
			return
		}

		var imageName string
		var env []apiv1.EnvVar

		job, err := job(jobLocation)
		if err != nil {
			log.Printf("Could not load job data (tasks) at location %q, error: %v", jobLocation, err)
			return
		}

		numTasks := len(job.Tasks())
		parallelism := parallelism(numTasks)
		isSubscriptionBased := (parallelism == 1)

		if isSubscriptionBased {
			jobPub, err := publishToJobTopic(jobLocation, job)
			if err != nil {
				log.Printf("Could not publish task locations, error: %v", err)
				return
			}
			imageName = "engine_subscription"
			env = subscriptionEnv(jobLocation.Bucket, subClient.Project, jobPub.SubscriptionID())
		} else {
			imageName = "engine_storage"
			env = storageEnv(jobLocation)

		}

		engineImage := fmt.Sprintf("%s/%s:latest", imageRegistry, imageName)

		options := k8.JobOptions{
			Name:        "engine-job-" + m.ID,
			Image:       engineImage,
			Parallelism: parallelism,
			Labels:      labels(jobLocation, fmt.Sprint(numTasks), m.ID),
			Env:         env}
		log.Printf("Creating Job with options: %v", options)
		_, err = k8Client.CreateJob(options)
		if err != nil {
			log.Printf("Could not create job with options: %v, error: %v", options, err)

			if isSubscriptionBased {
				deleteJobSubscription(jobLocation)
			}
			return
		}
		m.Ack()
	})

	if err != nil {
		panic(err)
	}
}

// TODO: this blindly tries to delete stuff already deleted when service 1st run if Jobs still around
func startJobWatcher() {
	log.Print("Startng job watcher")
	var err error

	options := metav1.ListOptions{LabelSelector: fmt.Sprintf("business=%s", string(business))}
	err = k8Client.Watch(options, kubernetes.SUCCESS, func(job *batchv1.Job) {

		labels := Labels{job.Labels}
		jobLocation := storage.Location{
			Bucket: labels.storageBucket(),
			Object: labels.storageObject(),
		}

		deleteJobStorageObjects(jobLocation)
		deleteJobSubscription(jobLocation)
	})
	if err != nil {
		log.Fatal("Could not start watching jobs", err)
	}
}

func deleteJobStorageObjects(jobLocation storage.Location) {
	log.Printf("Deleting cloud storage job data at location: %v", jobLocation)
	err := storageClient.Delete(jobLocation)
	if err != nil {
		log.Printf("Failed to delete job object at location: %v, error: %v", jobLocation, err)
	}
}

func deleteJobSubscription(jobLocation storage.Location) {
	ID := pubSubID(jobLocation.Object)
	log.Printf("Deleting job pubsub ID %q", ID)
	err := subClient.ExistingTempPubSub(ID).Delete()
	if err != nil {
		log.Printf("Error deleting job pubsub ID %q, error: %v", ID, err)
	}
}

// create client Job from informaiton in Location, a new task is added to the job for each storage object at location
func job(location storage.Location) (client.Job, error) {

	objectPath, err := client.ParseObjectPath(location.Object)
	if err != nil {
		return nil, err
	}

	session := client.NewSession(objectPath.Session, business)
	job := session.NewJob(objectPath.Job)

	err = storageClient.ForEachObject(location, func(attrs *cloudStorage.ObjectAttrs) error {
		objectPath, err := client.ParseObjectPath(attrs.Name)
		if err != nil {
			return err
		}
		job.AddTask(client.NewTaskProxy(job, objectPath.Task))
		return nil
	})
	return job, err
}

func pubSubID(object string) string {
	return strings.ReplaceAll(object, "/", "-")
}

// create temp topic (and subscription for engines) and publish task locations
func publishToJobTopic(location storage.Location, job client.Job) (*messaging.TempPubSub, error) {

	ID := pubSubID(location.Object)

	tmpPubSub, err := subClient.NewTempPubSub(ID)
	if err != nil {
		return nil, err
	}

	var taskLocations [][]byte
	for _, task := range job.Tasks() {
		location := storage.NewLocation(location.Bucket, task.ObjectPath().String())
		bytes, err := location.ToBytes()
		if err != nil {
			return nil, err
		}
		taskLocations = append(taskLocations, bytes)
	}

	err = tmpPubSub.PublishMany(taskLocations)
	if err != nil {
		log.Printf("Error publishing locations, deleting tmp pubsub %s, error: %v", tmpPubSub.SubscriptionID(), err)
		err := tmpPubSub.Delete()
		return nil, err
	}
	log.Printf("Published %d messages to topic %q", len(taskLocations), tmpPubSub.TopicName)

	return tmpPubSub, nil
}

// storageEnv env variables for storage based engine
func storageEnv(location storage.Location) []apiv1.EnvVar {
	return []apiv1.EnvVar{
		{
			Name:  "GOOGLE_APPLICATION_CREDENTIALS",
			Value: "/var/secrets/google/key.json",
		},
		{
			Name:  "CLOUD_STORAGE_BUCKET_NAME",
			Value: location.Bucket,
		},
		{
			Name:  "CLOUD_STORAGE_OBJECT_NAME",
			Value: location.Object,
		}}
}

// subscriptionEnv env variables for subscription based engine
func subscriptionEnv(bucket, project, subscription string) []apiv1.EnvVar {
	return []apiv1.EnvVar{
		{
			Name:  "GOOGLE_APPLICATION_CREDENTIALS",
			Value: "/var/secrets/google/key.json",
		},
		{
			Name:  "CLOUD_STORAGE_BUCKET_NAME",
			Value: bucket,
		},
		{
			Name:  "PROJECT_NAME",
			Value: project,
		},
		{
			Name:  "SUBSCRIPTION_NAME",
			Value: subscription,
		}}
}

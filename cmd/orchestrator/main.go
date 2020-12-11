//
package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"strings"

	"cloud.google.com/go/pubsub"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/stehrn/hpc-poc/client"
	messaging "github.com/stehrn/hpc-poc/gcp/pubsub"
	"github.com/stehrn/hpc-poc/gcp/storage"
	"github.com/stehrn/hpc-poc/internal/utils"
	"github.com/stehrn/hpc-poc/kubernetes"
	k8 "github.com/stehrn/hpc-poc/kubernetes"
)

// Global API clients used across function invocations.
var (
	k8Client       k8.ClientInterface
	subClient      *messaging.Client
	storageClient  storage.ClientInterface
	business       client.Business
	taskLoadFactor float64
	maxPodsPerJob  int64
)

const defaultTaskLoadFactor = 0.2
const defaultMaxPodsPerJob = 100

// init k8 client
func init() {
	var err error
	k8Client, err = k8.NewEnvClient()
	if err != nil {
		log.Fatalf("Could not create k8 client: %v", err)
	}
}

// init sub client
func init() {
	var err error
	business = client.BusinessFromEnv()
	project := utils.Env("PROJECT_NAME")
	subscription := business.SubscriptionName(project)
	subClient, err = messaging.NewSubClient(project, subscription)
	if err != nil {
		log.Fatalf("Could not create gcp sub client: %v", err)
	}
}

// init storage client
func init() {
	var err error
	storageClient, err = storage.NewEnvClient()
	if err != nil {
		log.Fatalf("Could not create gcp storage client: %v", err)
	}
}

// init task load factor
func init() {
	taskLoadFactor = utils.EnvAsFloat("TASK_LOAD_FACTOR", defaultTaskLoadFactor)
	maxPodsPerJob = utils.EnvAsInt("MAX_PODS_PER_JOB", defaultMaxPodsPerJob)
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

		taskLocations, err := loadJobData(jobLocation)
		if err != nil {
			log.Printf("Could not load job data (tasks) at location %q, error: %v", jobLocation, err)
			return
		}

		numTasks := len(taskLocations)
		parallelism := parallelism(numTasks)

		if parallelism == 1 {
			imageName = "engine_storage"
			env = storageEnv(jobLocation)
		} else {
			ID := pubSubID(jobLocation.Object)
			jobPub, err := publishToJobTopic(ID, taskLocations)
			if err != nil {
				log.Printf("Could not publish task locations, error: %v", err)
				return
			}
			imageName = "engine_subscription"
			env = subscriptionEnv(storageClient.BucketName(), subClient.Project, jobPub.SubscriptionID())
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
		object := labels.storageObject()

		jobLocation := storage.Location{
			Bucket: labels.storageBucket(),
			Object: object,
		}
		log.Printf("Deleting cloud storage data at location: %v", jobLocation)
		err = storageClient.Delete(jobLocation)
		if err != nil {
			log.Printf("Failed to delete object at location: %v, error: %v", jobLocation, err)
		}
		ID := pubSubID(jobLocation.Object)
		log.Printf("Deleting job pubsub ID %q", ID)
		err := subClient.ExistingTempPubSub(ID).Delete()
		if err != nil {
			log.Printf("Error deleting job pubsub ID %q, error: %v", ID, err)
		}
	})
	if err != nil {
		log.Fatal("Could not start watching jobs", err)
	}
}

func parallelism(numTasks int) int32 {
	if numTasks == 1 {
		return 1
	}
	parallelism := int32(math.Max(float64(numTasks)*float64(taskLoadFactor), 1.0))
	log.Printf("Parallelism set to %d, (numtasks * taskLoadFactor) = (%d * %f)", parallelism, numTasks, taskLoadFactor)
	return parallelism
}

// do an 'ls" on job storage directory and create location for each, convert to bytes
func loadJobData(location storage.Location) ([][]byte, error) {
	// get slice of storage locations
	directory := strings.Trim(location.Object, "/")
	objects, err := storageClient.ListObjects(directory)
	if err != nil {
		return nil, err
	}
	data, err := storageClient.ToLocationByteSlice(objects)
	if err != nil {
		return nil, err
	}
	log.Printf("Loaded %d storage objects from %s", len(data), directory)
	return data, nil
}

func pubSubID(object string) string {
	return strings.ReplaceAll(object, "/", "-")
}

// create temp topic (and subscription for engines) and publish tasks locations
func publishToJobTopic(ID string, taskLocations [][]byte) (*messaging.TempPubSub, error) {
	tmpPubSub, err := subClient.NewTempPubSub(ID)
	if err != nil {
		return nil, err
	}

	err = tmpPubSub.PublishMany(taskLocations)
	if err != nil {
		log.Printf("Error publishing locations, deleting tmp pubsub %s, error: %v", tmpPubSub.SubscriptionID(), err)
		err := tmpPubSub.Delete()
		return nil, err
	}
	log.Printf("Published %d messages to topic %q", len(taskLocations), tmpPubSub.TopicName())

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

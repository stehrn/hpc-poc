//
package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
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
	k8Client       *k8.Client
	subClient      *messaging.Client
	storageClient  *storage.Client
	business       client.Business
	taskLoadFactor float64
)

const defaultTaskLoadFactor = 0.2

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
	factorEnv := os.Getenv("TASK_LOAD_FACTOR")
	if factorEnv != "" {
		var err error
		taskLoadFactor, err = strconv.ParseFloat(factorEnv, 64)
		if err != nil {
			log.Fatalf("Could not convert SCALING_FACTOR %s into float64: %v", factorEnv, err)
		}
		log.Printf("Task load factor set to %f", taskLoadFactor)
	} else {
		taskLoadFactor = defaultTaskLoadFactor
		log.Printf("Task load factor set to default value of %f", taskLoadFactor)
	}
}

func main() {
	log.Print("Starting orchestrator")
	startJobWatcher()
	subscribe()
}

// TODO: this blindly tries to delete stuff already deleted when service 1st run
func startJobWatcher() {
	log.Print("Startng job watcher")
	var err error

	options := metav1.ListOptions{LabelSelector: fmt.Sprintf("business=%s", string(business))}
	err = k8Client.Watch(options, kubernetes.SUCCESS, func(job *batchv1.Job) {
		location := storage.Location{
			Bucket: job.Labels["gcp.storage.bucket"],
			Object: reverseClean(job.Labels["gcp.storage.object"]),
		}
		log.Printf("Deleting cloud storage data at location: %v", location)
		err = storageClient.Delete(location)
		if err != nil {
			log.Printf("Failed to delete object at location: %v, error: %v", location, err)
		}
	})
	if err != nil {
		log.Fatal("Could not start watching jobs", err)
	}
}

func subscribe() {
	log.Print("Startng subscriber")
	imageRegistry := utils.Env("IMAGE_REGISTRY")

	err := subClient.Subscribe(func(ctx context.Context, m *pubsub.Message) {
		location, err := storage.ToLocation(m.Data)
		if err != nil {
			log.Printf("Could not get location from message data (%v), error: %v", m.Data, err)
			return
		}

		strategy := strategy(location)
		var engineImage string
		var numTasks, parallelism int32
		var env []apiv1.EnvVar
		if strategy == "storage" {
			env = storageEnv(location)
			numTasks = 1
			parallelism = 1
		} else {
			var subscriptionID string
			subscriptionID, numTasks, err = publishToTempTopic(location)
			if err != nil {
				log.Printf("Could not create subscription for location %q, error: %v", location, err)
				return
			}
			parallelism = int32(math.Max(float64(numTasks)*float64(taskLoadFactor), 1.0))
			log.Printf("Parallelism set to %d, (numtasks * taskLoadFactor) = (%d * %f)", parallelism, numTasks, taskLoadFactor)
			env = subscriptionEnv(storageClient.BucketName, subClient.Project, subscriptionID)
		}
		engineImage = fmt.Sprintf("%s/engine_%s:latest", imageRegistry, strategy)

		options := k8.JobOptions{
			Name:        "engine-job-" + m.ID,
			Image:       engineImage,
			Parallelism: parallelism,
			Labels:      labels(location, fmt.Sprint(numTasks), m.ID),
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

// TODO
func strategy(location storage.Location) string {
	if !strings.HasPrefix(location.Object, "bu1") {
		return "subscription"
	}
	return "storage"
}

func publishToTempTopic(location storage.Location) (string, int32, error) {

	// get slice of storage locations
	directory := strings.Trim(location.Object, "/")
	objects, err := storageClient.ListStorageObjects(directory)
	if err != nil {
		return "", 0, err
	}
	data, err := storageClient.ToLocationByteSlice(objects)
	if err != nil {
		return "", 0, err
	}
	log.Printf("Loaded %d storage objects from %s", len(data), directory)

	// create temp topic and subscription and publish locations
	ID := strings.ReplaceAll(location.Object, "/", "-")
	tmpPubSub, err := subClient.NewTempPubSub(ID)
	if err != nil {
		return "", 0, err
	}

	err = tmpPubSub.Publish(data)
	if err != nil {
		return "", 0, err
	}
	log.Printf("Published to topic %q", tmpPubSub.TopicName())

	return tmpPubSub.SubscriptionID(), int32(len(data)), nil
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

func labels(location storage.Location, tasks, messageID string) map[string]string {
	labels := make(map[string]string)
	labels["business"] = string(business)
	labels["task.count"] = string(tasks)
	labels["task.load.factor"] = fmt.Sprint(taskLoadFactor)
	labels["k8.namespace"] = k8Client.Namespace
	labels["gcp.storage.bucket"] = location.Bucket
	labels["gcp.storage.object"] = clean(location.Object)
	labels["gcp.pubsub.project"] = subClient.Project
	labels["gcp.pubsub.subscription"] = subClient.SubscriptionID
	labels["gcp.pubsub.subscription_id"] = messageID
	return labels
}

func clean(item string) string {
	return strings.ReplaceAll(item, "/", "_")
}

func reverseClean(item string) string {
	return strings.ReplaceAll(item, "_", "/")
}

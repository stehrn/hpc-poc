package main

import (
	"fmt"
	"strings"

	"github.com/stehrn/hpc-poc/gcp/storage"
)

// Labels Job labels
type Labels struct {
	labels map[string]string
}

func (l *Labels) storageBucket() string {
	return l.labels["gcp.storage.bucket"]
}
func (l *Labels) storageObject() string {
	return reverseClean(l.labels["gcp.storage.object"])
}

func labels(location storage.Location, taskCount, messageID string) map[string]string {
	labels := make(map[string]string)
	labels["business"] = string(business)
	labels["task.count"] = string(taskCount)
	labels["task.load.factor"] = fmt.Sprint(taskLoadFactor)
	labels["k8.namespace"] = k8Client.Namespace()
	labels["gcp.storage.bucket"] = location.Bucket
	labels["gcp.storage.object"] = clean(location.Object)
	labels["gcp.pubsub.project"] = subClient.Project
	labels["gcp.pubsub.subscription"] = subClient.SubscriptionID
	labels["gcp.pubsub.subscription.id"] = messageID
	return labels
}

func clean(item string) string {
	return strings.ReplaceAll(item, "/", "_")
}

func reverseClean(item string) string {
	return strings.ReplaceAll(item, "_", "/")
}

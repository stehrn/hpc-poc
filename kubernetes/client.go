package kubernetes

import (
	"log"
	"os"

	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/batch/v1"

	// To support connecting to GKE from outside of cluster (if KUBE_CONFIG used)
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
)

// Client can be used to create and list kubernete Jobs
type Client struct {
	Namespace string
	clientSet *kubernetes.Clientset
}

// JobInfo details of job to create
type JobInfo struct {
	Name    string
	Image   string
	PayLoad string
}

// NewClientFromEnvironment create Client
func NewClientFromEnvironment() (*Client, error) {
	namespace := os.Getenv("NAMSPACE")
	if namespace == "" {
		return nil, errors.New("Namspace required")
	}
	return NewClient(namespace)
}

// NewClient create Client
func NewClient(namespace string) (*Client, error) {
	log.Printf("Creating k8 client for namespace: %s", namespace)
	clientset, err := clientset()
	if err != nil {
		return nil, err
	}
	return &Client{namespace, clientset}, nil
}

// create Job Batch client
func (c Client) jobsClient() v1.JobInterface {
	return c.clientSet.BatchV1().Jobs(c.Namespace)
}

func clientset() (*kubernetes.Clientset, error) {
	kubeConfig := os.Getenv("KUBE_CONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return nil, errors.Wrap(err, "Could not create k8 client")
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "Could not create clientset")
	}
	return clientset, nil
}

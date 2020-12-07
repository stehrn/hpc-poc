package kubernetes

import (
	"bytes"
	"io"

	"github.com/pkg/errors"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PodStatus represents status of (last) pod in Job
type PodStatus struct {
	Condition apiv1.PodCondition
	IsError   bool
}

// LatestPod get latest pod
func (c Client) LatestPod(jobName string) (apiv1.Pod, error) {
	pods, err := c.Pods(jobName)
	if err != nil {
		return apiv1.Pod{}, err
	}
	if len(pods) > 0 {
		return pods[0], nil
	}
	return apiv1.Pod{}, nil
}

// Pods get Poda for given job name
func (c Client) Pods(jobName string) ([]apiv1.Pod, error) {
	listOptions := metav1.ListOptions{
		LabelSelector: "job-name=" + jobName,
	}
	pods, err := c.clientSet.CoreV1().Pods(c.Namespace).List(listOptions)
	if err != nil {
		return []apiv1.Pod{}, errors.Wrapf(err, "Failed to get pod from job: '%s'", jobName)
	}

	return pods.Items, nil
}

// LogsForPod get logs for pod
func (c Client) LogsForPod(podName string) (string, error) {

	req := c.clientSet.CoreV1().Pods(c.Namespace).GetLogs(podName, &apiv1.PodLogOptions{})
	podLogs, err := req.Stream()
	if err != nil {
		return "", errors.Wrapf(err, "Failed to load job logs for pod '%s': error opening stream", podName)
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to load job logs for pod '%s': error in copy information from podLogs to buf", podName)
	}
	return buf.String(), nil
}

// LogsForJob get logs for job name
func (c Client) LogsForJob(jobName string) (string, error) {
	pod, err := c.LatestPod(jobName)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to load logs for job '%s'", jobName)
	}
	log, err := c.LogsForPod(pod.Name)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to load logs for job '%s'", jobName)
	}
	return log, nil
}

// NewPodStatus create new NewPodStatus from pod
func NewPodStatus(pod apiv1.Pod) PodStatus {
	if pod.Name != "" {
		conditions := pod.Status.Conditions
		if len(conditions) != 0 {
			condition := conditions[0]
			var jobError bool
			if condition.Reason == "Unschedulable" {
				jobError = true
			}
			return PodStatus{condition, jobError}
		}
	}
	return EmptyPodStatus()
}

// EmptyPodStatus create EmptyPodStatus
func EmptyPodStatus() PodStatus {
	return PodStatus{apiv1.PodCondition{}, false}
}

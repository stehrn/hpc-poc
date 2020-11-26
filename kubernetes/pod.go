package kubernetes

import (
	"bytes"
	"fmt"
	"io"

	"github.com/pkg/errors"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Pod load Pod for given job name
// For now, we just expect 1 pod per job
func (c Client) Pod(jobName string) (apiv1.Pod, error) {
	listOptions := metav1.ListOptions{
		LabelSelector: "job-name=" + jobName,
	}
	pods, err := c.clientSet.CoreV1().Pods(c.Namespace).List(listOptions)
	if err != nil {
		return apiv1.Pod{}, errors.Wrapf(err, "failed to get pod from job: '%s'", jobName)
	}

	if len(pods.Items) != 1 {
		return apiv1.Pod{}, fmt.Errorf("Expected 1 pod for job '%s', but got %d", jobName, len(pods.Items))
	}
	return pods.Items[0], nil
}

// LogsForPod get logs for pod
func (c Client) LogsForPod(pod apiv1.Pod) (string, error) {

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

// LogsForJob get logs for job name
func (c Client) LogsForJob(jobName string) (string, error) {
	pod, err := c.Pod(jobName)
	if err != nil {
		return "", errors.Wrap(err, "failed to load logs")
	}
	log, err := c.LogsForPod(pod)
	if err != nil {
		return "", errors.Wrap(err, "failed to load logs")
	}
	return log, nil
}

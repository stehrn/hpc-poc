package kubernetes

import apiv1 "k8s.io/api/core/v1"

// PodStatus represents status of (last) pod in Job
type PodStatus struct {
	Condition apiv1.PodCondition
	IsError   bool
}

// NewPodStatus create PodStatus from pod
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


TODO:
* add visibility on pubsub queue (use this as starting point for monitor?)
* think about jobs + tasks, right now, 1 job == 1 task
* jobs and pods accumulate - when to delete?
  * If we delete them, we lose engine logs
* write test harness to submit lots of jobs  
* add proper client interface/API
* Define cutom roles https://cloud.google.com/iam/docs/creating-custom-roles#iam-custom-roles-testable-permissions-gcloud
* add check that sub exists
* show events  

Look into quotas: https://cloud.google.com/kubernetes-engine/quotas

Create many Jobs in a batch might place high load on the Kubernetes control plane



Type: <span class="badge badge-info">{{ $condition.Type }}</span><p>Status: <span class="badge badge-info">{{ $condition.Status }}</span><p>LastProbeTime: <span class="badge badge-info">{{ $condition.LastProbeTime }}</span><p>Reason: <span class="badge badge-info">{{ $condition.Reason }}</span><p>Message: <span class="badge badge-info">{{ $condition.Message }}</span>


// Status status of job - assumes we only have 1 job
func Status(status batchv1.JobStatus) string {
	if status.Active > 0 {
		return "Running"
	} else if status.Succeeded > 0 {
		return "Successful"
	} else if status.Failed > 0 {
		if len(status.Conditions) > 0 {
			return fmt.Sprintf("Failed (%s)", status.Conditions[0].Reason)
		}
		return "Failed"
	}
	return "Unkonwn"
}

 {{ $lastPod := .LastPod }}
      {{if $lastPod.IsError }}
      <div class="alert alert-danger alert-dismissible fade show" role="alert">
        {{ $lastPod.Condition.Reason }}
        <button type="button" class="close" data-dismiss="alert" aria-label="Close">
        <span aria-hidden="true">&times;</span>
        </button>
      </div>
      {{end}}



open census

client


var ErrNotFound = errors.New("not found")
        return nil, fmt.Errorf("%q: %w", name, ErrNotFound)






kubectl create deployment nginx --image nginx:latest
kubectl get pods 
kubectl expose service nginx --port 80  --type LoadBalancer
kubectl scale deployment nginx --replicas 3

gcloud config set compute/zone us-central1-a




TODO:
* add visibility on pubsub queue (use this as starting point for monitor?)
* think about jobs + tasks, right now, 1 job == 1 task
* jobs and pods accumulate - when to delete?
  * If we delete them, we lose engine logs
* write test harness to submit lots of jobs  
* add proper client interface/API
* Dont create too mant topics
* Define cutom roles https://cloud.google.com/iam/docs/creating-custom-roles#iam-custom-roles-testable-permissions-gcloud
* add check that sub exists
* show events  

Look into quotas: https://cloud.google.com/kubernetes-engine/quotas

Create many Jobs in a batch might place high load on the Kubernetes control plane


build engine_storage Dockerfile
build engine_subscription Dockerfile

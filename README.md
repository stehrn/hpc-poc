
# POC of work queue

* Creation and set-up of GCP resources
  * Init gcloud
  * Create [GKE](https://cloud.google.com/kubernetes-engine/docs/quickstart) 
  * Create cloud storage bucket 
  * Create topic and subscription
  * Create GCP Service Account 
* Container build (with [Cloud Build](https://cloud.google.com/cloud-build)) and GKE deploymentmentof:  
  * Orchestrator
  * Engine (no deploy)
  * Client 
  * Monitor
* Test
  * Submit jobs
  * View activity

# Init gcloud
Note, some of following commands require: `export PROJECT_NAME=<GCP project>`

```
gcloud config set project ${PROJECT_NAME}
gcloud config set compute/zone europe-west2-a
```

https://console.cloud.google.com

# Create GKE 
```
gcloud container clusters create hpc-poc --num-nodes=1
gcloud container clusters get-credentials hpc-poc
```

View workload in [console](https://console.cloud.google.com/kubernetes/workload/)

# Create cloud storage bucket 
Data will be written by client, read by engine, and deleted by orchestrator
```
gsutil mb -p ${PROJECT_NAME} -c STANDARD -l europe-west2 -b on gs://stehrn_hpc-poc
```

# Create (test) topic and subscription
```
gcloud pubsub topics create test-topic
gcloud pubsub subscriptions create sub-test --topic=test-topic
```
The name of the subscription is passed into orchestrator container via `SUBSCRIPTION_NAME` env varible in [deployment.yaml](orchestrator/deployment.yaml)

# Create GCP Service Accounts 
## Summary
Set up following:
* Client: 
  * create storage objects [`storage.objectCreator`]
  * publish [`pubsub.publisher`]

* Orchestrator:
  * subscribe [`pubsub.subscriber`]
  * view [`pubsub.viewer`] (to check subscription exists)
  * delete storage objects [`storage.objects.delete`]

* Engine:
  * read storage objects [`storage.objectViewer`]

(see https://cloud.google.com/kubernetes-engine/docs/tutorials/authenticating-to-cloud-platform)

```
// create service sccount 
gcloud iam service-accounts create gke-sub-acc@hpc-poc.iam.gserviceaccount.com --description="GKE subscription account" --display-name="gke-subscription"
gcloud iam service-accounts list

// add roles
gcloud projects add-iam-policy-binding hpc-poc --member=serviceAccount:gke-sub-acc@hpc-poc.iam.gserviceaccount.com --role=roles/pubsub.subscriber 
gcloud projects add-iam-policy-binding hpc-poc --member=serviceAccount:gke-sub-acc@hpc-poc.iam.gserviceaccount.com --role=roles/pubsub.publisher
gcloud projects add-iam-policy-binding hpc-poc --member=serviceAccount:gke-sub-acc@hpc-poc.iam.gserviceaccount.com --role=roles/pubsub.viewer
gcloud projects add-iam-policy-binding hpc-poc --member=serviceAccount:gke-sub-acc@hpc-poc.iam.gserviceaccount.com --role=roles/storage.objectCreator
gcloud projects add-iam-policy-binding hpc-poc --member=serviceAccount:gke-sub-acc@hpc-poc.iam.gserviceaccount.com --role=roles/storage.objectViewer

// list roles
gcloud projects get-iam-policy hpc-poc --flatten="bindings[].members" --format='table(bindings.role)' --filter="bindings.members:gke-sub-acc@hpc-poc.iam.gserviceaccount.com"
```

* [storage iam-permissions](https://cloud.google.com/storage/docs/access-control/using-iam-permissions)

## Get GCP JSON key and create k8 secret
The key is injected into container env variable `GOOGLE_APPLICATION_CREDENTIALS` (also needed if running app locally)

```
gcloud iam service-accounts keys create ${HOME}/key.json --iam-account gke-sub-acc@hpc-poc.iam.gserviceaccount.com 
kubectl create secret generic pubsub-acc-key --from-file=key.json=${HOME}/key.json
```
...where path to download is location of key file.

# Build and deploy containers
Submit build to [cloud-build](https://cloud.google.com/cloud-build), which stores image in the [container-registry](https://cloud.google.com/container-registry); see [build-and-deploy](https://cloud.google.com/run/docs/quickstarts/build-and-deploy) quickstart.

## Orchestrator
see [orchestrator/README.md](orchestrator/README.md)

## Engine
see [engine/README.md](engine/README.md)

## Client
see [client/README.md](client/README.md)

## Monitor
see [monitor/README.md](monitor/README.md)

# Test
## Submit jobs
Use [web client](client/README.md) or gcloud shell:
```
gcloud pubsub topics publish test-topic --message="engine payload 1"
```

## View activity
* Use [web monitor](monitor/README.md) to view job/engine logs
* View workload in [console](https://console.cloud.google.com/kubernetes/workload/)
* Use `kubectl`:
```
kubectl logs --selector=app=orchestrator --tail 100
```

# Articles 
* https://blog.meain.io/2019/accessing-kubernetes-api-from-pod/
* https://cloud.google.com/pubsub/docs/quickstart-cli
* https://cloud.google.com/appengine/docs/flexible/go/writing-and-responding-to-pub-sub-messages
* https://github.com/googleapis/google-cloud-go/blob/master/pubsub/example_test.go
* https://pkg.go.dev/cloud.google.com/go/pubsub#example-Client.CreateSubscription
* https://cloud.google.com/run/docs/tutorials/pubsub
* https://cloud.google.com/kubernetes-engine/docs/tutorials/authenticating-to-cloud-platform

Argo looks interesting ....
https://github.com/argoproj/argo/blob/master/workflow/controller/workflowpod.go

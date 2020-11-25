
# POC of work queue

* Init gcloud
* Create [GKE](https://cloud.google.com/kubernetes-engine/docs/quickstart) 
* Create (test) topic and subscription
* Create GCP Service Account (so orchestrator container can subscribe to GCP)
* Build (with [Cloud Build](https://cloud.google.com/cloud-build)) and deployment  
  * Build and deploy orchestrator container into GKE 
  * Build engine container
  * Build and deploy monitor container into GKE 
* Testing

# Init gcloud
```
gcloud config set project hpc-poc
gcloud config set compute/zone europe-west2-a
```

https://console.cloud.google.com

# Create GKE 
Run:
```
gcloud container clusters create hpc-poc --num-nodes=1
gcloud container clusters get-credentials hpc-poc
```

View workload in [console](https://console.cloud.google.com/kubernetes/workload/)

# Create (test) topic and subscription
```
gcloud pubsub topics create test-topic
gcloud pubsub subscriptions create sub-test --topic=test-topic
```
The name of the subscription is passed into orchestrator container via `SUBSCRIPTION_NAME` env varible in [deployment.yaml](orchestrator/deployment.yaml)

# Create GCP Service Account (so orchestrator container can subscribe)
Create a service account to allow container running on GKP to subscribe.

see https://cloud.google.com/kubernetes-engine/docs/tutorials/authenticating-to-cloud-platform

## Create Service Account
```
gcloud iam service-accounts list
gcloud iam service-accounts create gke-sub-acc@hpc-poc.iam.gserviceaccount.com --description="GKE subscription account" --display-name="gke-subscription"
```
## Add `pubsub.subscriber` role
```
gcloud projects add-iam-policy-binding hpc-poc --member=serviceAccount:gke-sub-acc@hpc-poc.iam.gserviceaccount.com --role=roles/pubsub.subscriber
```
## Get JSON key and create k8 secret
The key is injected into container env variable `GOOGLE_APPLICATION_CREDENTIALS` 

```
gcloud iam service-accounts keys create key.json --iam-account gke-sub-acc@hpc-poc.iam.gserviceaccount.com 
kubectl create secret generic pubsub-acc-key --from-file=key.json=/Users/db/key.json
```
...where path to download is location of key file.

# Build containers
Submit build to [cloud-build](https://cloud.google.com/cloud-build), which stores image in the [container-registry](https://cloud.google.com/container-registry); see [build-and-deploy](https://cloud.google.com/run/docs/quickstarts/build-and-deploy) quickstart.

## Build and deploy orchestrator container into GKE 
see [orchestrator/README.md](orchestrator/README.md)

## Build engine container
see [engine/README.md](engine/README.md)

## Build and deploy monitor container into GKE 
see [monitor/README.md](monitor/README.md)

# Test
```
gcloud pubsub topics publish test-topic --message="engine payload 1"
kubectl logs --selector=app=orchestrator --tail 100
```

Use monitor to view job/engine logs, altenratively, view workload in [console](https://console.cloud.google.com/kubernetes/workload/)


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

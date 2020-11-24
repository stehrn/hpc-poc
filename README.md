
# POC of work queue

https://console.cloud.google.com

* Create GKE 
* Create GCP Service Account (so orchestrator container can subscribe to GCP)
* Create (test) topic and subscription
* Build and deploy orchestrator container into GKE 
* Build and deploy monitor container into GKE 
* Build engine container

# Create GCP Service Account so orchestrator container can subscribe
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

# Create (test) topic and subscription
```
gcloud pubsub topics create test-topic
gcloud pubsub subscriptions create sub-test --topic=test-topic
```
The name of the subscription is passed into orchestrator container via `SUBSCRIPTION_NAME` env varible in [deployment.yaml](orchestrator/deployment.yaml)

# Build and deploy orchestrator container into GKE 
see [orchestrator/README.md](orchestrator/README.md)

# Build and deploy monitor container into GKE 
see [monitor/README.md](monitor/README.md)

# Build engine container
see [engine/README.md](engine/README.md)

# Test
```
gcloud pubsub topics publish test-topic --message="engine payload 1"
kubectl logs --selector=app=orchestrator --tail 100
```

Use monitor to view job/engine logs, altenratively, view workload in [console](https://console.cloud.google.com/kubernetes/workload/)


# Articles that helped
* https://blog.meain.io/2019/accessing-kubernetes-api-from-pod/
* https://cloud.google.com/pubsub/docs/quickstart-cli
* https://cloud.google.com/appengine/docs/flexible/go/writing-and-responding-to-pub-sub-messages
* https://github.com/googleapis/google-cloud-go/blob/master/pubsub/example_test.go
* https://pkg.go.dev/cloud.google.com/go/pubsub#example-Client.CreateSubscription
* https://cloud.google.com/run/docs/tutorials/pubsub
* https://cloud.google.com/kubernetes-engine/docs/tutorials/authenticating-to-cloud-platform

Argo looks interesting ....
https://github.com/argoproj/argo/blob/master/workflow/controller/workflowpod.go

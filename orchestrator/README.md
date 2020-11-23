## Pub / Sub 

https://console.cloud.google.com


# Links
https://cloud.google.com/pubsub/docs/quickstart-cli
https://cloud.google.com/appengine/docs/flexible/go/writing-and-responding-to-pub-sub-messages
https://github.com/googleapis/google-cloud-go/blob/master/pubsub/example_test.go
https://pkg.go.dev/cloud.google.com/go/pubsub#example-Client.CreateSubscription
https://cloud.google.com/run/docs/tutorials/pubsub
https://cloud.google.com/kubernetes-engine/docs/tutorials/authenticating-to-cloud-platform


# Create Service Account to GKE can subscribe
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

## Create topic and subscription
```
gcloud pubsub topics create test-topic
gcloud pubsub subscriptions create sub-test --topic=test-topic
```
The name of the subscription is passed into container via `SUBSCRIPTION_NAME`

# Get valid version of libs
Get version of k8 
```
kubectl version
```
Add correct libs to `go.mod` (here, version is '1.16.3'):
```
go get k8s.io/client-go@kubernetes-1.16.3
```
(run from /Users/db/mygo/src/github.com/stehrn/hpc-poc/pub-sub)

# Create roles and service account and binding
To give pod access to k8 API and be able to create jobs
```
kubectl apply -f ./yaml
```

# Build container image and deploy
cd into app folder (e.g. /Users/db/mygo/src/github.com/stehrn/hpc-poc/pub-sub),
```
gcloud builds submit --tag gcr.io/hpc-poc/pub
kubectl apply -f pubsub-with-secret.yaml
```

gcloud container images list --repository=gcr.io/hpc-poc
gcloud container images list-tags gcr.io/hpc-poc/pub

# Create deployment
```
[kubectl delete deployment pubsub]
kubectl create deployment pubsub --image=gcr.io/hpc-poc/pub:latest
kubectl patch deployment pubsub -p '{"spec":{"template":{"spec":{"serviceAccountName":"job-engine"}}}}'
kubectl get deployment pubsub -o yaml 
```

View workload in [console](https://console.cloud.google.com/kubernetes/workload/), or:
```
kubectl get deployment hpc-server -o yaml
kubectl logs --selector=app=pubsub --tail 100

# Test
```
gcloud pubsub topics publish test-topic --message="hello nik?"
kubectl logs --selector=app=pubsub --tail 100
```



xxx
gcloud builds submit --tag gcr.io/hpc-poc/orchestrator
[kubectl delete deployment orchestrator]
kubectl apply -f yaml/deployment.yaml
kubectl logs --selector=app=orchestrator --tail 100

gcloud pubsub topics publish test-topic --message="engine job 1"


Check engines logs
kubectl get jobs

... engine-job-1734219594854205
.. but we need the pod .. so
engine-job-1734219594854205-kjcvc
kubectl logs engine-job-1734219594854205-kjcvc

Argo looks interesting ....
https://github.com/argoproj/argo/blob/master/workflow/controller/workflowpod.go
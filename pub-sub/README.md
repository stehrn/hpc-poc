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

# Build container image and deploy
```
gcloud builds submit --tag gcr.io/hpc-poc/pub
kubectl apply -f pubsub-with-secret.yaml
```

# Test
```
gcloud pubsub topics publish test-topic --message="hello nik"
kubectl logs --selector=app=pubsub --tail 100
```

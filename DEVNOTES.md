# Ad hoc notes


## Articles 
* https://blog.meain.io/2019/accessing-kubernetes-api-from-pod/
* https://cloud.google.com/pubsub/docs/quickstart-cli
* https://cloud.google.com/appengine/docs/flexible/go/writing-and-responding-to-pub-sub-messages
* https://github.com/googleapis/google-cloud-go/blob/master/pubsub/example_test.go
* https://pkg.go.dev/cloud.google.com/go/pubsub#example-Client.CreateSubscription
* https://cloud.google.com/run/docs/tutorials/pubsub
* https://cloud.google.com/kubernetes-engine/docs/tutorials/authenticating-to-cloud-platform
* Argo looks interesting : https://github.com/argoproj/argo/blob/master/workflow/controller/workflowpod.go

## Deploy to Cloud Run
```
gcloud run deploy --image gcr.io/hpc-poc/monitor --platform managed
```

## Go specific build
https://cloud.google.com/cloud-build/docs/building/build-go

Pushes go binary to storage bucket, not used for container

### Create storage bucket
https://cloud.google.com/storage/docs/creating-buckets
```
gsutil mb -p hpc-poc -l EUROPE-WEST2 -c STANDARD gs://nik-stehr-hpc
gsutil ls
gsutil ls -r gs://nik-stehr-hpc
```
### Submit build
```
gcloud builds submit --config cloudbuild-go.yaml go-app
```
[View build](https://console.cloud.google.com/cloud-build/builds/)






PROJECT_NUM=$(gcloud projects list --filter="$PROJECT" --format="value(PROJECT_NUMBER)" --project=$PROJECT)
SERVICE_ACCOUNT=${PROJECT_NUM}@cloudbuild.gserviceaccount.com
gcloud projects add-iam-policy-binding $PROJECT --member=serviceAccount:$SERVICE_ACCOUNT --role=roles/container.developer --project=$PROJECT

```


export PROJECT=hpc-poc
export CLUSTER=hpc-poc
export LOCATION=europe-west2-a
export PACKAGE=orchestrator

gcloud builds submit --project=$PROJECT --config=cloudbuild_with_deploy.yaml --substitutions=_PACKAGE=$PACKAGE,_GKE_CLUSTER=$CLUSTER,_GKE_LOCATION=$LOCATION .



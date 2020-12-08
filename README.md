
# POC of work queue

* Overview
* Creation and set-up of GCP resources
  * Init gcloud
  * Create [GKE](https://cloud.google.com/kubernetes-engine/docs/quickstart) 
  * Create cloud storage bucket 
  * Create topic and subscription
  * Create GCP Service Account 
* Build and deploy containers  
* Test
  * Submit jobs
  * View activity

# Overview 

## Flow
* `Client` submits a `Job`
  * write data to cloud storage bucket
  * publish message containing location of data on cloud storage (bucket/object)
* `Orchestrator` 
  * subscribes to topic
    * on message - create kubernetes/k8 job, passing in location of cloud storage data
  * watches (k8) jobs
    * on job success - delete cloud storage object
* `Engine`
   * read cloud storage data, do something with it, exit

Terms:
* `Client` - the thing submitting a `Job` (either the web client of driver app)
* `Job` - a unit of work 
  * can be broken down into tasks (for now  they represent same thing)
  * submitted for a given `Business`
  * has reference to data (in cloud storage) used by the `Engine`
* `Session` - a `Job` is submitted as part of a session
* `Business` - represents a given .. business (think desk)
  * each will have its own cloud storage sub-directory and topic
* `Orchestrator` - handles a `Job` for given `Business`
* `Engine` - the thing that computes something using the data referenced in the `Job`'

## Web applications
 * client - use to submit data 
 * monitor - view jobs and pods (extend to view pubsub details)

# Environment variables
Following env variables required for setup:
```
export CLIENT_NAME=client1
export BUSINESS_NAME=bu1
```
GCP resource specific:
```
export LOCATION=europe-west2
export ZONE=europe-west2-a
export PROJECT_NAME=hpc-poc
export CLOUD_STORAGE_BUCKET_NAME=${PROJECT_NAME}-bucket
export GKE_CLUSTER_NAME=${PROJECT_NAME}-cluster
```

# Init gcloud
```
gcloud config set project ${PROJECT_NAME}
gcloud config set compute/zone ${ZONE}
```

# Create GKE 
```
gcloud container clusters create ${GKE_CLUSTER_NAME} --num-nodes=1
gcloud container clusters get-credentials ${GKE_CLUSTER_NAME}
```

View workload in [console](https://console.cloud.google.com/kubernetes/workload/)

# Create cloud storage bucket 
Data will be written by client, read by engine, and deleted by orchestrator. 
```
gsutil mb -p ${PROJECT_NAME} -c STANDARD -l ${LOCATION} gs://${CLOUD_STORAGE_BUCKET_NAME}
```
The object/data will be stored in a business specific subdirectory: `${CLOUD_STORAGE_BUCKET_NAME}/${BUSINESS_NAME}/`
# Create topic and subscription (for given given business)
```
export TOPIC_NAME=${PROJECT_NAME}-${BUSINESS_NAME}-topic
export SUBSCRIPTION_NAME=${PROJECT_NAME}-${BUSINESS_NAME}-subscription
gcloud pubsub topics create ${TOPIC_NAME}
gcloud pubsub subscriptions create ${SUBSCRIPTION_NAME} --topic=${TOPIC_NAME}
```

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
  * publish [`roles/editor`]

* Engine:
  * read storage objects [`storage.objectViewer`]

* Monitor:
  * read storage objects [`storage.objectViewer`]


(see [authenticating-to-cloud-platform](https://cloud.google.com/kubernetes-engine/docs/tutorials/authenticating-to-cloud-platform))

For now, roles are v broad and grouped together, so more work needed here
```
export SERVICE_ACCOUNT=gke-sub-acc@hpc-poc.iam.gserviceaccount.com

// create service sccount 
gcloud iam service-accounts create ${SERVICE_ACCOUNT} --description="GKE subscription account" --display-name="gke-subscription"
gcloud iam service-accounts list

// add roles
gcloud projects add-iam-policy-binding ${PROJECT_NAME} --member=serviceAccount:${SERVICE_ACCOUNT} --role=roles/pubsub.subscriber 
gcloud projects add-iam-policy-binding ${PROJECT_NAME} --member=serviceAccount:${SERVICE_ACCOUNT} --role=roles/pubsub.publisher
gcloud projects add-iam-policy-binding ${PROJECT_NAME} --member=serviceAccount:${SERVICE_ACCOUNT} --role=roles/pubsub.viewer
gcloud projects add-iam-policy-binding ${PROJECT_NAME} --member=serviceAccount:${SERVICE_ACCOUNT} --role=roles/storage.objectAdmin
gcloud projects add-iam-policy-binding ${PROJECT_NAME} --member=serviceAccount:${SERVICE_ACCOUNT} --role=roles/pubsub.editor


// list roles
gcloud projects get-iam-policy ${PROJECT_NAME} --flatten="bindings[].members" --format='table(bindings.role)' --filter="bindings.members:${SERVICE_ACCOUNT}"
```

* [storage iam-permissions](https://cloud.google.com/storage/docs/access-control/using-iam-permissions)

## Get GCP JSON key and create k8 secret
The key is injected into container env variable `GOOGLE_APPLICATION_CREDENTIALS` (also needed if running app locally)

```
gcloud iam service-accounts keys create ${HOME}/key.json --iam-account ${SERVICE_ACCOUNT}
kubectl create secret generic pubsub-acc-key --from-file=key.json=${HOME}/key.json
```

# Build and deploy containers
Submit build to [cloud-build](https://cloud.google.com/cloud-build), which stores image in the [container-registry](https://cloud.google.com/container-registry)

[Build and deploy](BUILD_DEPLOY.md):
* orchestrator (service)
* engine (image only)
* client (web app)
* monitor (web app)

Further reading: 
* [build-and-deploy](https://cloud.google.com/run/docs/quickstarts/build-and-deploy) quickstart
# Test
## Submit jobs
Use web client or gcloud shell:
```
gcloud pubsub topics publish ${TOPIC_NAME} --message="engine payload 1"
```

## View activity
* Use web monito to view job/engine logs
* View workload in [console](https://console.cloud.google.com/kubernetes/workload/)
* Use `kubectl`:
```
kubectl logs --selector=app=orchestrator --tail 100
```

### List buckets:
```
gsutil ls -r gs://${CLOUD_STORAGE_BUCKET_NAME}/${BUSINESS_NAME}
```

### List jobs for given business:
```
kubectl get jobs -l business=${BUSINESS_NAME}
```
### List succesful jobs for given business:
```
kubectl get jobs -l business=${BUSINESS_NAME} --field-selector status.successful=1 
```
(or failed: `--field-selector status.successful=0`)


# Cleaup
In case you need to manually clear stuff down:
```
kubectl delete jobs --all
kubectl delete pods --all
gsutil rm gs://${CLOUD_STORAGE_BUCKET_NAME}/${BUSINESS_NAME}/*

```

gsutil ls -r gs://${CLOUD_STORAGE_BUCKET_NAME}/bu2*
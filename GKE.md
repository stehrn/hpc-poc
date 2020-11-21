

# Get started
1. create private github repo with go program
2. Build container with [Cloud Build](https://cloud.google.com/cloud-build)
3. deploy container to GKE

## Init gcloud
```
gcloud config set project hpc-poc
gcloud config set compute/zone europe-west2-a
```

## Build container
Submit build to [cloud-build](https://cloud.google.com/cloud-build), which stores image in the [container-registry](https://cloud.google.com/container-registry)

cd into app folder (e.g. src/github.com/stehrn/gcp/go-app), needs a Dockerfile and go source, and module (`go mod init`)

Submit build, giving name and location of container
```
gcloud builds submit --tag gcr.io/hpc-poc/server
```
List container (and tags):
```
gcloud container images list --repository=gcr.io/hpc-poc
gcloud container images list-tags gcr.io/hpc-poc/server
```

See [build-and-deploy](https://cloud.google.com/run/docs/quickstarts/build-and-deploy) quickstart.

## Deploy to Cloud Run
```
gcloud run deploy --image gcr.io/hpc-poc/server --platform managed
```

# Deploy container onto GKE
https://cloud.google.com/kubernetes-engine/docs/quickstart

```
gcloud container clusters create hpc-poc --num-nodes=1
gcloud container clusters get-credentials hpc-poc
kubectl create deployment hpc-server --image=gcr.io/hpc-poc/server:latest
```

View workload in [console](https://console.cloud.google.com/kubernetes/workload/), or:
```
kubectl get deployment hpc-server -o yaml
kubectl logs --selector=app=hpc-server --tail 100
```

# Destroy resources
```
kubectl delete service hpc-server

```

# Appendix
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
/Users/db/mygo/src/github.com/stehrn/gcp
gcloud builds submit --config cloudbuild-go.yaml go-app
```
[View build](https://console.cloud.google.com/cloud-build/builds/)

 https://hpc-poc-service-ygvnpjuaua-ew.a.run.app

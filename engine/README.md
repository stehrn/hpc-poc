# Engine container
A simple engine that downloads an object from GCP cloud storage bucket, prints it out, and exits.

Image is referenced in [orchestrator](../orchestrator/README.md) deployment (`ENGINE_IMAGE`)

The engine itself requires the following env variables:
* `BUCKET_NAME` 
* `OBJECT_NAME`
* `GOOGLE_APPLICATION_CREDENTIALS` to enable download payload from cloud storage

# Build
Run:
```
gcloud builds submit --config=cloudbuild.yaml --substitutions=_PACKAGE="engine" .
```
View image:
```
gcloud container images list --repository=gcr.io/hpc-poc
gcloud container images list-tags gcr.io/hpc-poc/engine
```

# Deploy
Engine image is not deployed as container into GKE, rather referenced in `Job` created by [orchestrator](../orchestrator/README.md)






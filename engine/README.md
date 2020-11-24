# Engine container
A simple engine, it actually just reads value of `PAYLOAD` env variable, prints it out, and exits.

Image is referenced in orchestrator deployment (`ENGINE_IMAGE`)

# Build
Run:
```
gcloud builds submit --tag gcr.io/hpc-poc/engine
```
View image:
```
gcloud container images list --repository=gcr.io/hpc-poc
gcloud container images list-tags gcr.io/hpc-poc/engine
```

# Deploy
Engine image is not deployed as container into GKE, rather referenced in Job created by orchestrator






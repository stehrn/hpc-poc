# Orchestrator container
Orchestrator:
* subscribes to subscription (`SUBSCRIPTION_NAME` env)- messages contain the location of a cloud storage object with data to send to engine
* on message, creates a k8 job with:
   * engine image (`ENGINE_IMAGE` env)
   * cloud storage object location
* set up a job watcher, when a job succeeds, delete the cloud storage object associated with the job

Note, some of following commands require: `export PROJECT_NAME=<GCP project>`

# Build (orchestrator)
Run:
```
gcloud builds submit --config=cloudbuild.yaml --substitutions=_PROJECT=${PROJECT_NAME},_PACKAGE="orchestrator" .
```
View image:
```
gcloud container images list --repository=gcr.io/${PROJECT_NAME}
gcloud container images list-tags gcr.io/${PROJECT_NAME}/orchestrator
```

# Deploy
Run:
```
kubectl apply -f orchestrator/yaml
kubectl get deployment orchestrator -o yaml 
```
Note the orchestrator container uses the k8 API to create Jobs, the required role, service account and bindings are also created.

# Monitor
View logs
```
kubectl logs --selector=app=orchestrator --tail 100
``` 
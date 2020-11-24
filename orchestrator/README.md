# Orchestrator container
Subscribes to configured subscription, when message recieved, create k8 job with engine image

# Build (orchestrator)
Run:
```
gcloud builds submit --tag gcr.io/hpc-poc/orchestrator
```
View image:
```
gcloud container images list --repository=gcr.io/hpc-poc
gcloud container images list-tags gcr.io/hpc-poc/orchestrator
```

# Deploy
Run:
```
cd orchestrator
kubectl apply -f /yaml
kubectl get deployment orchestrator -o yaml 
```
Note the orchestrator container uses the k8 API to create Jobs, the required role, service account and bindings are also created.

# Monitor
View logs
```
kubectl logs --selector=app=orchestrator --tail 100
``` 
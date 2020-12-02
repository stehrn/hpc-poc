# Monitor container
Simple web app to view jobs and logs 

Note, some of following commands require: `export PROJECT_NAME=<GCP project>`

# Build
Run (from base module dir):
```
gcloud builds submit --config=cloudbuild.yaml --substitutions=_PROJECT=${PROJECT_NAME},_PACKAGE="monitor",_DOCKERFILE="DockerfileForWeb" .
```
View image:
```
gcloud container images list --repository=gcr.io/${PROJECT_NAME}
gcloud container images list-tags gcr.io/${PROJECT_NAME}/monitor
```

# Deploy
Run:
```
kubectl apply -f monitor/yaml
kubectl get deployment monitor -o yaml 
```
Note the orchestrator container uses the k8 API to create Jobs, the required role, service account and bindings are also created.

Get port:
```
kubectl get service monitor
```
Open browswer at: http://<external-ip>:<port>/summary (use port number on left hand side of ':')

# Monitor
View logs
```
kubectl logs --selector=app=monitor --tail 100
``` 

# Delete
Run:
```
kubectl delete deployment monitor && kubectl delete pods,services -l app=monitor
```

# Run locally
```
export NAMESPACE=default
export KUBE_CONFIG=${HOME}/.kube/config
cd monitor
go run main.go
```
open http://localhost:8081/summary

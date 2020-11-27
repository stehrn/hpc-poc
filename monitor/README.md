# Monitor container
Simple web app to view jobs and logs 

# Build
Run (from base module dir):
```
gcloud builds submit --config cloudbuild_monitor.yaml
```
View image:
```
gcloud container images list --repository=gcr.io/hpc-poc
gcloud container images list-tags gcr.io/hpc-poc/monitor
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
Open browswer at: http://<external-ip>:<port>/jobs

(from above example - http://35.234.146.8:8081/jobs)

# Monitor
View logs
```
kubectl logs --selector=app=monitor --tail 100
``` 

# Delete
Run:
```
kubectl delete pods,services -l app=monitor
kubectl delete deployment monitor
```

# Run locally
```
export KUBE_CONFIG=${HOME}/.kube/config
cd monitor
go run main.go
```
open http://localhost:8081/jobs

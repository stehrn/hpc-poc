# Monitor container
Simple web app to view jobs and logs 

# Build
Run:
```
gcloud builds submit --tag gcr.io/hpc-poc/monitor
```
View image:
```
gcloud container images list --repository=gcr.io/hpc-poc
gcloud container images list-tags gcr.io/hpc-poc/monitor
```

# Deploy
```
kubectl create deployment monitor --image=gcr.io/hpc-poc/monitor:latest
kubectl get deployment monitor -o yaml 
```
TODO: expose and find URL

# Monitor
View logs
```
kubectl logs --selector=app=monitor --tail 100
``` 

# Run locally
```
export KUBE_CONFIG=${HOME}/.kube/config
go run main.go
```

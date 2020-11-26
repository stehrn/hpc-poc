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
# Deploy
Run:
```
cd monitor
kubectl apply -f yaml
kubectl get deployment monitor -o yaml 
```
Note the orchestrator container uses the k8 API to create Jobs, the required role, service account and bindings are also created.

Expose:
```
kubectl expose deployment monitor --type=LoadBalancer --port 8081 --name=monitor
kubectl get services monitor
NAME      TYPE           CLUSTER-IP      EXTERNAL-IP    PORT(S)          AGE
monitor   LoadBalancer   10.47.249.184   35.234.146.8   8081:31184/TCP   2m19s
```
Open browswer at: http://<external-ip>:<port>/jobs

(From above exmaple - http://35.234.146.8:8081/jobs)

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
open http://localhost:8081/jobs

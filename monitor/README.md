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
Expose:
```
kubectl expose deployment monitor --type=LoadBalancer --port 8081 --name=monitor
kubectl get services monitor
NAME      TYPE           CLUSTER-IP      EXTERNAL-IP     PORT(S)          AGE
monitor   LoadBalancer   10.47.250.242   35.189.67.125   8081:32233/TCP   36s
```
Open browswer at: http://<external-ip>:<port>/jobs

(From above exmaple - http://35.189.67.125:32233/jobs)


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

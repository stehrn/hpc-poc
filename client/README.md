# Client container
Simple web app to submit jobs

# Build
Run (from base module dir):
```
gcloud builds submit --config=cloudbuild.yaml --substitutions=_PACKAGE="client",_DOCKERFILE="DockerfileForWeb" .

```
View image:
```
gcloud container images list --repository=gcr.io/hpc-poc
gcloud container images list-tags gcr.io/hpc-poc/client
```

# Deploy
Run:
```
kubectl apply -f client/yaml
kubectl get deployment client -o yaml 
```
Get port:
```
kubectl get service client
```
Open browswer at: http://<external-ip>:<port>/client

(from above example - http://35.234.146.8:8082/client)

# Monitor
View logs:
```
kubectl logs --selector=app=client --tail 100
``` 

# Delete
Run:
```
kubectl delete pods,services -l app=client
kubectl delete deployment client
```

# Run locally
```
export PROJECT_NAME=hpc-poc
export GOOGLE_APPLICATION_CREDENTIALS=<path>key.json (see main README and 'Get GCP JSON key...')
export BUCKET_NAME=stehrn_hpc-poc
export TOPIC_NAME=test-topic
go run main.go
```
open http://localhost:8082/client

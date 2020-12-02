# Client container
Client API for submitting jobs, includes simple web app to submit jobs

Note, some of following commands require: `export PROJECT_NAME=<GCP project>`

# Build
Run (from base module dir):
```
gcloud builds submit --config=cloudbuild.yaml --substitutions=_PROJECT=${PROJECT_NAME},_PACKAGE="client",_DOCKERFILE="DockerfileForWeb" .

```
View image:
```
gcloud container images list --repository=gcr.io/${PROJECT_NAME}
gcloud container images list-tags gcr.io/${PROJECT_NAME}/client
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
Open browswer at: http://<external-ip>:<port>/client (use port number on left hand side of ':')

# Monitor
View logs:
```
kubectl logs --selector=app=client --tail 100
``` 

# Delete
Run:
```
kubectl delete deployment client && kubectl delete pods,services -l app=client 
```

# Run locally
```
export PROJECT_NAME=hpc-poc
export GOOGLE_APPLICATION_CREDENTIALS=<path>key.json (see main README and 'Get GCP JSON key...')
export BUCKET_NAME=stehrn_hpc-poc
export TOPIC_NAME=test-topic
go run *.go
```
open http://localhost:8082/client

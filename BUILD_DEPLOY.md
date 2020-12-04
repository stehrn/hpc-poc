# Build and deployment process 

* Set project name: `export PROJECT_NAME=<GCP project>`

## Everything
Run `bin/build_and_deploy.sh`

## Individually
* Set ${APP_NAME} & ${DOCKERFILE}:

| APP_NAME     | DOCKERFILE       |
|--------------|------------------|
| orchestrator | Dockerfile       |
| engine       | Dockerfile       |
| client       | DockerfileForWeb |
| monitor      | DockerfileForWeb |

## Build
Run (from base module dir):
```
gcloud builds submit --config=cloudbuild.yaml --substitutions=_PROJECT=${PROJECT_NAME},_PACKAGE=${APP_NAME},_DOCKERFILE=${DOCKERFILE} .
```
View image:
```
gcloud container images list --repository=gcr.io/${PROJECT_NAME}
gcloud container images list-tags gcr.io/${PROJECT_NAME}/${APP_NAME}
```

## Deploy
If web based app, deply usiing:
```
kubectl apply -f cmd/${APP_NAME}/yaml
kubectl get deployment ${APP_NAME} -o yaml 
```

Get port:
```
kubectl get service ${APP_NAME}
```
Open browswer at: http://<external-ip>:<port>/ (use port number on left hand side of ':')

## Monitor
View logs
```
kubectl logs --selector=app=${APP_NAME} --tail 100
``` 

## Delete
Run:
```
kubectl delete deployment ${APP_NAME} && kubectl delete pods,services -l app=${APP_NAME}
```

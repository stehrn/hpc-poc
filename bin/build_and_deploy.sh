
# build and deploy everything

echo "building orchestrator"
gcloud builds submit --config=cloudbuild.yaml --substitutions=_PROJECT=${PROJECT_NAME},_PACKAGE=orchestrator,_DOCKERFILE=Dockerfile .

echo "building engine"
gcloud builds submit --config=cloudbuild.yaml --substitutions=_PROJECT=${PROJECT_NAME},_PACKAGE=engine,_DOCKERFILE=Dockerfile .

echo "building client"
gcloud builds submit --config=cloudbuild.yaml --substitutions=_PROJECT=${PROJECT_NAME},_PACKAGE=client,_DOCKERFILE=DockerfileForWeb .

echo "building monitor"
gcloud builds submit --config=cloudbuild.yaml --substitutions=_PROJECT=${PROJECT_NAME},_PACKAGE=monitor,_DOCKERFILE=DockerfileForWeb .

echo "deploying orchestrator"
kubectl apply -f cmd/orchestrator/yaml

echo "deploying client"
kubectl apply -f cmd/client/yaml

echo "deploying monitor"
kubectl apply -f cmd/monitor/yaml
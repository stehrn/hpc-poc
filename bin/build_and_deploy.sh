#!/bin/bash 

# submit to gcp cloud build 
function build {
    package=$1 
    dockerfile=$2
    echo "building ${package} using ${dockerfile}"

    gcloud builds submit --config=cloudbuild.yaml --substitutions=_PROJECT=${PROJECT_NAME},_PACKAGE=${package},_DOCKERFILE=${dockerfile} .

    retVal=$?
    if [ $retVal -ne 0 ]; then
        echo "Error: $retVal"
    fi
}

# build and deploy everything

build orchestrator Dockerfile
build engine Dockerfile
build client DockerfileForWeb
build monitor DockerfileForWeb

echo "deploying orchestrator"
kubectl apply -f cmd/orchestrator/yaml

echo "deploying client"
kubectl apply -f cmd/client/yaml

echo "deploying monitor"
kubectl apply -f cmd/monitor/yaml
  
exit $?
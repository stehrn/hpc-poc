#!/bin/bash 

# edit as you wish

# kubernetes 
export NAMESPACE=default
export KUBE_CONFIG=${HOME}/.kube/config

# gcp
export GOOGLE_APPLICATION_CREDENTIALS=${HOME}/key.json
export PROJECT_NAME=hpc-poc
export CLOUD_STORAGE_BUCKET_NAME=${PROJECT_NAME}-bucket
export IMAGE_REGISTRY=gcr.io/${PROJECT_NAME}

#  app specific
export BUSINESS_NAME=bu2
export TASK_LOAD_FACTOR=0.6

go run *.go

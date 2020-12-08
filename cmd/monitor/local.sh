#!/bin/bash 

# edit as you wish

# kubernetes
export NAMESPACE=default
export KUBE_CONFIG=${HOME}/.kube/config

# gcp
export GOOGLE_APPLICATION_CREDENTIALS=${HOME}/key.json
export PROJECT_NAME=hpc-poc
export CLOUD_STORAGE_BUCKET_NAME=${PROJECT_NAME}-bucket

# app specific
export BUSINESS_NAMES=bu1,bu2,bu3
export PORT=8081

go run *.go

#!/bin/bash 

# edit as you wish

# gcp
export PROJECT_NAME=hpc-poc
export GOOGLE_APPLICATION_CREDENTIALS=${HOME}/key.json
export CLOUD_STORAGE_BUCKET_NAME=${PROJECT_NAME}-bucket

# app specific
export BUSINESS_NAMES=bu1,bu2,bu3
export PORT=8082

go run *.go

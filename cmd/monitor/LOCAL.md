# Monitor container
Simple web app to view jobs and logs 

# Run locally
```
export CLOUD_STORAGE_BUCKET_NAME=${PROJECT_NAME}-bucket
export NAMESPACE=default
export BUSINESS_NAMES=bu1,bu2
export KUBE_CONFIG=${HOME}/.kube/config
export PORT=8081
go run *.go
```
open http://localhost:8081/summary
# Monitor container
Simple web app to view jobs and logs 

# Run locally
```
export NAMESPACE=default
export BUSINESS_NAMES=bu1,bu2
export KUBE_CONFIG=${HOME}/.kube/config
export PORT=8081
go run main.go
```
open http://localhost:8081/summary
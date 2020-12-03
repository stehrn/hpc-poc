# Monitor container
Simple web app to view jobs and logs 

# Run locally
```
export NAMESPACE=default
export KUBE_CONFIG=${HOME}/.kube/config
cd monitor
go run main.go
```
open http://localhost:8081/summary
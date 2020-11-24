


To build image:
```
go mod init
gcloud builds submit --tag gcr.io/hpc-poc/monitor
```


subscriptions
|date/time|id|data|job|
|dddmmmttt|123|<link>|<link>|

jobs
|name|status|start time|completion time|duration|logs|


go get github.com/stehrn/hpc-poc/kubernetes@main
go run main.go 
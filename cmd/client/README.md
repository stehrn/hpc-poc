# Client container
Simple web app to submit jobs

# Run locally
```
export PROJECT_NAME=hpc-poc
export GOOGLE_APPLICATION_CREDENTIALS=<path>key.json (see main README and 'Get GCP JSON key...')
export CLOUD_STORAGE_BUCKET_NAME=${PROJECT_NAME}
export BUSINESS_NAMES=bu1,bu2
export PORT=8082
go run *.go
```
open http://localhost:8082/

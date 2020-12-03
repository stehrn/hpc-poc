# Client container
Simple web app to submit jobs

# Run locally
```
export PROJECT_NAME=hpc-poc
export GOOGLE_APPLICATION_CREDENTIALS=<path>key.json (see main README and 'Get GCP JSON key...')
export BUCKET_NAME=stehrn_hpc-poc
export TOPIC_NAME=test-topic
go run *.go
```
open http://localhost:8082/client

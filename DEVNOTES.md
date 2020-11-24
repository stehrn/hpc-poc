# Ad hoc notes

## Deploy to Cloud Run
```
gcloud run deploy --image gcr.io/hpc-poc/monitor --platform managed
```

## Go specific build
https://cloud.google.com/cloud-build/docs/building/build-go

Pushes go binary to storage bucket, not used for container

### Create storage bucket
https://cloud.google.com/storage/docs/creating-buckets
```
gsutil mb -p hpc-poc -l EUROPE-WEST2 -c STANDARD gs://nik-stehr-hpc
gsutil ls
gsutil ls -r gs://nik-stehr-hpc
```
### Submit build
```
gcloud builds submit --config cloudbuild-go.yaml go-app
```
[View build](https://console.cloud.google.com/cloud-build/builds/)

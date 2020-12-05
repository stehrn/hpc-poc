# Engine container
A simple engine for testing that downloads an object from GCP cloud storage bucket, prints it out, and exits.

Send it a `PANIC` message, and it crashes.

Image is referenced in orchestrator deployment (`ENGINE_IMAGE`)

The engine itself requires the following env variables:
* `CLOUD_STORAGE_BUCKET_NAME` 
* `CLOUD_STORAGE_OBJECT_NAME`
* `GOOGLE_APPLICATION_CREDENTIALS` to enable download payload from cloud storage

# Deploy
Engine image is not deployed as container into GKE, rather referenced in `Job` created by [orchestrator](../orchestrator/README.md)






# Storage API
API to upload, view, download and delete objects in cloud storage. Used to store data about jobs in cloud storage

Data is represented as a byte slice (`[]byte`), a `Location` defines _where_ the data object is stored, it has a bucket name and an object name:
```
type Location struct {
	Bucket string `json:"bucket"`
	Object string `json:"object"`
}
```

See [storage-integration_test.go](storage-integration_test.go) for full set of use cases, to follow is a basic user guide to get started.

# Creating locations

## Define own object name
```
location = client.LocationForObject("my/path/object_name") 
```
## Generate unique object name but define own path to object
```
location = client.Location("my/path") 
```
Note bucket is derived from `CLOUD_STORAGE_BUCKET_NAME`

## Create your own Location
```
location := Location{Bucket: "my-cuket", Object: "my/path/object_name" }
```

# Create new cloud storage client
```
client, err := NewEnvClient()
if err != nil {
    log.Fatal("Could not create client", err)
}
```

# Upload
Upload data to object called 'object_name' located at `gsutil ls gs://${CLOUD_STORAGE_BUCKET_NAME}/my/path/object_name`:
```
location = client.LocationForObject("my/path/object_name") 
data := []byte("payload")
err = client.Upload(location, data)
if err != nil {
   log.Fatal("Could not upload data", err)
}
```
Note bucket is derived from `CLOUD_STORAGE_BUCKET_NAME`

# Download
Using location defined previously:
```
download, err := client.Download(location)
if err != nil {
   log.Fatal("Could not download data", err)
}
```

# Delete
Using location defined previously:
```
err = client.Delete(location)
if err != nil {
   log.Errorf("Could not delete object at %q, error: %v", location, err)
}
```    

# Integration Test
See [storage-integration_test.go](storage-integration_test.go) 

One-off to set things up:
```
export CLOUD_STORAGE_BUCKET_NAME=<bucket name>
gsutil mb gs://${CLOUD_STORAGE_BUCKET_NAME}
export GOOGLE_APPLICATION_CREDENTIALS=${HOME}/integration_test_key.json
```
See [iam/README.md](iam/README.md) for info on creaitng the service account key 

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

# Create new cloud storage client
Create client:
```
client, err := NewClient("<unique-project-bucket-name>")
if err != nil {
    log.Fatal("Could not create client", err)
}
```

# Creating locations

## Create new Location
Create reference to Object location: `gs://my-unique-project-bucket/my/path/object_name`:
```
location := Location{Bucket: "my-unique-project-bucket", Object: "my/path/object_name" }
```

## Define own object name
Use client to create reference to Object location: `gs://${CLOUD_STORAGE_BUCKET_NAME}/my/path/object_name`:
```
location = client.LocationForObject("my/path/object_name") 
```

## Generate unique object name but define own path to object
Use client to create reference to Object location: `gs://${CLOUD_STORAGE_BUCKET_NAME}/my/path/<unique ID>`:
```
location = client.Location("my/path") 
```

# Upload
Upload data:
```
data := []byte("payload")
err = client.Upload(location, data)
if err != nil {
   log.Fatal("Could not upload data", err)
}
```

# Download
Download data:
```
download, err := client.Download(location)
if err != nil {
   log.Fatal("Could not download data", err)
}
```

# Delete
Delete an object:
```
err = client.Delete(location)
if err != nil {
   log.Errorf("Could not delete object at %q, error: %v", location, err)
}
```    

# Integration Test
See [storage_integration_test.go](storage_integration_test.go) 

## Set-up - create cloud storage bucket
One off task, test will upload objects into bucket, and delete data before exit
```
export CLOUD_STORAGE_BUCKET_NAME=<bucket name>
gsutil mb gs://${CLOUD_STORAGE_BUCKET_NAME}
```

## Run test
```
export GOOGLE_APPLICATION_CREDENTIALS=${HOME}/integration_test_key.json
go test -v
```
See [iam/README.md](iam/README.md) for info on creaitng the service account key 

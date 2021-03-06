// Integration test for cloud storage
//
// Following env variables required:
//
// export CLOUD_STORAGE_BUCKET_NAME=hpc-poc-bucket
// export GOOGLE_APPLICATION_CREDENTIALS=${HOME}/integration_test_key.json
//
package storage

import (
	"bytes"
	"os"
	"testing"

	"github.com/stehrn/hpc-poc/client"
)

var storageClient ClientInterface
var location Location
var uploaded bool

func TestMain(m *testing.M) {
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestStorage(t *testing.T) {

	var err error
	t.Log("Creating new client")
	storageClient, err = NewClient()
	if err != nil {
		t.Fatal("Could not create client:", err)
	}

	business := "integration/test-storage"
	bucketName := BucketNameFromEnv()
	bucketExists, err := storageClient.BucketExists(bucketName)
	if !bucketExists {
		t.Fatalf("Bucket %s does not exist: %v", bucketName, err)
	}

	location = NewLocation(bucketName, business)

	t.Logf("Uploading to %q", location)
	data := []byte("abc")
	err = storageClient.Upload(location, data)
	if err != nil {
		t.Fatal("Could not upload data ", err)
	}
	uploaded = true

	t.Logf("Downloading from %q", location)
	download, err := storageClient.Download(location)
	if err != nil {
		t.Fatal("Could not download data:", err)
	}

	if !bytes.Equal(data, download) {
		t.Errorf("Download looks odd, got: %s, want: %s", string(download), string(data))
	}

	t.Logf("Listing storage objects for location %q", location)
	objects, err := storageClient.ListObjects(location)
	if err != nil {
		t.Error("Could not list objects", err)
	}

	if len(objects) != 1 {
		t.Error("Expected just one object", err)
	} else {
		object := objects[0].Object
		if object != location.Object {
			t.Errorf("List objects failed, got: %s, want: %s.", location.Object, object)
		}
	}

	t.Logf("Deleting storage objects at %q", location)
	err = storageClient.Delete(location)
	if err != nil {
		t.Errorf("Could not delete object at %q, error: %v", location, err)
	}

	t.Logf("Checking storage object deleted")
	objects, err = storageClient.ListObjects(location)
	if err != nil {
		t.Error("Could not list objects", err)
	}

	if len(objects) != 0 {
		t.Errorf("Expected zero objects, got: %v", objects)
	}
}

func TestUploadMany(t *testing.T) {

	var err error
	t.Log("Creating new client")
	storageClient, err = NewClient()
	if err != nil {
		t.Fatal("Could not create client:", err)
	}

	bucketName := BucketNameFromEnv()

	t.Log("Uploading many objects")
	// var items client.DataSourceIterator
	dataSource := &TestDataSource{data: []byte("ABC€")}
	items := &TestDataSourceIterator{dataSource}
	uploaded := storageClient.UploadMany(bucketName, items)
	if uploaded != 1 {
		t.Error("Expected one upload")
	}

	// delete bucket...
	location := NewLocation(BucketNameFromEnv(), dataSource.ObjectPath().BusinessDir())
	t.Logf("Deleting location %s, is directory: %v", location, location.IsDirectory())
	storageClient.Delete(location)

	t.Log("Checking storage object deleted")
	objects, err := storageClient.ListObjects(location)
	if err != nil {
		t.Error("Could not list objects", err)
	}

	if len(objects) != 0 {
		t.Errorf("Expected zero objects, got: %v", objects)
	}

	t.Log("Test ok")
}

func teardown() {
}

type TestDataSourceIterator struct {
	dataSource client.DataSource
}

func (t *TestDataSourceIterator) Each(handler func(client.DataSource)) {
	handler(t.dataSource)
}

func (t *TestDataSourceIterator) Size() int {
	return 1
}

type TestDataSource struct {
	data   []byte
	errors []error
}

func (d *TestDataSource) ObjectPath() *client.ObjectPath {
	return client.ObjectPathForJob("integration/test-storage-many", "session-1", "job-1")
}

func (d *TestDataSource) Data() []byte {
	return d.data
}

func (d *TestDataSource) AddError(err error) {
	d.errors = append(d.errors, err)
}

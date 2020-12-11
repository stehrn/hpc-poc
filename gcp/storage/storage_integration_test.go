// Integration test for cloud storage
//
// Following env variables required:
// + CLOUD_STORAGE_BUCKET_NAME (and pre created)
// + GOOGLE_APPLICATION_CREDENTIALS
//
package storage

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

// export CLOUD_STORAGE_BUCKET_NAME=hpc-poc-bucket
// export GOOGLE_APPLICATION_CREDENTIALS=${HOME}/key.json

const business = "integration-test"

var client StorageClient
var location Location
var uploaded bool

func TestMain(m *testing.M) {
	code := m.Run()
	teardown()
	os.Exit(code)
}

func Test(t *testing.T) {

	var err error

	t.Log("Creating new client")
	client, err = NewEnvClient()
	if err != nil {
		t.Fatal("Could not create client", err)
	}

	location = client.Location(business)

	t.Logf("Uploading to %q", location)
	data := []byte("abc")
	err = client.Upload(location, data)
	if err != nil {
		t.Fatal("Could not upload data", err)
	}
	uploaded = true

	t.Logf("Downloading from %q", location)
	download, err := client.Download(location)
	if err != nil {
		t.Fatal("Could not download data", err)
	}

	if !bytes.Equal(data, download) {
		t.Errorf("Download looks odd, got: %s, want: %s", string(download), string(data))
	}

	t.Logf("Listing storage objects for business %q", business)
	objects, err := client.ListObjects(business)
	if err != nil {
		t.Error("Could not list objects", err)
	}

	if len(objects) != 1 {
		t.Error("Expected just one object", err)
	}

	object := objects[0].Object
	if object != location.Object {
		t.Errorf("List objects failed, got: %s, want: %s.", location.Object, object)
	}

	t.Logf("Deleting storage objects at %q", location)
	err = client.Delete(location)
	if err != nil {
		t.Errorf("Could not delete object at %q, error: %v", location, err)
	}

	t.Logf("Checking storage object deleted")
	objects, err = client.ListObjects(business)
	if err != nil {
		t.Error("Could not list objects", err)
	}

	if len(objects) != 0 {
		t.Errorf("Expected zero objects, got: %v", objects)
	}

	t.Log("Test ok")
}

func teardown() {
	if uploaded {
		fmt.Printf("Deleting %v", location)
		client.Delete(location)
	}
}

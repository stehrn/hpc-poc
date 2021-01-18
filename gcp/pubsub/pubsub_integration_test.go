// Integration test for pubsub
//
// Following env variables required:
//
// export CLOUD_STORAGE_BUCKET_NAME=hpc-poc-bucket
// export GOOGLE_APPLICATION_CREDENTIALS=${HOME}/integration_test_key.json
//
package pubsub

import (
	"context"
	"os"
	"sync"
	"testing"

	"cloud.google.com/go/pubsub"
)

const project = "hpc-poc"
const ID = "hpc-poc-int-test"

var tmpPub *TempPubSub

func TestMain(m *testing.M) {
	code := m.Run()
	teardown()
	os.Exit(code)
}

func Test(t *testing.T) {

	var err error

	t.Log("Creating new client")
	client, err := NewPubClient(project)
	if err != nil {
		t.Fatal("Could not create client", err)
	}

	t.Logf("Creating temp subscription %q", ID)
	tmpPub, err = client.NewTempPubSub(ID)
	if err != nil {
		t.Fatal("Could not create temp subscription client", err)
	}

	t.Logf("Publishing to topic %q", tmpPub.TopicName)
	payload := make([][]byte, 1)
	payload[0] = []byte("abc")

	err = tmpPub.PublishMany(payload)
	if err != nil {
		t.Fatal("Could not pulish to topic", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	t.Logf("Subscribing to subscription %q", tmpPub.client.SubscriptionID)
	go tmpPub.Subscribe(func(ctx context.Context, m *pubsub.Message) {
		t.Logf("Message recieved, ID: %q, data: %s", m.ID, string(m.Data))
		m.Ack()
		wg.Done()
	})

	t.Log("Waiting for subscriber...")
	wg.Wait()
}

func teardown() {
	if tmpPub != nil {
		tmpPub.Delete()
	}
}

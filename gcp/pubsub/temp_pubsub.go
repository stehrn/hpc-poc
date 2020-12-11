package pubsub

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
)

// TempPubSub a temporary pubsub client intended to be short lived
type TempPubSub struct {
	client *Client
}

// NewTempPubSub create temp pub sub
func (c Client) NewTempPubSub(ID string) (*TempPubSub, error) {
	ctx := context.Background()

	topicID := fmt.Sprintf("%s-topic", ID)
	topic, err := c.CreateTopic(ctx, topicID)
	if err != nil {
		return nil, fmt.Errorf("CreateTopic: %v", err)
	}

	subID := fmt.Sprintf("%s-subscription", ID)
	_, err = c.CreateSubscription(ctx, subID, pubsub.SubscriptionConfig{
		Topic: topic,
	})
	if err != nil {
		return nil, fmt.Errorf("CreateSubscription: %v", err)
	}

	return c.client(subID, topicID), nil
}

func (c Client) client(subID, topicID string) *TempPubSub {
	config := &ClientConfg{
		Project:        c.Project,
		SubscriptionID: subID,
		TopicName:      topicID}
	return &TempPubSub{&Client{config, c.Client}}
}

// ExistingTempPubSub load existing TempPubSub
func (c Client) ExistingTempPubSub(ID string) *TempPubSub {
	topicID := fmt.Sprintf("%s-topic", ID)
	subID := fmt.Sprintf("%s-subscription", ID)
	return c.client(subID, topicID)
}

// TopicName topic name
func (t *TempPubSub) TopicName() string {
	return t.client.TopicName
}

// SubscriptionID subscription ID
func (t *TempPubSub) SubscriptionID() string {
	return t.client.SubscriptionID
}

// PublishMany publish to temp topic
func (t *TempPubSub) PublishMany(payloads [][]byte) error {
	return t.client.PublishMany(payloads)
}

// Subscribe to temp subscription
func (t *TempPubSub) Subscribe(callback func(ctx context.Context, m *pubsub.Message)) error {
	return t.client.Subscribe(callback)
}

// Delete delete topic and subscription
func (t *TempPubSub) Delete() error {

	if t.client.TopicName != "" {
		log.Printf("deleting topic %q\n", t.client.TopicName)
		topic, err := t.client.topic()
		if err != nil {
			return fmt.Errorf("Could not load topic: %v", err)
		}
		ctx := context.Background()
		err = topic.Delete(ctx)
		if err != nil {
			return fmt.Errorf("Could not delete topic: %v", err)
		}
	}

	if t.client.SubscriptionID != "" {
		log.Printf("deleting subscription %q\n", t.client.SubscriptionID)
		sub, err := t.client.subscription()
		if err != nil {
			return fmt.Errorf("Could not load subscription: %v", err)
		}
		ctx := context.Background()
		err = sub.Delete(ctx)
		if err != nil {
			return fmt.Errorf("Could not delete subscription: %v", err)
		}
	}
	return nil
}

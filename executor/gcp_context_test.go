package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func context() *GcpContext {
	return &GcpContext{
		Project:  "project-1",
		Business: "bu-1",
	}
}

func TestTopicName(t *testing.T) {
	assert.Equal(t, "project-1-bu-1-topic", context().TopicName(), "wrong topic name")
}

func TestSubscriptionName(t *testing.T) {
	assert.Equal(t, "project-1-bu-1-subscription", context().SubscriptionName(), "wrong subscription name")
}

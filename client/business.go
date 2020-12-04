package client

import (
	"fmt"

	"github.com/stehrn/hpc-poc/internal/utils"
)

// Business business
type Business string

// BusinessFromEnv create a Business from 'BUSINESS_NAME' env variable
func BusinessFromEnv() Business {
	business := utils.Env("BUSINESS_NAME")
	return Business(business)
}

// TopicName derive name of topic
func (b Business) TopicName(project string) string {
	return fmt.Sprintf("%s-%s-topic", project, b)
}

// SubscriptionName derive name of subscription
func (b Business) SubscriptionName(project string) string {
	return fmt.Sprintf("%s-%s-subscription", project, b)
}

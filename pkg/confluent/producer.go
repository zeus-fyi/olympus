package confluent

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	bootstrapServers          = "<BOOTSTRAP_SERVERS>"
	ccloudAPIKey              = "<CCLOUD_API_KEY>"
	ccloudAPISecret           = "<CCLOUD_API_SECRET>"
	schemaRegistryAPIEndpoint = "<CCLOUD_SR_ENDPOINT>"
	schemaRegistryAPIKey      = "<CCLOUD_SR_API_KEY>"
	schemaRegistryAPISecret   = "<CCLOUD_SR_API_SECRET>"
)

func createProducer() (*kafka.Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"sasl.mechanisms":   "PLAIN",
		"security.protocol": "SASL_SSL",
		"sasl.username":     ccloudAPIKey,
		"sasl.password":     ccloudAPISecret},
	)
	return p, err
}

func createAdminClientFromProducer(p *kafka.Producer) (a *kafka.AdminClient, err error) {
	a, err = kafka.NewAdminClientFromProducer(p)
	if err != nil {
		fmt.Printf("Failed to create new admin client from producer: %s", err)
		return nil, err
	}
	return a, err
}

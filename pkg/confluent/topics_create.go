package confluent

import (
	"context"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type CreateTopicPayload struct {
	TopicName         string `json:"topic_name"`
	PartitionsCount   int    `json:"partitions_count"`
	ReplicationFactor int    `json:"replication_factor"`
	Configs           []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"configs"`
}

// CreateTopic creates a topic using the Admin Client API
func CreateTopic(ctx context.Context, p *kafka.Producer, topic CreateTopicPayload) error {
	a, err := createAdminClientFromProducer(p)
	if err != nil {
		fmt.Printf("Failed to create new admin client from producer: %s", err)
		return err
	}
	defer a.Close()
	// Create topics on cluster.
	// Set Admin options to wait up to 60s for the operation to finish on the remote cluster
	maxDur, err := time.ParseDuration("60s")
	if err != nil {
		fmt.Printf("ParseDuration(60s): %s", err)
		return err
	}
	results, err := a.CreateTopics(
		ctx,
		// Multiple topics can be created simultaneously
		// by providing more TopicSpecification structs here.
		[]kafka.TopicSpecification{{
			Topic:             topic.TopicName,
			NumPartitions:     topic.PartitionsCount,
			ReplicationFactor: topic.ReplicationFactor},
		})
	// Admin options
	kafka.SetAdminOperationTimeout(maxDur)
	if err != nil {
		fmt.Printf("Admin Client request error: %v\n", err)
		return err
	}
	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError && result.Error.Code() != kafka.ErrTopicAlreadyExists {
			fmt.Printf("Failed to create topic: %v\n", result.Error)
			return err
		}
		fmt.Printf("%v\n", result)
	}
	return err
}

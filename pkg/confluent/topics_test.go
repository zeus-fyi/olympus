package confluent

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
)

type TopicsTestSuite struct {
	base.TestSuite
}

func (s *TopicsTestSuite) TestCreateTopic() {
	producer, err := createProducer()
	s.Require().Nil(err)
	s.Require().NotEmpty(producer)

	ac, err := createAdminClientFromProducer(producer)
	s.Require().Nil(err)
	defer ac.Close()
	s.Require().NotEmpty(ac)

	topic := CreateTopicPayload{
		TopicName:         "test",
		PartitionsCount:   1,
		ReplicationFactor: 3,
		Configs:           nil,
	}
	// Contexts are used to abort or limit the amount of time
	// the Admin call blocks waiting for a result.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create topics on cluster.
	// Set Admin options to wait for the operation to finish (or at most 60s)
	maxDur, err := time.ParseDuration("60s")
	if err != nil {
		panic("ParseDuration(60s)")
	}
	results, err := ac.CreateTopics(
		ctx,
		// Multiple topics can be created simultaneously
		// by providing more TopicSpecification structs here.
		[]kafka.TopicSpecification{{
			Topic:             topic.TopicName,
			NumPartitions:     topic.PartitionsCount,
			ReplicationFactor: topic.ReplicationFactor}},
		// Admin options
		kafka.SetAdminOperationTimeout(maxDur),
	)

	s.Require().Nil(err)
	s.Assert().NotEmpty(results)

	// Print results
	for _, result := range results {
		fmt.Printf("%s\n", result)
	}

}

func TestTopicsTestSuite(t *testing.T) {
	suite.Run(t, new(TopicsTestSuite))
}

// RecordValue represents the struct of the value in a Kafka message
type RecordValue struct {
	Count int
}

// ParseArgs parses the command line arguments and
// returns the config file and topic on success, or exits on error
func ParseArgs() (*string, *string) {

	configFile := flag.String("f", "", "Path to Confluent Cloud configuration file")
	topic := flag.String("t", "", "Topic name")
	flag.Parse()
	if *configFile == "" || *topic == "" {
		flag.Usage()
		os.Exit(2) // the same exit code flag.Parse uses
	}

	return configFile, topic

}

// ReadCCloudConfig reads the file specified by configFile and
// creates a map of key-value pairs that correspond to each
// line of the file. ReadCCloudConfig returns the map on success,
// or exits on error
func ReadCCloudConfig(configFile string) map[string]string {

	m := make(map[string]string)

	file, err := os.Open(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %s", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "#") && len(line) != 0 {
			kv := strings.Split(line, "=")
			parameter := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			m[parameter] = value
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Failed to read file: %s", err)
		os.Exit(1)
	}

	return m

}

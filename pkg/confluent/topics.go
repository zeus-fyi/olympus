package confluent

type Topic struct {
	TopicName         string `json:"topic_name"`
	PartitionsCount   int    `json:"partitions_count"`
	ReplicationFactor int    `json:"replication_factor"`
	Configs           []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"configs"`
}

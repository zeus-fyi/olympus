package confluent

type TopicResponse struct {
	Kind     string `json:"kind"`
	Metadata struct {
		Self         string `json:"self"`
		ResourceName string `json:"resource_name"`
	} `json:"metadata"`
	ClusterId         string `json:"cluster_id"`
	TopicName         string `json:"topic_name"`
	IsInternal        bool   `json:"is_internal"`
	ReplicationFactor int    `json:"replication_factor"`
	PartitionsCount   int    `json:"partitions_count"`
	Partitions        struct {
		Related string `json:"related"`
	} `json:"partitions"`
	Configs struct {
		Related string `json:"related"`
	} `json:"configs"`
	PartitionReassignments struct {
		Related string `json:"related"`
	} `json:"partition_reassignments"`
}

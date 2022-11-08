package workload_state

type InternalWorkloadStatusUpdateRequest struct {
	TopologyID     int    `db:"topology_id" json:"topology_id"`
	TopologyStatus string `db:"topology_status" json:"topology_status"`
}

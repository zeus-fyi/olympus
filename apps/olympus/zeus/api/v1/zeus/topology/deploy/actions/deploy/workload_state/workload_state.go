package workload_state

type InternalWorkloadStatusUpdateRequest struct {
	TopologyID     int    `db:"topology_id" json:"topologyID"`
	TopologyStatus string `db:"topology_status" json:"topologyStatus"`
}

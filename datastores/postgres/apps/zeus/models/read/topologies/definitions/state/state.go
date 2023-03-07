package read_topology_deployment_status

import (
	"time"

	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

type ReadDeploymentStatus struct {
	TopologyID     int       `db:"topology_id" json:"topologyID,omitempty"`
	ClusterName    string    `db:"cluster_name" json:"clusterName,omitempty"`
	TopologyName   string    `db:"topology_name" json:"topologyName"` // TODO this is being used for component base name overriding the default
	TopologyStatus string    `db:"topology_status" json:"topologyStatus"`
	UpdatedAt      time.Time `db:"updated_at" json:"updatedAt"`
	CloudCtxNs     zeus_common_types.CloudCtxNs
}

type ReadDeploymentStatusesGroup struct {
	Slice []ReadDeploymentStatus
}

func NewReadDeploymentStatusesGroup() ReadDeploymentStatusesGroup {
	s := ReadDeploymentStatusesGroup{}
	s.Slice = []ReadDeploymentStatus{}
	return s
}

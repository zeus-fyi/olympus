package read_topology_deployment_status

import (
	"time"

	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
)

type ReadDeploymentStatus struct {
	TopologyName   string    `db:"topology_name" json:"topologyName"`
	TopologyStatus string    `db:"topology_status" json:"topologyStatus"`
	UpdatedAt      time.Time `db:"updated_at" json:"updatedAt"`
	CloudCtxNs     autogen_bases.TopologiesKns
}

type ReadDeploymentStatusesGroup struct {
	Slice []ReadDeploymentStatus
}

func NewReadDeploymentStatusesGroup() ReadDeploymentStatusesGroup {
	s := ReadDeploymentStatusesGroup{}
	s.Slice = []ReadDeploymentStatus{}
	return s
}

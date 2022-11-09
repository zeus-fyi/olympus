package read_topology_deployment_status

import (
	"time"

	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
)

type ReadDeploymentStatus struct {
	TopologyName   string    `db:"topology_name" json:"topology_name"`
	TopologyStatus string    `db:"topology_status" json:"topology_status"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
	Kns            autogen_bases.TopologiesKns
}

type ReadDeploymentStatusesGroup struct {
	Slice []ReadDeploymentStatus
}

func NewReadDeploymentStatusesGroup() ReadDeploymentStatusesGroup {
	s := ReadDeploymentStatusesGroup{}
	s.Slice = []ReadDeploymentStatus{}
	return s
}

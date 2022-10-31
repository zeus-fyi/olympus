package base

import (
	clusters "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/classes/cluster"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus/core"
)

type TopologyActionRequest struct {
	Action string
	core.K8sRequest

	clusters.Cluster
}

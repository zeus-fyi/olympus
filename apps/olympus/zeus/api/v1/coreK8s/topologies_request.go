package coreK8s

import (
	clusters "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/classes/cluster"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus/zeus_pkg"
)

type TopologyActionRequest struct {
	zeus_pkg.K8sRequest
	Action string

	clusters.Cluster
}

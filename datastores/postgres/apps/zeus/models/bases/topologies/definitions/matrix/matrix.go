package matrix

import (
	clusters "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/cluster"
)

type Matrix struct {
	TopologyClusters map[int]clusters.ClustersGroup
}

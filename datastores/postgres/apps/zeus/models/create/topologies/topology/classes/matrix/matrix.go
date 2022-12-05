package create_matrix

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/classes/clusters"
)

type Matrix struct {
	TopologyClusters map[int]clusters.ClustersGroup
}

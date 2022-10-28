package topologies

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/matrix"
)

type Topology struct {
	kns.Kns
	// use map->int to represent hierarchy
	// all same level uncoupled items are on the same depth

	*matrix.Matrix
}

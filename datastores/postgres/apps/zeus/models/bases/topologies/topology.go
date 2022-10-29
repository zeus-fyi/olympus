package topologies

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/classes/matrix"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
)

type Topology struct {
	autogen_bases.Topologies // -- specific sw names. eg. prysm_validator_client
	kns.Kns
	// use map->int to represent hierarchy
	// all same level uncoupled items are on the same depth

	*matrix.Matrix
}

func NewTopology() Topology {
	t := Topology{
		Kns:    kns.Kns{},
		Matrix: nil,
	}

	return t
}

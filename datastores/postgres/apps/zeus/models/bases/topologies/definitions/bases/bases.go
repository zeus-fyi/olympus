package bases

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/class_types"
)

type Base struct {
	autogen_bases.TopologyBaseComponents
}

func NewBaseClassTopologyType(orgID, systemID int, name string) Base {
	c := Base{
		TopologyBaseComponents: autogen_bases.TopologyBaseComponents{
			OrgID:                     orgID,
			TopologySystemComponentID: systemID,
			TopologyClassTypeID:       class_types.BaseClassTypeID,
			TopologyBaseName:          name,
		},
	}
	return c
}

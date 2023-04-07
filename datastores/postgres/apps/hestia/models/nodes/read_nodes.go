package hestia_nodes

import (
	"context"

	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	zeus_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
)

type NodeFilter struct {
	CloudProvider string                 `json:"cloudProvider"`
	Region        string                 `json:"region"`
	ResourceSums  zeus_core.ResourceSums `json:"resourceSums"`
}

func SelectNodes(ctx context.Context, nf NodeFilter) (hestia_autogen_bases.NodesSlice, error) {

	//q := `SELECT * FROM nodes WHERE `
	return hestia_autogen_bases.NodesSlice{}, nil
}

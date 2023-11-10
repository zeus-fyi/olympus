package zeus

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
)

func ReadUserTopologyConfig(ctx context.Context, topID int, ou org_users.OrgUser) (read_topology.InfraBaseTopology, error) {
	tr := read_topology.NewInfraTopologyReaderWithOrgUser(ou)
	tr.TopologyID = topID
	err := tr.SelectTopology(ctx)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("DeployTopology, SelectTopology error")
		return tr, err
	}
	return tr, err
}

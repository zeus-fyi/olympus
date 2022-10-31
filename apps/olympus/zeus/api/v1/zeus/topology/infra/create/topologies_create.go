package create_infra

import (
	"context"

	"github.com/labstack/echo/v4"
	create_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/topology"
	base_infra "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/base"
)

type TopologyActionCreateRequest struct {
	base_infra.TopologyInfraActionRequest
	TopologyCreateRequest
}
type TopologyCreateRequest struct {
	Name string
}

func (t *TopologyActionCreateRequest) CreateTopology(c echo.Context) error {
	ctx := context.Background()
	topCreate := create_topology.NewCreateOrgUsersInfraTopology(t.OrgUser)
	topCreate.Name = t.Name
	err := topCreate.InsertOrgUsersTopology(ctx)
	return err
}

package clean_deploy_request

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
)

func CleanNamespaceAtTopology(c echo.Context, tar zeus_req_types.TopologyDeployRequest) error {
	log.Debug().Msg("CleanNamespaceAtTopology")
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	tr := read_topology.NewInfraTopologyReaderWithOrgUser(ou)
	tr.TopologyID = tar.TopologyID
	err := tr.SelectTopology(ctx)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("CleanNamespaceAtTopology, SelectTopology error")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return zeus.ExecuteCleanDeployWorkflow(c, ctx, ou, tar, tr.GetTopologyBaseInfraWorkload())
}

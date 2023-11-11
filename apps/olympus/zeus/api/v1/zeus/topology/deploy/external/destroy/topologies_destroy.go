package destroy_deploy_request

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
)

func DestroyDeployedTopology(c echo.Context, tar zeus_req_types.TopologyDeployRequest) error {
	log.Debug().Msg("DestroyDeployedTopology")
	ctx := context.Background()
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Warn().Interface("ou", ou).Msg("DestroyDeployedTopology, orgUser not found")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	tr := read_topology.NewInfraTopologyReaderWithOrgUser(ou)
	tr.TopologyID = tar.TopologyID
	err := tr.SelectTopology(ctx)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("DestroyDeployedTopology, SelectTopology error")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return zeus.ExecuteDestroyDeployWorkflow(c, ctx, ou, tar, tr.GetTopologyBaseInfraWorkload())
}

type TopologyUIDestroyDeployRequest struct {
	zeus_common_types.CloudCtxNs `json:"cloudCtxNs"`
}

func DestroyNamespaceCluster(c echo.Context, cloudCtxNs zeus_common_types.CloudCtxNs) error {
	log.Debug().Msg("DestroyNamespaceCluster")
	ctx := context.Background()
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Warn().Interface("ou", ou).Msg("DestroyDeployedTopology, orgUser not found")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if cloudCtxNs.CheckIfEmpty() {
		log.Warn().Interface("orgUser", ou).Interface("cloudCtxNs", cloudCtxNs).Msg("DestroyNamespaceCluster, CloudCtxNs is empty")
		return c.JSON(http.StatusBadRequest, nil)
	}
	tr := read_topology.NewInfraTopologyReaderWithOrgUser(ou)
	tar := zeus_req_types.TopologyDeployRequest{
		CloudCtxNs: cloudCtxNs,
	}
	return zeus.ExecuteDestroyNamespaceWorkflow(c, ctx, ou, tar, tr.GetTopologyBaseInfraWorkload())
}

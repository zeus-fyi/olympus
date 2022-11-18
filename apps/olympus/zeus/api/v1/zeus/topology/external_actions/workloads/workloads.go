package external_api_workloads

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	zeus_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

type TopologyCloudCtxNsQueryRequest struct {
	zeus_core.CloudCtxNs
}

func (t *TopologyCloudCtxNsQueryRequest) ReadDeployedWorkloads(c echo.Context) error {
	log.Debug().Msg("TopologyCloudCtxNsQueryRequest")
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)

	// TODO move to middleware
	authed, err := read_topology.IsOrgCloudCtxNsAuthorized(ctx, ou.OrgID, t.CloudCtxNs)
	if authed != true {
		return c.JSON(http.StatusInternalServerError, err)
	}

	workload, err := zeus.K8Util.GetWorkloadAtNamespace(ctx, t.CloudCtxNs)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("TopologyCloudCtxNsQueryRequest: GetWorkloadAtNamespace")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, workload)
}

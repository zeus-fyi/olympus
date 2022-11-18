package destroy_deploy_request

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

type TopologyDestroyDeployRequest struct {
	kns.TopologyKubeCtxNs
}

func (t *TopologyDestroyDeployRequest) DestroyDeployedTopology(c echo.Context) error {
	log.Debug().Msg("DestroyDeployedTopology")
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	tr := read_topology.NewInfraTopologyReaderWithOrgUser(ou)
	tr.TopologyID = t.TopologyID
	err := tr.SelectTopology(ctx)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("DestroyDeployedTopology, SelectTopology error")
		return c.JSON(http.StatusInternalServerError, err)
	}
	// validate context kns
	authed, err := tr.IsOrgCloudCtxNsAuthorized(ctx, t.CloudCtxNs)
	if authed != true {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return zeus.ExecuteDestroyDeployWorkflow(c, ctx, ou, t.TopologyKubeCtxNs, tr.GetNativeK8s())
}

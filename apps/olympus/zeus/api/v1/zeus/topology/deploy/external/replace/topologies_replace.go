package replace_topology

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

type TopologyReplaceRequest struct {
	kns.TopologyKubeCtxNs
}

func (t *TopologyReplaceRequest) ReplaceTopology(c echo.Context) error {
	log.Debug().Msg("TopologyReplaceTopology")
	nk, err := zeus.DecompressUserInfraWorkload(c)
	if err != nil {
		log.Err(err).Interface("kubernetesWorkload", nk).Msg("TopologyReplaceTopology: DecompressUserInfraWorkload")
		return c.JSON(http.StatusBadRequest, nil)
	}
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	tr, err := zeus.ReadUserTopologyConfig(ctx, t.TopologyID, ou)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("ReplaceTopology, SelectTopology error")
		return c.JSON(http.StatusInternalServerError, err)
	}
	diffReplacement := zeus.DiffChartUpdate(nk, tr.GetTopologyBaseInfraWorkload())
	return zeus.ExecuteDeployWorkflow(c, ctx, ou, t.TopologyKubeCtxNs, diffReplacement, false)
}

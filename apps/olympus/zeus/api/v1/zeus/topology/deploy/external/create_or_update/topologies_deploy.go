package create_or_update_deploy

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

type TopologyDeployRequest struct {
	kns.TopologyKubeCtxNs
	ClusterName                     string `json:"clusterClassName,omitempty"`
	RequestChoreographySecretDeploy bool   `json:"requestChoreographySecretDeploy,omitempty"`
}

func (t *TopologyDeployRequest) DeployTopology(c echo.Context) error {
	log.Debug().Msg("TopologyDeployRequest")
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	tr, err := zeus.ReadUserTopologyConfig(ctx, t.TopologyID, ou)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("DeployTopology, ReadUserTopologyConfig error")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return zeus.ExecuteDeployWorkflow(c, ctx, ou, t.TopologyKubeCtxNs, tr.GetTopologyBaseInfraWorkload(), t.RequestChoreographySecretDeploy, t.ClusterName)
}

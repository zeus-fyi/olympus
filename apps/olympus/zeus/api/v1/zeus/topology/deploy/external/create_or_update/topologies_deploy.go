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
	TopologyKubeCtxNs               kns.TopologyKubeCtxNs `json:"topologyKubeCtxNs"`
	ClusterName                     string                `json:"clusterClassName,omitempty"`
	SecretRef                       string                `json:"secretRef,omitempty"`
	RequestChoreographySecretDeploy bool                  `json:"requestChoreographySecretDeploy,omitempty"`
}

func (t *TopologyDeployRequest) DeployTopology(c echo.Context) error {
	log.Debug().Msg("TopologyDeployRequest")
	ctx := context.Background()
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Warn().Interface("ou", ou).Interface("cloudCtxNs", t.TopologyKubeCtxNs).Msg("ClusterTopologyDeploymentHandler: orgUser")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	tr, err := zeus.ReadUserTopologyConfig(ctx, t.TopologyKubeCtxNs.TopologyID, ou)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("DeployTopology, ReadUserTopologyConfig error")
		return c.JSON(http.StatusInternalServerError, err)
	}
	nk := tr.GetTopologyBaseInfraWorkload()
	if nk.Job != nil {
		return zeus.ExecuteDeployJobWorkflow(c, ctx, ou, t.TopologyKubeCtxNs, tr.GetTopologyBaseInfraWorkload(), t.RequestChoreographySecretDeploy, t.ClusterName, t.SecretRef)
	}
	if nk.CronJob != nil {
		return zeus.ExecuteDeployCronJobWorkflow(c, ctx, ou, t.TopologyKubeCtxNs, tr.GetTopologyBaseInfraWorkload(), t.RequestChoreographySecretDeploy, t.ClusterName, t.SecretRef)
	}
	return zeus.ExecuteDeployWorkflow(c, ctx, ou, t.TopologyKubeCtxNs, tr.GetTopologyBaseInfraWorkload(), t.RequestChoreographySecretDeploy, t.ClusterName, t.SecretRef)
}

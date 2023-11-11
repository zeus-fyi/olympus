package create_or_update_deploy

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
)

func DeployTopology(c echo.Context, tar zeus_req_types.TopologyDeployRequest) error {
	log.Debug().Msg("TopologyDeployRequest")
	ctx := context.Background()
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Warn().Interface("ou", ou).Interface("cloudCtxNs", tar.CloudCtxNs).Msg("ClusterTopologyDeploymentHandler: orgUser")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	tr, err := zeus.ReadUserTopologyConfig(ctx, tar.TopologyID, ou)
	if err != nil {
		log.Err(err).Interface("req", tar).Interface("orgUser", ou).Msg("DeployTopology, ReadUserTopologyConfig error")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	nk := tr.GetTopologyBaseInfraWorkload()
	if nk.Job != nil {
		return zeus.ExecuteDeployJobWorkflow(c, ctx, ou, tar, tr.GetTopologyBaseInfraWorkload(), tar.RequestChoreographySecretDeploy, tar.ClusterClassName, tar.SecretRef)
	}
	if nk.CronJob != nil {
		return zeus.ExecuteDeployCronJobWorkflow(c, ctx, ou, tar, tr.GetTopologyBaseInfraWorkload(), tar.RequestChoreographySecretDeploy, tar.ClusterClassName, tar.SecretRef)
	}
	return zeus.ExecuteDeployWorkflow(c, ctx, ou, tar, tr.GetTopologyBaseInfraWorkload(), tar.RequestChoreographySecretDeploy, tar.ClusterClassName, tar.SecretRef)
}

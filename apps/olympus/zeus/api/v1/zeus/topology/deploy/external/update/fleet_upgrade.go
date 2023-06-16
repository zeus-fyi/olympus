package deploy_updates

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

type FleetUpgradeRequest struct {
	ClusterClassName string `json:"clusterClassName"`
	AppTaint         bool   `json:"appTaint,omitempty"`
}

func (t *FleetUpgradeRequest) UpgradeFleet(c echo.Context) error {
	log.Debug().Msg("UpgradeFleet")
	ou := c.Get("orgUser").(org_users.OrgUser)
	ctx := context.Background()
	params := base_deploy_params.FleetUpgradeWorkflowRequest{
		ClusterName: t.ClusterClassName,
		OrgUser:     ou,
		AppTaint:    t.AppTaint,
	}
	return zeus.ExecuteDeployFleetUpgradeWorkflow(c, ctx, params)
}

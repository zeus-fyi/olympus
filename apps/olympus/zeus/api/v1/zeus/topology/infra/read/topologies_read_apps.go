package read_infra

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
)

type TopologyReadPrivateAppsRequest struct {
}

func (t *TopologyReadPrivateAppsRequest) ListPrivateAppsRequest(c echo.Context) error {
	ou := c.Get("orgUser").(org_users.OrgUser)
	ctx := context.Background()
	apps, err := read_topology.SelectOrgApps(ctx, ou.OrgID)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("ListPrivateAppsRequest: SelectOrgApps")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, apps)
}

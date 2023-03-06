package read_infra

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_topologies "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies"
)

func (t *TopologyReadRequest) ReadTopologiesOrgCloudCtxNs(c echo.Context) error {
	ou := c.Get("orgUser").(org_users.OrgUser)
	ctx := context.Background()
	resp, err := read_topologies.SelectTopologiesMetadata(ctx, ou.OrgID)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("ReadTopologyChart: SelectTopology")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, resp)
}

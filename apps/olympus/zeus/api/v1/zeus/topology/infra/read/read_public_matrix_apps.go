package read_infra

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
)

type PublicAppsMatrixRequest struct {
}

func PublicAppsMatrixRequestHandler(c echo.Context) error {
	request := new(PublicAppsMatrixRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.GetPublicAppFamily(c)
}

func (a *PublicAppsMatrixRequest) GetPublicAppFamily(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Err(nil).Msg("GetAppFamily: orgUser not found")
		return c.JSON(http.StatusUnauthorized, nil)
	}
	ctx := context.Background()
	appList, err := read_topology.SelectPublicMatrixApps(ctx, AppsOrgID)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("ListPrivateAppsRequest: SelectOrgApps")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, appList)
}

package zeus_v1_clusters_api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/authorized_clusters"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

type ReadAuthorizedClustersRequest struct {
}

func ReadAuthorizedClustersRequestHandler(c echo.Context) error {
	request := new(ReadAuthorizedClustersRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.Read(c)
}

func (t *ReadAuthorizedClustersRequest) Read(c echo.Context) error {
	ou := c.Get("orgUser").(org_users.OrgUser)
	if ou.OrgID == 0 {
		return c.JSON(http.StatusUnauthorized, "Unauthorized")
	}
	extCfgs, err := authorized_clusters.SelectAuthedClusterConfigsByOrgID(c.Request().Context(), ou)
	if err != nil {
		log.Err(err).Msg("ReadExtKubeConfig: SelectAuthedClusterConfigsByOrgID")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, extCfgs)
}

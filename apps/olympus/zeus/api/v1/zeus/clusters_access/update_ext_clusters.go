package zeus_v1_clusters_api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/ext_clusters"
)

type UpdateExtClustersRequest struct {
	ExtClusterConfigs []ext_clusters.ExtClusterConfig `json:"extClusterConfigs"`
}

func UpdateExtClustersRequestHandler(c echo.Context) error {
	request := new(UpdateExtClustersRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.UpdateExtClusterConfigs(c)
}

func (t *UpdateExtClustersRequest) UpdateExtClusterConfigs(c echo.Context) error {
	ou := c.Get("orgUser").(org_users.OrgUser)
	if ou.OrgID == 0 {
		return c.JSON(http.StatusUnauthorized, "Unauthorized")
	}
	err := ext_clusters.InsertOrUpdateExtClusterConfigs(c.Request().Context(), ou, t.ExtClusterConfigs)
	if err != nil {
		log.Err(err).Msg("UpdateExtClusterConfigs: InsertOrUpdateExtClusterConfigs")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, nil)
}

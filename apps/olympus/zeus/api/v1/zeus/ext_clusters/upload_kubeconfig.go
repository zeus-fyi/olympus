package zeus_v1_clusters_api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

type CreateOrUpdateKubeConfigsRequest struct {
}

func CreateOrUpdateKubeConfigsHandler(c echo.Context) error {
	request := new(CreateOrUpdateKubeConfigsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateOrUpdateKubeConfig(c)
}

func (t *CreateOrUpdateKubeConfigsRequest) CreateOrUpdateKubeConfig(c echo.Context) error {
	ou := c.Get("orgUser").(org_users.OrgUser)
	if ou.OrgID == 0 {
		return c.JSON(http.StatusUnauthorized, "Unauthorized")
	}
	fileResp, err := DecompressUserKubeConfigsWorkload(c)
	if err != nil {
		log.Err(err).Msg("CreateOrUpdateKubeConfig: DecompressUserKubeConfigsWorkload")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	fmt.Println(fileResp)
	return c.JSON(http.StatusOK, "ok")
}

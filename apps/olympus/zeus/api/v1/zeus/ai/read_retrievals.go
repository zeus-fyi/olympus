package zeus_v1_ai

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

type GetRetrievalsRequest struct {
}

func GetRetrievalsRequestHandler(c echo.Context) error {
	request := new(GetRetrievalsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.GetRetrievals(c)
}
func (t *GetRetrievalsRequest) GetRetrievals(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if ou.OrgID <= 0 || ou.UserID <= 0 {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	ret, err := artemis_orchestrations.SelectRetrievals(c.Request().Context(), ou, 0)
	if err != nil {
		log.Err(err).Msg("failed to get retrievals")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, ret)
}

func GetRetrievalRequestHandler(c echo.Context) error {
	request := new(GetRetrievalsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	// Extracting the :id parameter from the route
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Err(err).Msg("invalid ID parameter")
		return c.JSON(http.StatusBadRequest, "invalid ID parameter")
	}
	return request.GetRetrieval(c, id)
}

func (t *GetRetrievalsRequest) GetRetrieval(c echo.Context, id int) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if ou.OrgID <= 0 || ou.UserID <= 0 {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	ret, err := artemis_orchestrations.SelectRetrievals(c.Request().Context(), ou, id)
	if err != nil {
		log.Err(err).Msg("failed to get retrievals")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, ret)
}

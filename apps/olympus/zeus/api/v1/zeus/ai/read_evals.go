package zeus_v1_ai

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

type GetEvalsRequest struct {
}

func GetEvalsRequestHandler(c echo.Context) error {
	request := new(GetEvalsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.GetEvalFns(c)
}

func (t *GetEvalsRequest) GetEvalFns(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if ou.OrgID <= 0 || ou.UserID <= 0 {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	evalFns, err := artemis_orchestrations.SelectEvalFnsByOrgIDAndID(c.Request().Context(), ou, 0)
	if err != nil {
		log.Err(err).Msg("failed to get evals")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, evalFns)
}

func GetEvalRequestHandler(c echo.Context) error {
	// Extracting the ID from the route parameter
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Err(err).Msg("invalid ID parameter")
		return c.JSON(http.StatusBadRequest, "invalid ID parameter")
	}

	request := new(GetEvalsRequest)
	if err = c.Bind(request); err != nil {
		return err
	}
	return request.GetEvalFn(c, id) // Pass the ID to the method
}

func (t *GetEvalsRequest) GetEvalFn(c echo.Context, id int) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if ou.OrgID <= 0 || ou.UserID <= 0 {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	// Use the provided ID in the SelectEvalFnsByOrgIDAndID call
	evalFns, err := artemis_orchestrations.SelectEvalFnsByOrgIDAndID(c.Request().Context(), ou, id)
	if err != nil {
		log.Err(err).Msg("failed to get evals")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, evalFns)
}

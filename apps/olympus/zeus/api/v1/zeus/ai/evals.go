package zeus_v1_ai

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

type CreateOrUpdateEvalsRequest struct {
	artemis_orchestrations.AITaskLibrary
}

func CreateOrUpdateEvalsRequestHandler(c echo.Context) error {
	request := new(CreateOrUpdateEvalsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateOrUpdateEval(c)
}

func (t *CreateOrUpdateEvalsRequest) CreateOrUpdateEval(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if ou.OrgID <= 0 || ou.UserID <= 0 {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if t.TaskType == "" || t.TaskName == "" || t.TaskGroup == "" {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	t.OrgID = ou.OrgID
	t.UserID = ou.UserID
	err := artemis_orchestrations.InsertTask(c.Request().Context(), &t.AITaskLibrary)
	if err != nil {
		log.Err(err).Msg("failed to insert task")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, nil)
}

package zeus_v1_ai

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

type CreateOrUpdateTaskRequest struct {
	artemis_orchestrations.AITaskLibrary
}

func CreateOrUpdateTaskRequestHandler(c echo.Context) error {
	request := new(CreateOrUpdateTaskRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateOrUpdateTask(c)
}

func (t *CreateOrUpdateTaskRequest) CreateOrUpdateTask(c echo.Context) error {
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
	if t.TaskStrID != "" {
		tid, err := strconv.Atoi(t.TaskStrID)
		if err != nil {
			log.Err(err).Msg("failed to parse int")
			return c.JSON(http.StatusBadRequest, nil)
		}
		t.TaskID = tid
	}
	if t.MarginBuffer == 0 {
		t.MarginBuffer = 0.5
	}
	if t.MarginBuffer < 0.2 {
		t.MarginBuffer = 0.2
	}
	if t.MarginBuffer > 0.8 {
		t.MarginBuffer = 0.8
	}

	err := artemis_orchestrations.InsertTask(c.Request().Context(), &t.AITaskLibrary)
	if err != nil {
		log.Err(err).Msg("failed to insert task")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, t)
}

type TaskRequest struct {
	artemis_orchestrations.AITaskLibrary
}

package zeus_v1_ai

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

type WorkflowsActionsRequest struct {
	Action      string   `json:"action"`
	WorkflowIDs []string `json:"workflowIDs"`
}

func WorkflowsActionsRequestHandler(c echo.Context) error {
	request := new(WorkflowsActionsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.Process(c)
}

func (w *WorkflowsActionsRequest) Process(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	log.Info().Interface("ou", ou)
	return nil
}

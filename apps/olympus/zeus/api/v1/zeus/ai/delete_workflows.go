package zeus_v1_ai

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

type WorkflowsDeletionRequest struct {
	Workflows []artemis_orchestrations.WorkflowTemplate `json:"workflows"`
}

func WorkflowsDeletionRequestHandler(c echo.Context) error {
	request := new(WorkflowsDeletionRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.DeleteWorkflows(c)
}

func (w *WorkflowsDeletionRequest) DeleteWorkflows(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Info().Interface("ou", ou)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	err := artemis_orchestrations.DeleteWorkflowTemplates(c.Request().Context(), ou, w.Workflows)
	if err != nil {
		log.Err(err).Msg("failed to delete workflows")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, nil)
}

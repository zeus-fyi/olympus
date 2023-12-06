package zeus_v1_ai

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	ai_platform_service_orchestrations "github.com/zeus-fyi/olympus/pkg/zeus/ai/orchestrations"
)

type WorkflowsActionsRequest struct {
	Action        string                                    `json:"action"`
	UnixStartTime int                                       `json:"unixStartTime,omitempty"`
	UnixEndTime   int                                       `json:"unixEndTime,omitempty"`
	Workflows     []artemis_orchestrations.WorkflowTemplate `json:"workflows"`
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
		log.Info().Interface("ou", ou)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	switch w.Action {
	case "start":
		resp, err := artemis_orchestrations.GetAiOrchestrationParams(c.Request().Context(), ou, w.UnixStartTime, w.UnixEndTime, w.Workflows)
		if err != nil {
			log.Err(err).Interface("ou", ou).Interface("[]WorkflowTemplate", w.Workflows).Msg("WorkflowsActionsRequestHandler: GetAiOrchestrationParams failed")
			return c.JSON(http.StatusInternalServerError, nil)
		}
		for _, v := range resp {
			err = ai_platform_service_orchestrations.ZeusAiPlatformWorker.ExecuteRunAiWorkflowProcess(c.Request().Context(), ou, v)
			if err != nil {
				log.Err(err).Interface("ou", ou).Interface("WorkflowExecParams", resp).Msg("WorkflowsActionsRequestHandler: ExecuteRunAiWorkflowProcess failed")
				return c.JSON(http.StatusInternalServerError, nil)
			}
		}
	case "stop":
		// do y
	}
	return c.JSON(http.StatusOK, nil)
}

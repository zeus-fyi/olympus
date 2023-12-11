package zeus_v1_ai

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_stripe "github.com/zeus-fyi/olympus/pkg/hestia/stripe"
	ai_platform_service_orchestrations "github.com/zeus-fyi/olympus/pkg/zeus/ai/orchestrations"
)

type RunsActionsRequest struct {
	Action string                                 `json:"action"`
	Runs   []artemis_autogen_bases.Orchestrations `json:"runs,omitempty"`
}

func RunsActionsRequestHandler(c echo.Context) error {
	request := new(RunsActionsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.Process(c)
}

func (w *RunsActionsRequest) Process(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Info().Interface("ou", ou)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	isBillingSetup, err := hestia_stripe.DoesUserHaveBillingMethod(c.Request().Context(), ou.UserID)
	if err != nil {
		log.Error().Err(err).Msg("failed to check if user has billing method")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if !isBillingSetup {
		return c.JSON(http.StatusPreconditionFailed, nil)
	}
	switch w.Action {
	case "start":
	case "stop":
		var runIDs []string
		for _, run := range w.Runs {
			runIDs = append(runIDs, run.OrchestrationName)
		}
		err = ai_platform_service_orchestrations.ZeusAiPlatformWorker.ExecuteCancelWorkflowRuns(c.Request().Context(), ou, runIDs)
		if err != nil {
			log.Error().Err(err).Msg("failed to cancel workflow runs")
			return c.JSON(http.StatusInternalServerError, nil)
		}
	}
	return c.JSON(http.StatusOK, nil)
}

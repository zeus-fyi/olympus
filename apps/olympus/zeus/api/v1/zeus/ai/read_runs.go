package zeus_v1_ai

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_stripe "github.com/zeus-fyi/olympus/pkg/hestia/stripe"
)

func GetRunReportsRequestHandler(c echo.Context) error {
	request := new(RunsActionsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.GetRuns(c)
}

func (w *RunsActionsRequest) GetRuns(c echo.Context) error {
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
	ojsRuns, err := artemis_orchestrations.SelectAiSystemOrchestrations(c.Request().Context(), ou)
	if err != nil {
		log.Err(err).Msg("failed to get runs")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, ojsRuns)
}

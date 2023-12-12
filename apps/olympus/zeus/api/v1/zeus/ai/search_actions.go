package zeus_v1_ai

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_stripe "github.com/zeus-fyi/olympus/pkg/hestia/stripe"
	ai_platform_service_orchestrations "github.com/zeus-fyi/olympus/pkg/zeus/ai/orchestrations"
)

func SearchIndexerActionsRequestHandler(c echo.Context) error {
	request := new(ai_platform_service_orchestrations.SearchIndexerActionsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return ProcessIndexerRequest(c, request)
}

func ProcessIndexerRequest(c echo.Context, w *ai_platform_service_orchestrations.SearchIndexerActionsRequest) error {
	if w == nil || len(w.SearchIndexers) == 0 {
		return c.JSON(http.StatusBadRequest, nil)
	}
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
	case "start", "stop":
		err = ai_platform_service_orchestrations.ZeusAiPlatformWorker.ExecuteAiSearchIndexerActionsWorkflow(c.Request().Context(), ou, *w)
		if err != nil {
			log.Err(err).Msg("failed to execute ai search indexer actions workflow")
			return c.JSON(http.StatusInternalServerError, nil)
		}
	default:
		log.Info().Interface("action", w.Action).Msg("unknown action")
		return c.JSON(http.StatusBadRequest, nil)
	}
	return c.JSON(http.StatusOK, nil)
}

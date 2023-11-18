package zeus_webhooks

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	ai_platform_service_orchestrations "github.com/zeus-fyi/olympus/pkg/zeus/ai/ai/orchestrations"
)

func (a *AIServiceRequest) SupportAcknowledgeTelegramAiTask(c echo.Context) error {
	log.Info().Msg("Zeus: SupportAcknowledgeTelegramAiTask")
	group := c.Param("group")
	if len(group) == 0 {
		return c.JSON(http.StatusBadRequest, nil)
	}
	msgs, err := ai_platform_service_orchestrations.GetPandoraMessages(c.Request().Context(), group)
	if err != nil {
		log.Err(err).Msg("Zeus: CreateAIServiceTaskRequestHandler")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	internalOrgID := 7138983863666903883
	ou := org_users.NewOrgUserWithID(internalOrgID, 7138958574876245567)
	err = ai_platform_service_orchestrations.ZeusAiPlatformWorker.ExecuteAiTelegramWorkflow(c.Request().Context(), ou, msgs)
	if err != nil {
		log.Err(err).Msg("Zeus: CreateAIServiceTaskRequestHandler")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, nil)
}

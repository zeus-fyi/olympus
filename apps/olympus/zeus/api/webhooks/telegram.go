package zeus_webhooks

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	ai_platform_service_orchestrations "github.com/zeus-fyi/olympus/pkg/zeus/ai/ai/orchestrations"
)

func AiTelegramSupportAcknowledgeTelegramAiTaskHandler(c echo.Context) error {
	log.Info().Msg("Zeus: SupportAcknowledgeTelegramAiTask")
	request := new(AIServiceRequest)
	if err := c.Bind(request); err != nil {
		log.Err(err).Msg("SupportAcknowledgeTelegramAiTask")
		return err
	}
	return request.SupportAcknowledgeTelegramAiTask(c)
}

func (a *AIServiceRequest) SupportAcknowledgeTelegramAiTask(c echo.Context) error {
	log.Info().Msg("Zeus: SupportAcknowledgeTelegramAiTask")
	group := c.Param("group")
	if len(group) == 0 {
		return c.JSON(http.StatusBadRequest, nil)
	}
	internalOrgID := 7138983863666903883
	ou := org_users.NewOrgUserWithID(internalOrgID, 7138958574876245567)
	msgs, err := ai_platform_service_orchestrations.GetPandoraMessages(c.Request().Context(), group)
	if err != nil {
		log.Err(err).Msg("Zeus: CreateAIServiceTaskRequestHandler")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	err = ai_platform_service_orchestrations.ZeusAiPlatformWorker.ExecuteAiTelegramWorkflow(c.Request().Context(), ou, msgs)
	if err != nil {
		log.Err(err).Msg("Zeus: CreateAIServiceTaskRequestHandler")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, nil)
}

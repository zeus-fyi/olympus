package zeus_webhooks

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	ai_platform_service_orchestrations "github.com/zeus-fyi/olympus/pkg/zeus/ai/orchestrations"
)

func SupportAcknowledgeSearchIndexerHandler(c echo.Context) error {
	log.Info().Msg("Zeus: SupportAcknowledgeSearchIndexerHandler")
	request := new(AIServiceRequest)
	if err := c.Bind(request); err != nil {
		log.Err(err).Msg("SupportAcknowledgeSearchIndexerHandler")
		return err
	}
	return request.SupportAcknowledgeSearchIndexer(c)
}

func (a *AIServiceRequest) SupportAcknowledgeSearchIndexer(c echo.Context) error {
	log.Info().Msg("Zeus: SupportAcknowledgeSearchIndexer")
	err := ai_platform_service_orchestrations.ZeusAiPlatformWorker.ExecuteAiSearchIndexerWorkflow(c.Request().Context())
	if err != nil {
		log.Err(err).Msg("Zeus: ExecuteAiSearchIndexerWorkflow")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, nil)
}

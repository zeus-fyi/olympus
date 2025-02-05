package zeus_webhooks

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	ai_platform_service_orchestrations "github.com/zeus-fyi/olympus/pkg/zeus/ai/orchestrations"
)

const (
	defaultTwitterSearchGroupName = "zeusfyi"
)

func SupportAcknowledgeTwitterAiTaskRequestHandler(c echo.Context) error {
	log.Info().Msg("Zeus: SupportAcknowledgeTwitterAiTaskRequestHandler")
	request := new(AIServiceRequest)
	if err := c.Bind(request); err != nil {
		log.Err(err).Msg("SupportAcknowledgeTwitterAiTaskRequestHandler")
		return err
	}
	return request.SupportAcknowledgeTwitterAiTask(c)
}

func (a *AIServiceRequest) SupportAcknowledgeTwitterAiTask(c echo.Context) error {
	log.Info().Msg("Zeus: SupportAcknowledgeTwitterAiTask")
	group := c.Param("group")
	if len(group) == 0 {
		group = defaultTwitterSearchGroupName
	}
	internalOrgID := 7138983863666903883
	ou := org_users.NewOrgUserWithID(internalOrgID, 7138958574876245567)
	err := ai_platform_service_orchestrations.ZeusAiPlatformWorker.ExecuteAiTwitterWorkflow(c.Request().Context(), ou, group)
	if err != nil {
		log.Err(err).Msg("Zeus: ExecuteAiTwitterWorkflow")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, nil)
}

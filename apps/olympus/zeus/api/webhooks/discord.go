package zeus_webhooks

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hera_discord "github.com/zeus-fyi/olympus/pkg/hera/discord"
	ai_platform_service_orchestrations "github.com/zeus-fyi/olympus/pkg/zeus/ai/orchestrations"
)

func SupportAcknowledgeDiscordAiTaskRequestHandler(c echo.Context) error {
	log.Info().Msg("Zeus: SupportAcknowledgeDiscordAiTaskRequestHandler")
	request := new(AIServiceRequest)
	if err := c.Bind(request); err != nil {
		log.Err(err).Msg("SupportAcknowledgeDiscordAiTaskRequestHandler")
		return err
	}
	return request.SupportAcknowledgeDiscordAiTask(c)
}

func (a *AIServiceRequest) SupportAcknowledgeDiscordAiTask(c echo.Context) error {
	log.Info().Msg("Zeus: SupportAcknowledgeDiscordAiTask")
	group := c.Param("group")
	if len(group) == 0 {
		group = defaultTwitterSearchGroupName
	}
	internalOrgID := 7138983863666903883
	ou := org_users.NewOrgUserWithID(internalOrgID, 7138958574876245567)
	err := ai_platform_service_orchestrations.ZeusAiPlatformWorker.ExecuteAiFetchDataToIngestDiscordWorkflow(c.Request().Context(), ou, group)
	if err != nil {
		log.Err(err).Msg("Zeus: SupportAcknowledgeDiscordAiTask")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, nil)
}

func RequestDiscordAiTaskStartRequestHandler(c echo.Context) error {
	log.Info().Msg("Zeus: RequestDiscordAiTaskStartRequestHandler")
	request := new(DiscordRequest)
	if err := c.Bind(request); err != nil {
		log.Err(err).Msg("RequestDiscordAiTaskStartRequestHandler")
		return err
	}
	return request.RequestDiscordAiTaskStart(c)
}

type DiscordRequest struct {
	Body echo.Map `json:"body"`
}

func (a *DiscordRequest) RequestDiscordAiTaskStart(c echo.Context) error {
	internalOrgID := 7138983863666903883
	ou := org_users.NewOrgUserWithID(internalOrgID, 7138958574876245567)
	b, err := json.Marshal(a.Body)
	if err != nil {
		log.Err(err).Msg("Zeus: RequestDiscordAiTaskStart")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	cms := hera_discord.ChannelMessages{}
	err = json.Unmarshal(b, &cms)
	if err != nil {
		log.Err(err).Interface("body", a.Body).Msg("Zeus: RequestDiscordAiTaskStart")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if cms.Guild.Id == "" {
		log.Err(err).Interface("body", a.Body).Msg("Zeus: RequestDiscordAiTaskStart")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if cms.Channel.Id == "" {
		log.Err(err).Interface("body", a.Body).Msg("Zeus: RequestDiscordAiTaskStart")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	err = ai_platform_service_orchestrations.ZeusAiPlatformWorker.ExecuteAiIngestDiscordWorkflow(c.Request().Context(), ou, cms)
	if err != nil {
		log.Err(err).Msg("Zeus: RequestDiscordAiTaskStart")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, nil)
}

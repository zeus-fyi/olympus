package zeus_webhooks

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/keys"
	hermes_email_notifications "github.com/zeus-fyi/olympus/pkg/hermes/email"
	ai_platform_service_orchestrations "github.com/zeus-fyi/olympus/pkg/zeus/ai/orchestrations"
)

type AIServiceRequest struct {
}

func SupportEmailAIServiceTaskRequestHandler(c echo.Context) error {
	log.Info().Msg("Zeus: SupportEmailAIServiceTaskRequestHandler")
	request := new(AIServiceRequest)
	if err := c.Bind(request); err != nil {
		log.Err(err).Msg("SupportEmailAIServiceTaskRequestHandler")
		return err
	}
	return request.SupportAcknowledgeAITask(c, "support@zeus.fyi")
}

func AlexEmailAIServiceTaskRequestHandler(c echo.Context) error {
	log.Info().Msg("Zeus: SupportEmailAIServiceTaskRequestHandler")
	request := new(AIServiceRequest)
	if err := c.Bind(request); err != nil {
		log.Err(err).Msg("SupportEmailAIServiceTaskRequestHandler")
		return err
	}
	return request.SupportAcknowledgeAITask(c, "alex@zeus.fyi")
}

func AiEmailAIServiceTaskRequestHandler(c echo.Context) error {
	log.Info().Msg("Zeus: SupportEmailAIServiceTaskRequestHandler")
	request := new(AIServiceRequest)
	if err := c.Bind(request); err != nil {
		log.Err(err).Msg("SupportEmailAIServiceTaskRequestHandler")
		return err
	}
	return request.SupportAcknowledgeAITask(c, "ai@zeus.fyi")
}

func (a *AIServiceRequest) SupportAcknowledgeAITask(c echo.Context, email string) error {
	log.Info().Msg("Zeus: CreateAIServiceTaskRequestHandler")
	ou := org_users.OrgUser{}
	key := read_keys.NewKeyReader()
	key.GetUserFromEmail(c.Request().Context(), email)
	ou = org_users.NewOrgUserWithID(key.OrgID, key.UserID)
	var emc hermes_email_notifications.GmailServiceClient
	switch email {
	case "support@zeus.fyi":
		emc = hermes_email_notifications.SupportEmailUser
	case "alex@zeus.fyi":
		emc = hermes_email_notifications.MainEmailUser
	default:
		emc = hermes_email_notifications.AIEmailUser
	}
	msgs, err := emc.GetReadEmails(email, 10)
	if err != nil {
		log.Err(err).Msg("Zeus: CreateAIServiceTaskRequestHandler")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	err = ai_platform_service_orchestrations.ZeusAiPlatformWorker.ExecuteAiTaskWorkflow(c.Request().Context(), ou, msgs)
	if err != nil {
		log.Err(err).Msg("Zeus: CreateAIServiceTaskRequestHandler")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, nil)
}

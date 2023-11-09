package zeus_webhooks

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/keys"
	hermes_email_notifications "github.com/zeus-fyi/olympus/pkg/hermes/email"
	ai_platform_service_orchestrations "github.com/zeus-fyi/olympus/pkg/zeus/ai/ai/orchestrations"
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
	content := ""
	ou := org_users.OrgUser{}
	key := read_keys.NewKeyReader()

	key.GetUserFromEmail(c.Request().Context(), email)
	ou = org_users.NewOrgUserWithID(key.OrgID, key.UserID)

	msgs, err := hermes_email_notifications.SupportEmailUser.GetReadEmails(email, 10)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	for _, msg := range msgs {
		content = "write a bullet point summary of the email contents below, and then suggest some responses unless the message is from a no-reply address\n"
		content += "message is from " + msg.From + "\n"
		content += msg.Subject + "\n"
		content += msg.Body + "\n"
		fmt.Println(content)
		fmt.Println(ou.UserID, ou.OrgID)
		err = ai_platform_service_orchestrations.ZeusAiPlatformWorker.ExecuteAiTaskWorkflow(c.Request().Context(), ou, msg.From, content)
		if err != nil {
			log.Err(err).Msg("CreateAIServiceTaskRequestHandler")
			return err
		}
	}

	return c.JSON(http.StatusOK, msgs)
}

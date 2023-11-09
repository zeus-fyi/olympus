package zeus_webhooks

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	hermes_email_notifications "github.com/zeus-fyi/olympus/pkg/hermes/email"
)

type AIServiceRequest struct {
}

func SupportEmailAIServiceTaskRequestHandler(c echo.Context) error {
	log.Info().Msg("Hestia: SupportEmailAIServiceTaskRequestHandler")
	request := new(AIServiceRequest)
	if err := c.Bind(request); err != nil {
		log.Err(err).Msg("SupportEmailAIServiceTaskRequestHandler")
		return err
	}
	return request.SupportAcknowledgeAITask(c, "support@zeus.fyi")
}

func AlexEmailAIServiceTaskRequestHandler(c echo.Context) error {
	log.Info().Msg("Hestia: SupportEmailAIServiceTaskRequestHandler")
	request := new(AIServiceRequest)
	if err := c.Bind(request); err != nil {
		log.Err(err).Msg("SupportEmailAIServiceTaskRequestHandler")
		return err
	}
	return request.SupportAcknowledgeAITask(c, "alex@zeus.fyi")
}

func (a *AIServiceRequest) SupportAcknowledgeAITask(c echo.Context, email string) error {
	log.Info().Msg("Hestia: CreateAIServiceTaskRequestHandler")
	//content := ""
	//ou := org_users.OrgUser{}
	//key := read_keys.NewKeyReader()
	//if len(a.Email) <= 0 {
	//	return c.JSON(http.StatusBadRequest, nil)
	//}
	msgs, err := hermes_email_notifications.SupportEmailUser.GetReadEmails(email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	//for k, v := range a.Subject {
	//	content += k + ": " + v.(string) + "\n"
	//}
	//content = "write a bullet point summary of the email contents below, and then suggest a few different reponses they can choose to respond if it calls for a response\n"
	//content += a.Subject + "\n"
	//content += a.Body + "\n"
	//fmt.Println(a.Email)
	//fmt.Println(content)
	//fmt.Println(ou.UserID, ou.OrgID)
	//err = ai_platform_service_orchestrations.HestiaAiPlatformWorker.ExecuteAiTaskWorkflow(c.Request().Context(), ou, a.Email, content)
	//if err != nil {
	//	log.Err(err).Msg("CreateAIServiceTaskRequestHandler")
	//	return err
	//}
	return c.JSON(http.StatusOK, msgs)
}

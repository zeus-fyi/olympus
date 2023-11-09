package hestia_v1_ai

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/keys"
	ai_platform_service_orchestrations "github.com/zeus-fyi/olympus/pkg/hestia/platform/ai/orchestrations"
)

type AIServiceRequest struct {
	Email   string `json:"email"`
	Subject string `json:"content"`
	Body    string `json:"body"`
}

func CreateAIServiceTaskRequestHandler(c echo.Context) error {
	log.Info().Msg("Hestia: CreateAIServiceTaskRequestHandler")
	request := new(AIServiceRequest)
	if err := c.Bind(request); err != nil {
		log.Err(err).Msg("CreateAIServiceTaskRequestHandler")
		return err
	}
	return request.AcknowledgeAITask(c)
}

func (a *AIServiceRequest) AcknowledgeAITask(c echo.Context) error {
	log.Info().Msg("Hestia: CreateAIServiceTaskRequestHandler")
	content := ""
	ou := org_users.OrgUser{}
	key := read_keys.NewKeyReader()
	if len(a.Email) <= 0 {
		return c.JSON(http.StatusBadRequest, nil)
	}

	err := key.GetUserFromEmail(c.Request().Context(), a.Email)
	if err == nil && key.OrgID > 0 && key.UserID > 0 {
		ou = org_users.NewOrgUserWithID(key.OrgID, key.UserID)
		c.Set("orgUser", ou)
	}
	err = nil

	//for k, v := range a.Subject {
	//	content += k + ": " + v.(string) + "\n"
	//}
	content = "write a bullet point summary of the email contents below, and then suggest a few different reponses they can choose to respond if it calls for a response\n"
	content += a.Subject + "\n"
	content += a.Body + "\n"
	fmt.Println(a.Email)
	fmt.Println(content)
	fmt.Println(ou.UserID, ou.OrgID)
	err = ai_platform_service_orchestrations.HestiaAiPlatformWorker.ExecuteAiTaskWorkflow(c.Request().Context(), ou, a.Email, content)
	if err != nil {
		log.Err(err).Msg("CreateAIServiceTaskRequestHandler")
		return err
	}
	return c.JSON(http.StatusOK, "ok")
}

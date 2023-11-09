package hestia_v1_ai

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/keys"
	kronos_helix "github.com/zeus-fyi/olympus/pkg/kronos/helix"
)

type AIServiceRequest struct {
	Email   map[string]interface{} `json:"email"`
	Subject map[string]interface{} `json:"content"`
	Body    map[string]interface{} `json:"body"`
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
	for k, v := range a.Email {
		fmt.Println(k, v)

		em, ok := v.(string)
		if !ok {
			continue
		}
		key := read_keys.NewKeyReader()
		err := key.GetUserFromEmail(c.Request().Context(), em)
		if err == nil && key.OrgID > 0 && key.UserID > 0 {
			ou = org_users.NewOrgUserWithID(key.OrgID, key.UserID)
			c.Set("orgUser", ou)
		}
	}
	for k, v := range a.Subject {
		content += k + ": " + v.(string) + "\n"
	}
	for k, v := range a.Body {
		content += k + ": " + v.(string) + "\n"
	}
	fmt.Println(content)
	fmt.Println(ou.UserID, ou.OrgID)
	err := kronos_helix.KronosServiceWorker.ExecuteAiTaskWorkflow(c.Request().Context(), ou, content)
	if err != nil {
		log.Err(err).Msg("CreateAIServiceTaskRequestHandler")
		return err
	}
	return c.JSON(http.StatusOK, "ok")
}

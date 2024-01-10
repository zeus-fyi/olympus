package zeus_v1_ai

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/aegis/aws_secrets"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
)

type CreateOrUpdateAssistantRequest struct {
	openai.Assistant
}

func CreateOrUpdateAssistantRequestHandler(c echo.Context) error {
	request := new(CreateOrUpdateAssistantRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateOrUpdateAssistant(c)
}

func (t *CreateOrUpdateAssistantRequest) CreateOrUpdateAssistant(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if ou.OrgID <= 0 || ou.UserID <= 0 {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if t.Model == "" {
		return c.JSON(http.StatusBadRequest, nil)
	}

	sv, err := aws_secrets.GetMockingbirdPlatformSecrets(c.Request().Context(), ou, "openai")
	if err != nil {
		log.Err(err).Msg("failed to get openai secrets")
		return c.JSON(http.StatusUnprocessableEntity, nil)
	}
	if sv.ApiKey == "" {
		return c.JSON(http.StatusUnprocessableEntity, nil)
	}
	oac := hera_openai.InitOrgHeraOpenAI(sv.ApiKey)
	av, err := hera_openai.CreateOrUpdateAssistant(c.Request().Context(), oac, &t.Assistant)
	if err != nil {
		log.Err(err).Msg("failed to create or update assistant")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if av == nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	err = artemis_orchestrations.InsertAssistant(c.Request().Context(), ou, *av)
	if err != nil {
		log.Err(err).Msg("failed to insert task")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, nil)
}

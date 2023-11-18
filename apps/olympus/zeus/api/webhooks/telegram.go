package zeus_webhooks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_access_keygen "github.com/zeus-fyi/olympus/hestia/web/access"
	artemis_hydra_orchestrations_aws_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/aws_auth"
	ai_platform_service_orchestrations "github.com/zeus-fyi/olympus/pkg/zeus/ai/ai/orchestrations"
)

var cah = cache.New(time.Hour, cache.DefaultExpiration)

func (a *AIServiceRequest) SupportAcknowledgeTelegramAiTask(c echo.Context) error {
	log.Info().Msg("Zeus: SupportAcknowledgeTelegramAiTask")
	group := c.Param("group")
	if len(group) == 0 {
		return c.JSON(http.StatusBadRequest, nil)
	}
	internalOrgID := 7138983863666903883
	ou := org_users.NewOrgUserWithID(internalOrgID, 7138958574876245567)
	sv, err := artemis_hydra_orchestrations_aws_auth.GetOrgSecret(c.Request().Context(), hestia_access_keygen.FormatSecret(ou.OrgID))
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg(fmt.Sprintf("%s", err.Error()))
		return c.JSON(http.StatusInternalServerError, nil)
	}
	m := make(map[string]hestia_access_keygen.SecretsKeyValue)
	err = json.Unmarshal(sv, &m)
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg(fmt.Sprintf("%s", err.Error()))
		return c.JSON(http.StatusInternalServerError, nil)
	}

	token := ""
	tv, ok := cah.Get("telegram_token")
	if ok {
		token = tv.(string)
	}
	if len(token) == 0 {
		for k, v := range m {
			if k == "telegram" {
				if v.Key == "token" {
					cah.Set("telegram_token", v.Value, cache.DefaultExpiration)
					token = v.Value
				}
			}
		}
	}
	msgs, err := ai_platform_service_orchestrations.GetPandoraMessages(c.Request().Context(), token, group)
	if err != nil {
		cah.Delete("telegram_token")
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

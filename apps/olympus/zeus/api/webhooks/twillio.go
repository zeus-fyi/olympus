package zeus_webhooks

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func SupportAcknowledgeTwillioTaskHandler(c echo.Context) error {
	log.Info().Msg("Zeus: SupportAcknowledgeTwillioTask")
	request := new(json.RawMessage)
	if err := c.Bind(request); err != nil {
		log.Err(err).Msg("SupportAcknowledgeTwillioTask")
		return err
	}
	return SupportAcknowledgeTwillioTask(c, *request)
}

func SupportAcknowledgeTwillioTask(c echo.Context, request json.RawMessage) error {
	log.Info().Msg("Zeus: SupportAcknowledgeTwillioTask")
	group := c.Param("group")
	if len(group) == 0 {
		return c.JSON(http.StatusBadRequest, nil)
	}
	internalOrgID := 7138983863666903883
	ou := org_users.NewOrgUserWithID(internalOrgID, 7138958574876245567)
	urw := &artemis_entities.UserEntityWrapper{
		UserEntity: artemis_entities.UserEntity{
			Nickname: "sms_acknowledgement",
			Platform: "twillio",
			MdSlice: []artemis_entities.UserEntityMetadata{
				{
					JsonData: request,
				},
			},
		},
		Ou: ou,
	}
	err := artemis_entities.InsertUserEntityLabeledMetadata(c.Request().Context(), urw)
	if err != nil {
		log.Err(err).Msg("Zeus: CreateAIServiceTaskRequestHandler")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, nil)
}

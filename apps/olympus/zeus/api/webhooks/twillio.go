package zeus_webhooks

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/twilio/twilio-go/twiml"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func SupportAcknowledgeTwillioTaskHandler(c echo.Context) error {
	log.Info().Msg("Zeus: SupportAcknowledgeTwillioTask")
	request := new(twiml.MessagingMessage)
	if err := c.Bind(request); err != nil {
		log.Err(err).Msg("SupportAcknowledgeTwillioTask")
		return err
	}
	return SupportAcknowledgeTwillioTask(c, *request)
}

func SupportAcknowledgeTwillioTask(c echo.Context, request twiml.MessagingMessage) error {
	log.Info().Msg("Zeus: SupportAcknowledgeTwillioTask")
	internalOrgID := 7138983863666903883
	ou := org_users.NewOrgUserWithID(internalOrgID, 7138958574876245567)

	j, err := json.Marshal(request)
	if err != nil {
		log.Err(err).Msg("Zeus: CreateAIServiceTaskRequestHandler")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	urw := &artemis_entities.UserEntityWrapper{
		UserEntity: artemis_entities.UserEntity{
			Nickname: request.From,
			Platform: "twillio",
			MdSlice: []artemis_entities.UserEntityMetadata{
				{
					TextData: aws.String(request.Body),
					JsonData: j,
					Labels: []artemis_entities.UserEntityMetadataLabel{
						{
							Label: "from:" + request.From,
						},
						{
							Label: "to:" + request.To,
						},
						{
							Label: "twillio",
						},
					},
				},
			},
		},
		Ou: ou,
	}
	err = artemis_entities.InsertUserEntityLabeledMetadata(c.Request().Context(), urw)
	if err != nil {
		log.Err(err).Msg("Zeus: CreateAIServiceTaskRequestHandler")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, nil)
}

package zeus_webhooks

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/keys"
	hermes_email_notifications "github.com/zeus-fyi/olympus/pkg/hermes/email"
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
	msgs, err := emc.GetReadEmails(email, 25)
	if err != nil {
		log.Err(err).Msg("Zeus: CreateAIServiceTaskRequestHandler")
		return c.JSON(http.StatusInternalServerError, nil)
	}

	for _, msg := range msgs {
		bd, brr := json.Marshal(msg)
		if brr != nil {
			log.Err(brr).Msg("Zeus: SupportAcknowledgeAITask")
			return c.JSON(http.StatusInternalServerError, nil)
		}
		urw := &artemis_entities.UserEntityWrapper{
			UserEntity: artemis_entities.UserEntity{
				Nickname: email,
				Platform: "email",
				MdSlice: []artemis_entities.UserEntityMetadata{
					{
						JsonData: bd,
						TextData: aws.String(msg.Body),
						Labels: []artemis_entities.UserEntityMetadataLabel{
							{
								Label: "from:" + msg.From,
							},
							{
								Label: "to:" + email,
							},
							{
								Label: "action:respond:" + HashContents(msg.Body),
							},
							{
								Label: "indexer:email",
							},
							{
								Label: "email",
							},
							{
								Label: fmt.Sprintf("email:id:%d", msg.MsgId),
							},
							{
								Label: "email:subject:" + msg.Subject,
							},
							{
								Label: "mockingbird",
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
	}

	//err = ai_platform_service_orchestrations.ZeusAiPlatformWorker.ExecuteAiTaskWorkflow(c.Request().Context(), ou, msgs)
	//if err != nil {
	//	log.Err(err).Msg("Zeus: CreateAIServiceTaskRequestHandler")
	//	return c.JSON(http.StatusInternalServerError, nil)
	//}
	return c.JSON(http.StatusOK, nil)
}

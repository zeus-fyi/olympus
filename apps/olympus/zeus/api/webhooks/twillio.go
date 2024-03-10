package zeus_webhooks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/keys"
	"github.com/zeus-fyi/olympus/pkg/aegis/aws_secrets"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"golang.org/x/crypto/sha3"
)

func SupportAcknowledgeTwillioTaskHandler(c echo.Context) error {
	log.Info().Msg("Zeus: SupportAcknowledgeTwillioTask")
	return SupportAcknowledgeTwillioTask(c)
}

func SupportAcknowledgeTwillioTask(c echo.Context) error {
	log.Info().Msg("Zeus: SupportAcknowledgeTwillioTask")
	internalOrgID := 7138983863666903883

	var user, pw string
	sn := "api-twillio"
	ps, err := aws_secrets.GetMockingbirdPlatformSecrets(context.Background(), org_users.NewOrgUserWithID(internalOrgID, internalOrgID), sn)
	if len(ps.TwillioAccount) > 0 {
		user = ps.TwillioAccount
	}
	if len(ps.TwillioAuth) > 0 {
		pw = ps.TwillioAuth
	}
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: user,
		Password: pw,
	})
	ouInternal := org_users.NewOrgUserWithID(internalOrgID, internalOrgID)
	tvUnix, err := artemis_entities.SelectHighestLabelIdForLabelAndPlatform(context.Background(), ouInternal, "twillio", "indexer:twillio")
	if err != nil {
		log.Err(err).Msg("Zeus: SupportAcknowledgeTwillioTask")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	ch := chronos.Chronos{}
	tv := ch.ConvertUnixTimeStampToDate(tvUnix)
	// add query param
	resp, err := client.Api.ListMessage(&twilioApi.ListMessageParams{
		DateSentAfter: &tv,
		PageSize:      aws.Int(1000),
		Limit:         aws.Int(100),
	})
	for _, record := range resp {
		j, jerr := json.Marshal(record)
		if jerr != nil {
			log.Err(jerr).Msg("Zeus: SupportAcknowledgeTwillioTask")
			return c.JSON(http.StatusInternalServerError, nil)
		}
		log.Info().Interface("msg", aws.ToString(record.From)).Msg("Zeus: SupportAcknowledgeTwillioTask: From")
		key := read_keys.NewKeyReader()
		err = key.GetUserFromPhone(c.Request().Context(), aws.ToString(record.From))
		var ou org_users.OrgUser
		if err != nil && key.OrgID > 0 && key.UserID > 0 {
			ou = org_users.NewOrgUserWithID(key.OrgID, key.UserID)
		} else {
			log.Info().Interface("msg", record).Msg("Zeus: SupportAcknowledgeTwillioTask: no user found")
			continue
		}
		urw := &artemis_entities.UserEntityWrapper{
			UserEntity: artemis_entities.UserEntity{
				Nickname: aws.ToString(record.From),
				Platform: "twillio",
				MdSlice: []artemis_entities.UserEntityMetadata{
					{
						TextData: record.Body,
						JsonData: j,
						Labels: []artemis_entities.UserEntityMetadataLabel{
							{
								Label: "from:" + aws.ToString(record.From),
							},
							{
								Label: "to:" + aws.ToString(record.To),
							},
							{
								Label: "action:respond:" + HashContents(aws.ToString(record.Body)),
							},
							{
								Label: "indexer:twillio",
							},
							{
								Label: "twillio",
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
	return c.JSON(http.StatusOK, nil)
}

func HashContents(content string) string {
	hash := sha3.New256()
	_, _ = hash.Write([]byte(content))
	// Get the resulting encoded byte slice
	sha3v := hash.Sum(nil)
	return fmt.Sprintf("%x", hash.Sum(sha3v))
}

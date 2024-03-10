package ai_platform_service_orchestrations

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/keys"
	"github.com/zeus-fyi/olympus/pkg/aegis/aws_secrets"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
)

// if api-twillio, split sptring

func (t *ZeusWorkerTestSuite) TestDo() {
	sn := "api-twillio"
	var pw, user string
	ps, err := aws_secrets.GetMockingbirdPlatformSecrets(context.Background(), org_users.NewOrgUserWithID(internalOrgID, internalOrgID), sn)
	if len(ps.TwillioAccount) > 0 {
		user = ps.TwillioAccount
	}
	if len(ps.TwillioAuth) > 0 {
		pw = ps.TwillioAuth
	}
	t.Require().NotEmpty(user)
	t.Require().NotEmpty(pw)

	accountSid := fmt.Sprintf("AC3dd03d09b1ffc5ddff47e451b93f542a:%s", t.Tc.TwillioAuth)
	ss := strings.Split(accountSid, ":")
	fmt.Println(ss[0])
	fmt.Println(ss[1])
	r := resty.New()
	r.SetBasicAuth(ss[0], ss[1])
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + ss[0] + "/Messages.json"
	msg := "test message2"
	body := map[string]string{
		"To":   "17575828406",
		"From": "+18667953797", // Replace with a Twilio phone number from your account
		"Body": msg,
	}
	resp, err := r.R().
		SetFormData(body).
		Post(urlStr)
	fmt.Println(resp, err)
	t.Require().Nil(err)
}

func (t *ZeusWorkerTestSuite) TestGet() {
	accountSid := fmt.Sprintf("AC3dd03d09b1ffc5ddff47e451b93f542a:%s", t.Tc.TwillioAuth)

	ss := strings.Split(accountSid, ":")
	fmt.Println(ss[0])
	fmt.Println(ss[1])
	r := resty.New()
	r.SetBasicAuth(ss[0], ss[1])
	//urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + ss[0] + "/Messages.json"

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: ss[0],
		Password: ss[1],
	})

	tb := time.Now().Add(-time.Hour * 24)
	resp, err := client.Api.ListMessage(&twilioApi.ListMessageParams{
		PathAccountSid: nil,
		To:             nil,
		From:           nil,
		DateSent:       nil,
		DateSentBefore: nil,
		DateSentAfter:  &tb,
		PageSize:       nil,
		Limit:          aws.Int(10),
	})
	fmt.Println(resp, err)
	t.Require().Nil(err)
}

func (t *ZeusWorkerTestSuite) TestGetAfter() {
	apps.Pg.InitPG(context.Background(), t.Tc.ProdLocalDbPgconn)
	oi := 7138983863666903883

	var user, pw string
	sn := "api-twillio"
	ps, err := aws_secrets.GetMockingbirdPlatformSecrets(context.Background(), org_users.NewOrgUserWithID(oi, oi), sn)
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
	ouInternal := org_users.NewOrgUserWithID(oi, oi)
	tvUnix, err := artemis_entities.SelectHighestLabelIdForLabelAndPlatform(context.Background(), ouInternal, "twillio", "indexer:twillio")
	t.Require().Nil(err)
	if tvUnix == 0 {
		tvUnix = int(time.Now().Add(-time.Hour * 24 * 7).UnixNano())
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
		t.Require().Nil(jerr)
		log.Info().Interface("msg", aws.ToString(record.From)).Msg("Zeus: SupportAcknowledgeTwillioTask: From")
		key := read_keys.NewKeyReader()
		err = key.GetUserFromPhone(context.Background(), aws.ToString(record.From))
		var ou org_users.OrgUser
		fmt.Println(key.UserID)
		fmt.Println(key.OrgID)
		if err == nil && key.OrgID > 0 && key.UserID > 0 {
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
		err = artemis_entities.InsertUserEntityLabeledMetadata(context.Background(), urw)
		t.Require().Nil(err)
	}
}

/*
res, err := SelectHighestLabelIdForLabelAndPlatform(ctx, s.Ou, "twillio", "indexer:twillio")
	s.Require().Nil(err)
	//s.Require().NotZero(res)

*/

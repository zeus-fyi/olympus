package ai_platform_service_orchestrations

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/go-resty/resty/v2"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/keys"
	"github.com/zeus-fyi/olympus/pkg/aegis/aws_secrets"
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
	key := read_keys.NewKeyReader()
	err := key.GetUserFromPhone(context.Background(), "+17575828406")
	t.Require().Nil(err)
	//res, err := artemis_entities.SelectHighestLabelIdForLabelAndPlatform(ctx, t.Ou, "twillio", "indexer:twillio")
	//t.Require().Nil(err)
	//s.Require().NotZero(res)

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

/*
res, err := SelectHighestLabelIdForLabelAndPlatform(ctx, s.Ou, "twillio", "indexer:twillio")
	s.Require().Nil(err)
	//s.Require().NotZero(res)

*/

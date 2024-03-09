package ai_platform_service_orchestrations

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/go-resty/resty/v2"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

// if api-twillio, split sptring

func (t *ZeusWorkerTestSuite) TestDo() {
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

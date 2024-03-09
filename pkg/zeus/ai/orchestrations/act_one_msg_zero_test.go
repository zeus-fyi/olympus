package ai_platform_service_orchestrations

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

// if api-twillio, split sptring

func (t *ZeusWorkerTestSuite) TestDo() {
	accountSid := "AC3dd03d09b1ffc5ddff47e451b93f542a:a20cffd057418b83c996d80b31f02c48"

	ss := strings.Split(accountSid, ":")
	fmt.Println(ss[0])
	fmt.Println(ss[1])
	r := resty.New()
	r.SetBasicAuth(ss[0], ss[1])
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + ss[0] + "/Messages.json"
	msg := "test message"
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

package apollo_pagerduty

import (
	"context"
	"testing"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

var ctx = context.Background()

type PagerDutyTestSuite struct {
	pd PagerDutyClient
	test_suites_base.TestSuite
}

func (t *PagerDutyTestSuite) SetupTest() {
	t.InitLocalConfigs()
	t.pd = NewPagerDutyClient(t.Tc.PagerDutyApiKey)
}

func (t *PagerDutyTestSuite) TestSendAlert() {
	testEvent := pagerduty.V2Event{
		RoutingKey: t.Tc.PagerDutyRoutingKey,
		Action:     Action.Trigger(),
		DedupKey:   "",
		Payload: &pagerduty.V2Payload{
			Summary:   "This is a test alert",
			Source:    "PAGERDUTY_TEST",
			Severity:  Sev.Info(),
			Component: "This is a test component",
			Group:     "This is a test group",
			Class:     "Test Class",
			Details:   nil,
		},
	}

	resp, err := t.pd.SendAlert(ctx, testEvent)
	t.Require().NoError(err)
	t.Require().NotEmpty(resp)
}

func TestPagerDutyTestSuite(t *testing.T) {
	suite.Run(t, new(PagerDutyTestSuite))
}

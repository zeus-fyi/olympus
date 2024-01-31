package kronos_helix

import "github.com/PagerDuty/go-pagerduty"

func (t *KronosWorkerTestSuite) TestAlert() {
	InitPagerDutyAlertClient(t.Tc.PagerDutyApiKey)
	t.Require().NotNil(PdAlertClient)
	PdAlertGenericWfIssuesEvent.RoutingKey = t.Tc.PagerDutyRoutingKey
	ka := NewKronosActivities()
	ev := &pagerduty.V2Event{
		RoutingKey: t.Tc.PagerDutyRoutingKey,
		Action:     "trigger",
		DedupKey:   "aa349a9b-65f8-44ed-8cad-f3835f7b11e1",
		Images:     nil,
		Links:      nil,
		Client:     "your_client_name",
		ClientURL:  "your_client_url",
		Payload: &pagerduty.V2Payload{
			Summary:   "A QuickNode services workflow is stuck",
			Source:    "TEMPORAL_ALERTS",
			Severity:  "critical",
			Timestamp: "",
			Component: "This is a workflow component",
			Group:     "This is a workflow group",
			Class:     "Temporal Workflows",
			Details: map[string]interface{}{
				"details": "TEMPORAL_ALERTS",
			},
		},
	}

	err := ka.ExecuteTriggeredAlert(ctx, ev)
	t.Require().NoError(err)
}

package kronos_helix

import (
	"time"

	"github.com/PagerDuty/go-pagerduty"
	apollo_pagerduty "github.com/zeus-fyi/olympus/pkg/apollo/pagerduty"
)

type Instructions struct {
	Alerts  AlertInstructions   `json:"alerts,omitempty"`
	Trigger TriggerInstructions `json:"trigger,omitempty"`
}

var (
	PdAlertClient               apollo_pagerduty.PagerDutyClient
	PdAlertGenericWfIssuesEvent = pagerduty.V2Event{
		RoutingKey: "",
		Action:     "",
		Payload:    PdAlertGenericWfIssuesPayload,
	}
	PdAlertGenericWfIssuesPayload = &pagerduty.V2Payload{
		Summary:   "There is a stuck workflow",
		Source:    "TEMPORAL_ALERTS",
		Severity:  apollo_pagerduty.CRITICAL,
		Component: "This is a workflow component",
		Group:     "This is a workflow group",
		Class:     "Temporal Workflows",
		Details:   nil,
	}
)

// AlertAfterTime: eg set a 20 minutes duration, and for each wf watched it will check for completion <20 minutes or will alert

type TriggerInstructions struct {
	AlertAfterTime      time.Duration `json:"alertAfterTime,omitempty"`
	ResetAlertAfterTime time.Time     `json:"resetAlertAfterTime,omitempty"`
}

type AlertInstructions struct {
	Severity     apollo_pagerduty.Severity `json:"alertSeverity"`
	AlertMessage string                    `json:"alertMessage"`
	Source       string                    `json:"source"`
	Component    string                    `json:"component"`
}

func InitPagerDutyAlertClient(pdApiKey string) {
	PdAlertClient = apollo_pagerduty.NewPagerDutyClient(pdApiKey)
}

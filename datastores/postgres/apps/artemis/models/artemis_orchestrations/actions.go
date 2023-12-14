package artemis_orchestrations

import (
	"context"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

type Action struct {
	ActionID        int                     `json:"actionID"`
	ActionName      string                  `json:"actionName"`
	ActionGroupName string                  `json:"actionGroupName"`
	ActionType      string                  `json:"actionType"`
	ActionStatus    string                  `json:"actionStatus"`
	ActionMetrics   []ActionMetric          `json:"actionMetrics"`
	ActionPlatforms []ActionPlatformAccount `json:"actionPlatformAccounts"`
}

type ActionMetric struct {
	MetricName                 string  `json:"metricName"`
	MetricScoreThreshold       float64 `json:"metricScoreThreshold"`
	MetricPostActionMultiplier float64 `json:"metricPostActionMultiplier"`
}

type ActionPlatformAccount struct {
	ActionPlatformName    string `json:"actionPlatformName"`
	ActionPlatformAccount string `json:"actionPlatformAccount"`
}

func CreateOrUpdateAction(ctx context.Context, act Action) error {
	return nil
}

func SelectActions(ctx context.Context, ou org_users.OrgUser) ([]Action, error) {
	return nil, nil
}

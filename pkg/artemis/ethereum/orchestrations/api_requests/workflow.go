package artemis_api_requests

import (
	"time"

	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"go.temporal.io/sdk/workflow"
)

type ArtemisApiRequestsWorkflow struct {
	temporal_base.Workflow
	ArtemisApiRequestsActivities
}

const defaultTimeout = 6 * time.Second

func NewArtemisApiRequestsWorkflow() ArtemisApiRequestsWorkflow {
	deployWf := ArtemisApiRequestsWorkflow{
		Workflow:                     temporal_base.Workflow{},
		ArtemisApiRequestsActivities: ArtemisApiRequestsActivities{},
	}
	return deployWf
}

func (a *ArtemisApiRequestsWorkflow) GetWorkflows() []interface{} {
	return []interface{}{a.ProxyRequest}
}

type ApiProxyRequest struct {
	Url     string
	Payload []byte
}

func (a *ArtemisApiRequestsWorkflow) ProxyRequest(ctx workflow.Context, pr *ApiProxyRequest) (*ApiProxyRequest, error) {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 300,
	}
	sendCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(sendCtx, a.RelayRequest, pr).Get(sendCtx, &pr.Payload)
	if err != nil {
		log.Error("Failed to relay api request", "Error", err)
		return pr, err
	}
	return pr, err
}

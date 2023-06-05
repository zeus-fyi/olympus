package artemis_api_requests

import (
	"time"

	"github.com/labstack/echo/v4"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"go.temporal.io/sdk/temporal"
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
	Url        string
	Payload    echo.Map
	Response   echo.Map
	IsInternal bool
	Timeout    time.Duration
}

func (a *ArtemisApiRequestsWorkflow) ProxyRequest(ctx workflow.Context, pr *ApiProxyRequest) (*ApiProxyRequest, error) {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: pr.Timeout,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    50 * time.Millisecond,
			BackoffCoefficient: 1.2,
		},
	}
	sendCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(sendCtx, a.RelayRequest, pr).Get(sendCtx, &pr)
	if err != nil {
		log.Error("Failed to relay api request", "Error", err)
		return pr, err
	}
	return pr, err
}

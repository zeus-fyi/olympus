package iris_api_requests

import (
	"net/http"
	"net/url"
	"time"

	"github.com/labstack/echo/v4"
	iris_usage_meters "github.com/zeus-fyi/olympus/pkg/iris/proxy/usage_meters"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type IrisApiRequestsWorkflow struct {
	temporal_base.Workflow
	IrisApiRequestsActivities
}

const defaultTimeout = 6 * time.Second

func NewIrisApiRequestsWorkflow() IrisApiRequestsWorkflow {
	deployWf := IrisApiRequestsWorkflow{
		Workflow:                  temporal_base.Workflow{},
		IrisApiRequestsActivities: IrisApiRequestsActivities{},
	}
	return deployWf
}

func (i *IrisApiRequestsWorkflow) GetWorkflows() []interface{} {
	return []interface{}{i.ProxyRequest, i.ProxyInternalRequest,
		i.CacheRefreshAllOrgRoutingTablesWorkflow, i.CacheRefreshOrgRoutingTablesWorkflow, i.CacheRefreshOrgGroupTableWorkflow,
		i.DeleteRoutingGroupWorkflow, i.DeleteAllOrgRoutingGroupsWorkflow,
	}
}

type ApiProxyRequest struct {
	Url              string
	ServicePlan      string
	PayloadTypeREST  string
	Referrers        []string
	QueryParams      url.Values
	Payload          echo.Map
	Response         echo.Map
	RequestHeaders   http.Header
	ResponseHeaders  http.Header
	RawResponse      []byte
	StatusCode       int
	IsInternal       bool
	Timeout          time.Duration
	ReceivedAt       time.Time
	Latency          time.Duration
	PayloadSizeMeter *iris_usage_meters.PayloadSizeMeter
}

func (i *IrisApiRequestsWorkflow) ProxyRequest(ctx workflow.Context, pr *ApiProxyRequest) (*ApiProxyRequest, error) {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: pr.Timeout,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    100 * time.Millisecond,
			BackoffCoefficient: 2,
			MaximumAttempts:    20,
		},
	}
	sendCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(sendCtx, i.RelayRequest, pr).Get(sendCtx, &pr)
	if err != nil {
		log.Error("failed to relay api request", "Error", err)
		return pr, err
	}
	return pr, err
}

func (i *IrisApiRequestsWorkflow) ProxyInternalRequest(ctx workflow.Context, pr *ApiProxyRequest) (*ApiProxyRequest, error) {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: pr.Timeout,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    100 * time.Millisecond,
			BackoffCoefficient: 2,
			MaximumAttempts:    20,
		},
	}
	sendCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(sendCtx, i.InternalSvcRelayRequest, pr).Get(sendCtx, &pr)
	if err != nil {
		log.Error("failed to relay api request", "Error", err)
		return pr, err
	}
	return pr, err
}

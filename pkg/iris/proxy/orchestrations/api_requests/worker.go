package iris_api_requests

import (
	"context"

	"github.com/rs/zerolog/log"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"go.temporal.io/sdk/client"
)

type IrisApiRequestsWorker struct {
	temporal_base.Worker
}

var (
	IrisProxyWorker IrisApiRequestsWorker
	IrisCacheWorker IrisApiRequestsWorker
)

const (
	ApiRequestsTaskQueue         = "ApiRequestsTaskQueue"
	CacheUpdateRequestsTaskQueue = "CacheUpdateRequestsTaskQueue"
)

func (t *IrisApiRequestsWorker) ExecuteIrisProxyWorkflow(ctx context.Context, pr *ApiProxyRequest) (*ApiProxyRequest, error) {
	tc := t.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
	}
	txWf := NewIrisApiRequestsWorkflow()
	wf := txWf.ProxyRequest
	workflowRun, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, pr)
	if err != nil {
		log.Err(err).Msg("ExecuteIrisProxyWorkflow")
		return pr, err
	}
	err = workflowRun.Get(ctx, &pr)
	if err != nil {
		log.Err(err).Msg("Get ExecuteIrisProxyWorkflow")
		return pr, err
	}
	return pr, err
}

func (t *IrisApiRequestsWorker) ExecuteIrisInternalSvcApiProxyWorkflow(ctx context.Context, pr *ApiProxyRequest) (*ApiProxyRequest, error) {
	tc := t.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
	}
	txWf := NewIrisApiRequestsWorkflow()
	wf := txWf.ProxyInternalRequest
	workflowRun, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, pr)
	if err != nil {
		log.Err(err).Msg("ExecuteIrisProxyWorkflow")
		return pr, err
	}
	err = workflowRun.Get(ctx, &pr)
	if err != nil {
		log.Err(err).Msg("Get ExecuteIrisProxyWorkflow")
		return pr, err
	}
	return pr, err
}

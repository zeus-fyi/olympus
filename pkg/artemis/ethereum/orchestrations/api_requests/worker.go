package artemis_api_requests

import (
	"context"

	"github.com/rs/zerolog/log"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"go.temporal.io/sdk/client"
)

type ArtemisApiRequestsWorker struct {
	temporal_base.Worker
}

var (
	ArtemisProxyWorker ArtemisApiRequestsWorker
)

const (
	ApiRequestsTaskQueue = "ApiRequestsTaskQueue"
)

func (t *ArtemisApiRequestsWorker) ExecuteArtemisApiProxyWorkflow(ctx context.Context, pr *ApiProxyRequest) (*ApiProxyRequest, error) {
	tc := t.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
	}
	txWf := NewArtemisApiRequestsWorkflow()
	wf := txWf.ProxyRequest
	workflowRun, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, pr)
	if err != nil {
		log.Err(err).Msg("ExecuteArtemisApiProxyWorkflow")
		return pr, err
	}
	err = workflowRun.Get(ctx, &pr)
	if err != nil {
		log.Err(err).Msg("Get ExecuteArtemisApiProxyWorkflow")
		return pr, err
	}
	return pr, err
}

func (t *ArtemisApiRequestsWorker) ExecuteArtemisInternalSvcApiProxyWorkflow(ctx context.Context, pr *ApiProxyRequest) (*ApiProxyRequest, error) {
	tc := t.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
	}
	txWf := NewArtemisApiRequestsWorkflow()
	wf := txWf.ProxyInternalRequest
	workflowRun, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, pr)
	if err != nil {
		log.Err(err).Msg("ExecuteArtemisApiProxyWorkflow")
		return pr, err
	}
	err = workflowRun.Get(ctx, &pr)
	if err != nil {
		log.Err(err).Msg("Get ExecuteArtemisApiProxyWorkflow")
		return pr, err
	}
	return pr, err
}

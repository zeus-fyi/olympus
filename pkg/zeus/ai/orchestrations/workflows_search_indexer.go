package ai_platform_service_orchestrations

import (
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type SearchIndexerActionsRequest struct {
	Action         string                            `json:"action"`
	SearchIndexers []hera_search.SearchIndexerParams `json:"searchIndexers,omitempty"`
}

func (z *ZeusAiPlatformServiceWorkflows) AiSearchIndexerActionsWorkflow(ctx workflow.Context, wfID string, ou org_users.OrgUser, params SearchIndexerActionsRequest) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
		},
	}
	oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(ou.OrgID, wfID, "ZeusAiPlatformServiceWorkflows", "AiSearchIndexerActionsWorkflow")
	jobCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(jobCtx, "UpsertAssignment", oj).Get(jobCtx, nil)
	if err != nil {
		logger.Error("failed to update ai orch services", "Error", err)
		return err
	}

	for _, si := range params.SearchIndexers {
		if params.Action == "start" {
			si.Active = true
		} else {
			si.Active = false
		}
		mockingBirdSearchesUpdateCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(mockingBirdSearchesUpdateCtx, z.PlatformIndexerGroupStatusUpdate, ou, si).Get(mockingBirdSearchesUpdateCtx, nil)
		if err != nil {
			logger.Error("failed to update indexer services", "Error", err)
			return err
		}
	}
	finishedCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(finishedCtx, "UpdateAndMarkOrchestrationInactive", oj).Get(finishedCtx, nil)
	if err != nil {
		logger.Error("failed to mark services action complete", "Error", err)
		return err
	}

	return nil
}

func (z *ZeusAiPlatformServiceWorkflows) AiSearchIndexerWorkflow(ctx workflow.Context, wfID string) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
		},
	}
	oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(internalOrgID, wfID, "ZeusAiPlatformServiceWorkflows", "AiSearchIndexerWorkflow")
	jobCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(jobCtx, "UpsertAssignment", oj).Get(jobCtx, nil)
	if err != nil {
		logger.Error("failed to update ai orch services", "Error", err)
		return err
	}
	var sis []hera_search.SearchIndexerParams
	mockingBirdSearchesCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(mockingBirdSearchesCtx, z.SelectActiveSearchIndexerJobs).Get(mockingBirdSearchesCtx, &sis)
	if err != nil {
		logger.Error("failed to update ai orch services", "Error", err)
		return err
	}

	for _, si := range sis {
		wfIDCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(wfIDCtx, z.StartIndexingJob, si).Get(wfIDCtx, nil)
		if err != nil {
			logger.Error("failed to start job", "Error", err)
			return err
		}
	}
	finishedCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(finishedCtx, "UpdateAndMarkOrchestrationInactive", oj).Get(finishedCtx, nil)
	if err != nil {
		logger.Error("failed to update cache for qn services", "Error", err)
		return err
	}
	return nil
}

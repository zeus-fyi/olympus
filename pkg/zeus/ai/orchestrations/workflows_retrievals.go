package ai_platform_service_orchestrations

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/jackc/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const (
	lbStrategyPollTable = "poll-table"
)

func (z *ZeusAiPlatformServiceWorkflows) RetrievalsWorkflow(ctx workflow.Context, tte TaskToExecute) (*hera_search.SearchResultGroup, error) {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
			MaximumAttempts:    10,
		},
	}
	oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(tte.Ou.OrgID, tte.WfID, "ZeusAiPlatformServiceWorkflows", "RetrievalsWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if err != nil {
		logger.Error("failed to update ai orch services", "Error", err)
		return nil, err
	}

	switch tte.Wft.RetrievalPlatform {
	case twitterPlatform, redditPlatform, discordPlatform, telegramPlatform:
		retrievalCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(retrievalCtx, z.AiRetrievalTask, tte.Ou, tte.Wft, tte.Sg.Window).Get(retrievalCtx, &tte.Sg)
		if err != nil {
			logger.Error("failed to run retrieval", "Error", err)
			return nil, err
		}
		tte.Sg.SourceTaskID = tte.Wft.AnalysisTaskID
	case webPlatform:
		var routes []iris_models.RouteInfo
		retrievalWebCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(retrievalWebCtx, z.AiWebRetrievalGetRoutesTask, tte.Ou.OrgID, tte.Wft).Get(retrievalWebCtx, &routes)
		if err != nil {
			logger.Error("failed to run retrieval", "Error", err)
			return nil, err
		}
		for _, route := range routes {
			rt := RouteTask{
				Ou:          tte.Ou,
				RetrievalDB: tte.Wft.RetrievalDB,
				RouteInfo:   route,
			}
			fetchedResult := &hera_search.SearchResult{}
			retrievalWebTaskCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(retrievalWebTaskCtx, z.ApiCallRequestTask, rt).Get(retrievalWebTaskCtx, &fetchedResult)
			if err != nil {
				logger.Error("failed to run api call request task retrieval", "Error", err)
				return nil, err
			}
			if fetchedResult != nil && len(fetchedResult.WebResponse.Body) > 0 {
				tte.Sg.ApiResponseResults = append(tte.Sg.ApiResponseResults, *fetchedResult)
			}
			if fetchedResult != nil && fetchedResult.WebResponse.WebFilters != nil &&
				fetchedResult.WebResponse.WebFilters.LbStrategy != nil &&
				*fetchedResult.WebResponse.WebFilters.LbStrategy != lbStrategyPollTable {
				break
			}
		}
	case apiApproval:
		var echoResp []echo.Map
		// part 1
		payloadMaps := artemis_orchestrations.CreateMapInterfaceFromAssignedSchemaFields(tte.Ec.JsonResponseResults)
		for _, m := range payloadMaps {
			echoMap := echo.Map{}
			for k, v := range m {
				echoMap[k] = v
			}
			echoResp = append(echoResp, echoMap)
		}
		if len(echoResp) <= 0 {
			return nil, nil
		}
		// part 2
		var rets []artemis_orchestrations.RetrievalItem
		retrievalCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(retrievalCtx, z.SelectRetrievalTask, tte.Ou, tte.Tc.AIWorkflowTriggerResultApiResponse.RetrievalID).Get(retrievalCtx, &rets)
		if err != nil {
			logger.Error("failed to save trigger response", "Error", err)
			return nil, err
		}
		if len(rets) <= 0 {
			return nil, nil
		}
		retrieval := rets[0]
		if retrieval.RetrievalItemInstruction.Instructions.Status == pgtype.Null {
			return nil, nil
		}
		tte.Wft.RetrievalID = retrieval.RetrievalID
		tte.Wft.RetrievalName = retrieval.RetrievalName
		tte.Wft.RetrievalGroup = retrieval.RetrievalGroup
		tte.Wft.RetrievalInstructions = retrieval.RetrievalItemInstruction.Instructions.Bytes
		tte.Wft.RetrievalPlatform = retrieval.RetrievalPlatform
		var routes []iris_models.RouteInfo
		retrievalWebCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(retrievalWebCtx, z.AiWebRetrievalGetRoutesTask, tte.Ou.OrgID, tte.Wft).Get(retrievalWebCtx, &routes)
		if err != nil {
			logger.Error("failed to run retrieval", "Error", err)
			return nil, err
		}

		var breakTwoLoop bool
		for _, route := range routes {
			for _, payload := range echoResp {
				rt := RouteTask{
					Ou:          tte.Ou,
					RetrievalDB: tte.Wft.RetrievalDB,
					RouteInfo:   route,
					Payload:     payload,
				}
				fetchedResult := &hera_search.SearchResult{}
				// TODO, set the retry policy from instructions
				retrievalWebTaskCtx := workflow.WithActivityOptions(ctx, ao)
				err = workflow.ExecuteActivity(retrievalWebTaskCtx, z.ApiCallRequestTask, rt).Get(retrievalWebTaskCtx, &fetchedResult)
				if err != nil {
					logger.Error("failed to run api call request task retrieval", "Error", err)
					return nil, err
				}
				if fetchedResult != nil && len(fetchedResult.WebResponse.Body) > 0 {
					tte.Sg.ApiResponseResults = append(tte.Sg.ApiResponseResults, *fetchedResult)
					trrr := &artemis_orchestrations.AIWorkflowTriggerResultApiResponse{
						ApprovalID:  tte.Tc.AIWorkflowTriggerResultApiResponse.ApprovalID,
						TriggerID:   tte.Tc.AIWorkflowTriggerResultApiResponse.TriggerID,
						RetrievalID: aws.ToInt(tte.Wft.RetrievalID),
						ReqPayload:  payload,
						RespPayload: fetchedResult.WebResponse.Body,
					}
					saveApiRespCtx := workflow.WithActivityOptions(ctx, ao)
					err = workflow.ExecuteActivity(retrievalWebTaskCtx, z.SaveTriggerApiResponseOutput, trrr).Get(saveApiRespCtx, &trrr)
					if err != nil {
						logger.Error("failed to save trigger response retrieval", "Error", err)
						return nil, err
					}
				}
				if fetchedResult != nil && fetchedResult.WebResponse.WebFilters != nil &&
					fetchedResult.WebResponse.WebFilters.LbStrategy != nil &&
					*fetchedResult.WebResponse.WebFilters.LbStrategy != lbStrategyPollTable {
					breakTwoLoop = true
					break
				}
			}
			if breakTwoLoop {
				break
			}
		}
	}
	finishedCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(finishedCtx, "UpdateAndMarkOrchestrationInactive", oj).Get(finishedCtx, nil)
	if err != nil {
		logger.Error("failed to update orch for retrieval services", "Error", err)
		return nil, err
	}
	return tte.Sg, nil
}

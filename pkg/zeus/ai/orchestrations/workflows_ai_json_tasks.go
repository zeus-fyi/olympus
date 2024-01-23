package ai_platform_service_orchestrations

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const (
	AnalysisTask = "analysis"
	AggTask      = "aggregation"
)

type TaskToExecute struct {
	WfID string                                           `json:"wfID"`
	Ou   org_users.OrgUser                                `json:"ou"`
	Tc   TaskContext                                      `json:"taskContext"`
	Wft  artemis_orchestrations.WorkflowTemplateData      `json:"wft"`
	Sg   *hera_search.SearchResultGroup                   `json:"sg"`
	Wr   *artemis_orchestrations.AIWorkflowAnalysisResult `json:"wr"`
}

type TaskContext struct {
	TaskName       string `json:"taskName"`
	TaskType       string `json:"taskType"`
	ResponseFormat string `json:"responseFormat"`
	Model          string `json:"model"`
	TaskID         int    `json:"taskID"`
}

func (z *ZeusAiPlatformServiceWorkflows) JsonOutputTaskWorkflow(ctx workflow.Context, tte TaskToExecute) (*ChatCompletionQueryResponse, error) {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
			MaximumAttempts:    10,
		},
	}
	oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(tte.Ou.OrgID, tte.WfID, "ZeusAiPlatformServiceWorkflows", "JsonOutputTaskWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if err != nil {
		logger.Error("failed to update ai orch services", "Error", err)
		return nil, err
	}

	var fullTaskDef []artemis_orchestrations.AITaskLibrary
	selectTaskCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(selectTaskCtx, z.SelectTaskDefinition, tte.Ou, tte.Sg.SourceTaskID).Get(selectTaskCtx, &fullTaskDef)
	if err != nil {
		logger.Error("failed to run task", "Error", err)
		return nil, err
	}
	if len(fullTaskDef) == 0 {
		return nil, nil
	}
	var jdef []*artemis_orchestrations.JsonSchemaDefinition
	for _, taskDef := range fullTaskDef {
		jdef = append(jdef, taskDef.Schemas...)
	}
	fd := artemis_orchestrations.ConvertToFuncDef(tte.Tc.TaskName, jdef)
	jsonTaskCtx := workflow.WithActivityOptions(ctx, ao)
	maxAttempts := ao.RetryPolicy.MaximumAttempts
	var aiResp *ChatCompletionQueryResponse
	for attempt := 0; attempt < int(maxAttempts); attempt++ {
		jsonTaskCtx = workflow.WithActivityOptions(ctx, ao)
		params := hera_openai.OpenAIParams{
			Model:              tte.Tc.Model,
			Prompt:             tte.Sg.GetPromptBody(),
			FunctionDefinition: fd,
		}
		jsd := artemis_orchestrations.ConvertToJsonSchema(params.FunctionDefinition)
		tte.Wr.IterationCount = attempt
		err = workflow.ExecuteActivity(jsonTaskCtx, z.CreateJsonOutputModelResponse, tte.Ou, params).Get(jsonTaskCtx, &aiResp)
		if err != nil {
			logger.Error("failed to run analysis json", "Error", err)
			return nil, err
		}
		analysisCompCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(analysisCompCtx, z.RecordCompletionResponse, tte.Ou, aiResp).Get(analysisCompCtx, &aiResp.ResponseTaskID)
		if err != nil {
			logger.Error("failed to record completion response", "Error", err)
			return nil, err
		}
		wr := tte.Wr
		wr.SourceTaskID = tte.Tc.TaskID
		wr.IterationCount = attempt
		wr.ResponseID = aiResp.ResponseTaskID
		var m any
		var anyErr error
		if len(aiResp.Response.Choices) > 0 && len(aiResp.Response.Choices[0].Message.ToolCalls) > 0 {
			m, anyErr = UnmarshallOpenAiJsonInterfaceSlice(params.FunctionDefinition.Name, aiResp)
			// ok no err
			if anyErr != nil {
				log.Err(anyErr).Interface("m", m).Msg("1_UnmarshallFilteredMsgIdsFromAiJson: UnmarshallOpenAiJsonInterfaceSlice failed")
			}
		} else {
			m, anyErr = UnmarshallOpenAiJsonInterface(params.FunctionDefinition.Name, aiResp)
			if anyErr != nil {
				log.Err(anyErr).Interface("m", m).Msg("2_UnmarshallFilteredMsgIdsFromAiJson: UnmarshallOpenAiJsonInterface failed")
			}
		}
		var tmpResp [][]*artemis_orchestrations.JsonSchemaDefinition
		if anyErr == nil {
			tmpResp, anyErr = artemis_orchestrations.AssignMapValuesMultipleJsonSchemasSlice(jsd, m)
		}
		if anyErr != nil {
			log.Err(anyErr).Interface("m", m).Msg("JsonOutputTaskWorkflow: AssignMapValuesMultipleJsonSchemasSlice: failed")
			recordTaskCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(recordTaskCtx, z.SaveTaskOutput, tte.Wr, aiResp.Response).Get(recordTaskCtx, nil)
			if err != nil {
				logger.Error("failed to save task output", "Error", err)
				return nil, err
			}
			continue
		}
		switch tte.Tc.ResponseFormat {
		case socialMediaExtractionResponseFormat:
			mm := tte.Sg.GetMessageMap()
			seen := make(map[int]bool)
			notFound := make(map[int]int)
			duplicateCount := make(map[int]int)
			for ssi, schemas := range tmpResp {
				for si, sch := range schemas {
					for findex, fi := range sch.Fields {
						switch fi.FieldName {
						case "msg_id":
							msgID := aws.IntValue(fi.IntValue)
							if _, ok := seen[msgID]; ok {
								duplicateCount[msgID]++
								continue
							}
							if srv, ok := mm[aws.IntValue(fi.IntValue)]; ok {
								srv.Verified = &ok
								seen[msgID] = true
								tte.Sg.FilteredSearchResultMap[msgID] = srv
								tmpResp[ssi][si].Fields[findex].IsValidated = ok
							} else {
								notFound[msgID]++
							}
						case "msg_ids":
							for _, msgID := range fi.IntValueSlice {
								if _, ok := seen[msgID]; ok {
									duplicateCount[msgID]++
									continue
								}
								if srv, ok := mm[msgID]; ok {
									srv.Verified = &ok
									seen[msgID] = true
									tte.Sg.FilteredSearchResultMap[msgID] = srv
									tmpResp[ssi][si].Fields[findex].IsValidated = ok
								} else {
									notFound[msgID]++
								}
							}
						}
					}
				}
			}
			if len(notFound) > 0 {
				logger.Info("JsonOutputTaskWorkflow: socialMediaExtractionResponseFormat", "notFound", notFound)
			}
			if len(duplicateCount) > 0 {
				logger.Info("JsonOutputTaskWorkflow: socialMediaExtractionResponseFormat", "duplicateCount", duplicateCount)
			}
			logger.Info("JsonOutputTaskWorkflow: socialMediaExtractionResponseFormatStats", "seen", len(seen), "notFound", len(notFound), "duplicateCount", len(duplicateCount))
			aiResp.JsonResponseResults = append(aiResp.JsonResponseResults, tmpResp...)
		default:
			aiResp.JsonResponseResults = append(aiResp.JsonResponseResults, tmpResp...)
		}
		aiResp.JsonResponseResults = append(aiResp.JsonResponseResults, tmpResp...)
		recordTaskCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(recordTaskCtx, z.SaveTaskOutput, tte.Wr, aiResp.JsonResponseResults).Get(recordTaskCtx, nil)
		if err != nil {
			logger.Error("failed to save task output", "Error", err)
			return nil, err
		}
		break
	}

	finishedCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(finishedCtx, "UpdateAndMarkOrchestrationInactive", oj).Get(finishedCtx, nil)
	if err != nil {
		logger.Error("failed to update cache for qn services", "Error", err)
		return nil, err
	}
	return aiResp, nil
}

package artemis_orchestrations

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type AggregatedData struct {
	AIWorkflowAnalysisResult
	TaskName         string          `json:"taskName"`
	TaskType         string          `json:"taskType"`
	Model            string          `json:"model"`
	Prompt           json.RawMessage `json:"prompt,omitempty"`
	PromptTokens     int             `json:"promptTokens"`
	CompletionTokens int             `json:"completionTokens"`
	TotalTokens      int             `json:"totalTokens"`
}

type OrchestrationsAnalysis struct {
	TotalWorkflowTokenUsage              int                         `db:"total_workflow_token_usage" json:"totalWorkflowTokenUsage"`
	RunCycles                            int                         `db:"max_run_cycle" json:"runCycles"`
	CompleteApiRequests                  int                         `json:"completeApiRequests"`
	TotalApiRequests                     int                         `json:"totalApiRequests"`
	TotalCsvCells                        int                         `json:"totalCsvCells"`
	Progress                             float64                     `json:"progress"`
	AggregatedData                       []AggregatedData            `db:"aggregated_data" json:"aggregatedData,omitempty"`
	AggregatedEvalResults                []EvalMetric                `json:"aggregatedEvalResults,omitempty"`
	AggregatedRetrievalResults           []AIWorkflowRetrievalResult `json:"aggregatedRetrievalResults,omitempty"`
	artemis_autogen_bases.Orchestrations `json:"orchestration,omitempty"`
}

func SelectAiSystemOrchestrations(ctx context.Context, ou org_users.OrgUser, rid int) ([]OrchestrationsAnalysis, error) {
	var ojs []OrchestrationsAnalysis
	q := sql_query_templates.QueryParams{}
	args := []interface{}{ou.OrgID}

	var limit string
	queryByRunID := ""
	if rid > 0 {
		queryByRunID = " AND ar.workflow_run_id = $2"
		args = append(args, rid)
	} else {
		limit = " LIMIT 100"
	}

	// uses main for unique id, so type == real name for related workflow
	q.RawQuery = `WITH cte_a AS (
						SELECT
							o.orchestration_id,
							o.orchestration_name,
							o.group_name AS group_name,
							o.type AS orchestration_type,
							o.active,
							o.org_id,
							ar.total_api_requests,
							ar.total_csv_cells
						FROM 
							public.ai_workflow_runs AS ar 
					  	JOIN
							public.orchestrations AS o ON o.orchestration_id = ar.orchestration_id
						WHERE 
							o.org_id = $1 AND ar.is_archived = false` + queryByRunID + ` 
						ORDER BY
							o.orchestration_id DESC
						` + limit + `
					), cte_ret_status AS (
						SELECT
							o.orchestration_id,
							o.orchestration_name,
							o.group_name,
							o.orchestration_type,
							o.active,
							ai_io.workflow_result_id,
							ai_io.status,
							ai_io.retrieval_id AS retrieval_id, 
							ai_io.iteration_count,
							ai_io.chunk_offset,
							ai_io.running_cycle_number,
							ai_io.skip_retrieval AS skip_retrieval, 
							ai_io.search_window_unix_start,
							ai_io.search_window_unix_end,
							ai_io.metadata
						FROM 
							cte_a o 
						JOIN
							public.ai_workflow_io_results ai_io ON ai_io.orchestration_id = o.orchestration_id
						WHERE 
							o.org_id = $1 
						GROUP BY                             
							o.orchestration_id,
							o.orchestration_name,
							o.group_name,
							o.orchestration_type,
							o.active,
							ai_io.workflow_result_id,
							ai_io.retrieval_id,
							ai_io.iteration_count,
							ai_io.chunk_offset,
							ai_io.running_cycle_number,
							ai_io.skip_retrieval,
							ai_io.search_window_unix_start,
							ai_io.search_window_unix_end,
							ai_io.metadata
					), cte_0 AS (
						SELECT
							o.orchestration_id,
							o.orchestration_name,
							o.group_name,
							o.orchestration_type,
							o.active,
							ai_res.workflow_result_id,
							ai_res.response_id,
							ai_res.source_task_id,
							ai_res.iteration_count,
							ai_res.chunk_offset,
							ai_res.running_cycle_number,
							ai_res.skip_analysis,
							ai_res.search_window_unix_start,
							ai_res.search_window_unix_end,
							ai_res.metadata
						FROM 
							cte_a o 
						JOIN
							public.ai_workflow_analysis_results ai_res ON ai_res.orchestration_id = o.orchestration_id
						WHERE 
							o.org_id = $1 
						GROUP BY 								
							o.orchestration_id,
							o.orchestration_name,
							o.group_name,
							o.orchestration_type,
							o.active,
							ai_res.workflow_result_id,
							ai_res.response_id,
							ai_res.source_task_id,
							ai_res.iteration_count,
							ai_res.chunk_offset,
							ai_res.running_cycle_number,
							ai_res.skip_analysis,
							ai_res.search_window_unix_start,
							ai_res.search_window_unix_end,
							ai_res.metadata
						), cte_00 AS (
							SELECT 
								ai_res.orchestration_id,
								MAX(ai_res.running_cycle_number) AS max_run_cycle,
								SUM(comp_resp.total_tokens) AS total_workflow_token_usage
							FROM cte_0 ai_res
							JOIN 
								public.completion_responses AS comp_resp ON comp_resp.response_id = ai_res.response_id
							GROUP BY
								ai_res.orchestration_id
						), cte_ret_agg AS (
							SELECT 
								ai_io.orchestration_id,
								ai_io.orchestration_type,
								ai_io.orchestration_name,
								ai_io.active,
								ai_io.group_name,
								JSONB_AGG(
									JSON_BUILD_OBJECT(
										'orchestrationID', ai_io.orchestration_id,
										'workflowResultID', ai_io.workflow_result_id,
										'workflowResultStrID', ai_io.workflow_result_id::text,
										'sourceRetrievalID', ai_io.retrieval_id,  
										'iterationCount', ai_io.iteration_count,
										'status', ai_io.status,
										'chunkOffset', ai_io.chunk_offset,
										'skipRetrieval', ai_io.skip_retrieval, 
										'retrievalName', task_lib.retrieval_name,
										'retrievalPlatform', task_lib.retrieval_platform,
										'retrievalGroup', task_lib.retrieval_group,
										'runningCycleNumber', ai_io.running_cycle_number,
										'searchWindowUnixStart', ai_io.search_window_unix_start,
										'searchWindowUnixEnd', ai_io.search_window_unix_end,
										'metadata', ai_io.metadata
									) ORDER BY ai_io.running_cycle_number DESC, ai_io.iteration_count DESC, ai_io.retrieval_id DESC
								) AS ret_aggregated_data
							FROM cte_ret_status ai_io  
							JOIN 
								public.ai_retrieval_library AS task_lib ON task_lib.retrieval_id = ai_io.retrieval_id 
							GROUP BY
								ai_io.orchestration_id, ai_io.orchestration_type, ai_io.group_name, ai_io.orchestration_name, ai_io.active
						), cte_1 AS (
							SELECT 
							ai_res.orchestration_id,
							ai_res.orchestration_type,
							ai_res.orchestration_name,
							ai_res.active,
							ai_res.group_name,
									JSONB_AGG(
										JSON_BUILD_OBJECT(
											'orchestrationID', ai_res.orchestration_id,
											'workflowResultStrID', ai_res.workflow_result_id::text,
											'workflowResultID', ai_res.workflow_result_id,
											'responseID', ai_res.response_id,
											'sourceTaskID', ai_res.source_task_id,
											'iterationCount', ai_res.iteration_count,
											'chunkOffset', ai_res.chunk_offset,
											'skipAnalysis', ai_res.skip_analysis,
											'taskName', task_lib.task_name,
											'taskType', task_lib.task_type,
											'model', task_lib.model,
											'runningCycleNumber', ai_res.running_cycle_number,
											'searchWindowUnixStart', ai_res.search_window_unix_start,
											'searchWindowUnixEnd', ai_res.search_window_unix_end,
											'metadata', ai_res.metadata,
											'completionChoices', comp_resp.completion_choices,
											'prompt', comp_resp.prompt,
											'promptTokens', comp_resp.prompt_tokens,
											'completionTokens', comp_resp.completion_tokens,
											'totalTokens', comp_resp.total_tokens
										) ORDER BY ai_res.running_cycle_number DESC, ai_res.iteration_count DESC, ai_res.response_id DESC
									) AS aggregated_data
							FROM cte_0 ai_res
							JOIN 
								public.ai_task_library AS task_lib ON task_lib.task_id = ai_res.source_task_id
							JOIN 
								public.completion_responses AS comp_resp ON comp_resp.response_id = ai_res.response_id
							GROUP BY
								ai_res.orchestration_id,ai_res.orchestration_type,ai_res.group_name,ai_res.orchestration_name,ai_res.active
							), cte_2 AS (
								SELECT
									ai_res.orchestration_id,
									ai_res.running_cycle_number,
									JSONB_AGG(
										CASE 
											WHEN eval_res.eval_metrics_results_id IS NOT NULL THEN
												JSONB_BUILD_OBJECT(
													'evalName', ef.eval_name,
													'evalField', JSON_BUILD_OBJECT(
															'fieldName', af.field_name,
															'dataType', af.data_type
													),
													'evalMetricID', eval_met.eval_metric_id,	
													'evalExpectedResultState', eval_met.eval_metric_result,
													'evalMetricResult', JSON_BUILD_OBJECT(
															'evalMetricsResultID', eval_res.eval_metrics_results_id,
															'evalResultOutcomeBool', eval_res.eval_result_outcome,
															'runningCycleNumber', eval_res.running_cycle_number,
															'evalIterationCount', eval_res.eval_iteration_count,
															'searchWindowUnixStart', eval_res.search_window_unix_start,
															'searchWindowUnixEnd', eval_res.search_window_unix_end,
															'evalMetadata', eval_res.eval_metadata
													),
													'evalMetricComparisonValues', JSON_BUILD_OBJECT(
														'evalComparisonInteger', eval_met.eval_comparison_integer,
														'evalComparisonBoolean', eval_met.eval_comparison_boolean,
														'evalComparisonNumber', eval_met.eval_comparison_number,
														'evalComparisonString', eval_met.eval_comparison_string
													),	
													'evalOperator', eval_met.eval_operator,
													'evalState', eval_met.eval_state
												) 
										END
										ORDER BY eval_res.running_cycle_number DESC, eval_res.eval_metrics_results_id DESC
									) AS aggregated_eval_results
								FROM 
									cte_0 AS ai_res
								JOIN 
									public.eval_metrics_results AS eval_res ON eval_res.orchestration_id = ai_res.orchestration_id
								JOIN 
									public.eval_metrics AS eval_met ON eval_met.eval_metric_id = eval_res.eval_metric_id
								JOIN 
									public.ai_fields AS af ON af.field_id = eval_met.field_id
								JOIN 
									public.eval_fns AS ef ON ef.eval_id = eval_met.eval_id
								GROUP BY 
									ai_res.orchestration_id, ai_res.running_cycle_number
								ORDER BY 
									ai_res.orchestration_id DESC
							)
							SELECT 
								ca.orchestration_id,
								ca.orchestration_id::text,
								ca.orchestration_name,
								ca.group_name,
								ca.orchestration_type,
								ca.active,
								ca.total_api_requests,
								ca.total_csv_cells,
							  	COALESCE(c00.max_run_cycle, 0),
							  	COALESCE(c00.total_workflow_token_usage, 0) AS total_workflow_token_usage,
								COALESCE(aggregated_data, '[]'::jsonb) AS aggregated_data,
								COALESCE(aggregated_eval_results, '[]'::jsonb) AS aggregated_eval_results,
								COALESCE(cret.ret_aggregated_data, '[]'::jsonb) AS ret_aggregated_data
							 FROM cte_a ca 
							 LEFT JOIN cte_ret_agg cret ON cret.orchestration_id = ca.orchestration_id
							 LEFT JOIN cte_1 c1 ON ca.orchestration_id = c1.orchestration_id
							 LEFT JOIN cte_00 c00 ON c00.orchestration_id = c1.orchestration_id
							 LEFT JOIN cte_2 c2 ON c2.orchestration_id = c1.orchestration_id
							 GROUP BY ca.orchestration_id,
								ca.orchestration_id::text,
								ca.orchestration_name,
								ca.group_name,
								ca.orchestration_type,
								ca.active, c00.max_run_cycle, c00.total_workflow_token_usage,
								aggregated_data, aggregated_eval_results, ret_aggregated_data,
 								ca.total_api_requests,
								ca.total_csv_cells
							  ORDER BY ca.orchestration_id DESC;`

	log.Debug().Interface("SelectSystemOrchestrationsWithInstructionsByGroup", q.LogHeader(Orchestrations))
	rows, err := apps.Pg.Query(ctx, q.RawQuery, args...)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Orchestrations)); returnErr != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		oj := OrchestrationsAnalysis{}
		var evals []EvalMetric
		var agdd []AggregatedData
		var retd []AIWorkflowRetrievalResult
		rowErr := rows.Scan(&oj.OrchestrationID, &oj.OrchestrationStrID, &oj.OrchestrationName, &oj.GroupName,
			&oj.Type, &oj.Active, &oj.TotalApiRequests, &oj.TotalCsvCells, &oj.RunCycles, &oj.TotalWorkflowTokenUsage, &agdd, &evals, &retd)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(Orchestrations))
			return nil, rowErr
		}
		oj.CompleteApiRequests = len(retd)
		for i, _ := range agdd {
			if len(agdd[i].CompletionChoices) == 4 {
				b, berr := agdd[i].CompletionChoices.MarshalJSON()
				if berr != nil {
					log.Err(berr).Msg("failed to marshal completion choices")
					return nil, berr
				}
				if string(b) == "null" {
					agdd[i].CompletionChoices = agdd[i].Metadata
				}
			}
		}
		oj.AggregatedRetrievalResults = retd
		oj.AggregatedData = agdd
		var filteredResults []EvalMetric
		seen := make(map[int]bool)
		for j, _ := range evals {
			if evals[j].EvalMetricResult == nil || aws.ToInt(evals[j].EvalMetricResult.EvalMetricResultID) == 0 {
				continue
			}
			if evals[j].EvalMetricID != nil {
				evals[j].EvalMetricStrID = aws.String(fmt.Sprintf("%d", *evals[j].EvalMetricID))
			}
			const pass = "pass"
			const cPass = "Pass"
			const fail = "fail"
			const cFail = "Fail"
			const ignore = "ignore"
			const cIgnore = "Ignore"
			if evals[j].EvalMetricResult.EvalResultOutcomeBool != nil {
				var resultBool bool
				if evals[j].EvalMetricResult.EvalResultOutcomeBool != nil {
					resultBool = *evals[j].EvalMetricResult.EvalResultOutcomeBool
				}
				switch evals[j].EvalExpectedResultState {
				case pass:
					evals[j].EvalExpectedResultState = cPass
					if resultBool {
						evals[j].EvalMetricResult.EvalResultOutcomeStateStr = aws.String(cPass)
					} else {
						evals[j].EvalMetricResult.EvalResultOutcomeStateStr = aws.String(cFail)
					}
				case fail:
					evals[j].EvalExpectedResultState = cFail
					if resultBool {
						evals[j].EvalMetricResult.EvalResultOutcomeStateStr = aws.String(cFail)
					} else {
						evals[j].EvalMetricResult.EvalResultOutcomeStateStr = aws.String(cPass)
					}
				case ignore:
					evals[j].EvalExpectedResultState = cIgnore
					evals[j].EvalMetricResult.EvalResultOutcomeStateStr = aws.String(cIgnore)
				default:
					evals[j].EvalExpectedResultState = cIgnore
					evals[j].EvalMetricResult.EvalResultOutcomeStateStr = aws.String(cIgnore)
				}
			}
			evals[j].EvalMetricResult.EvalMetricResultStrID = aws.String(fmt.Sprintf("%d", aws.ToInt(evals[j].EvalMetricResult.EvalMetricResultID)))
			if _, ok := seen[aws.ToInt(evals[j].EvalMetricResult.EvalMetricResultID)]; !ok {
				filteredResults = append(filteredResults, evals[j])
				seen[aws.ToInt(evals[j].EvalMetricResult.EvalMetricResultID)] = true
			}
		}
		if oj.Active {
			mad := len(oj.AggregatedData)
			if mad > oj.TotalCsvCells {
				// alert?
				mad = oj.TotalCsvCells
			}
			rad := len(oj.AggregatedRetrievalResults)
			if rad > oj.TotalApiRequests {
				// todo alert?
				rad = oj.TotalApiRequests
			}
			oj.Progress = (float64(mad+rad) / float64(oj.TotalApiRequests+oj.TotalApiRequests)) * float64(100)
		} else {
			oj.Progress = 100
		}
		oj.AggregatedEvalResults = filteredResults
		ojs = append(ojs, oj)
	}
	var ojsRunsActions []OrchestrationsAnalysis
	for _, oj := range ojs {
		if !oj.Active && oj.RunCycles == 0 && oj.TotalWorkflowTokenUsage == 0 {
			continue
		}
		ojsRunsActions = append(ojsRunsActions, oj)
	}
	sortRunsByID(ojsRunsActions)
	return ojsRunsActions, err
}
func sortRunsByID(rs []OrchestrationsAnalysis) {
	sort.Slice(rs, func(i, j int) bool {
		return rs[i].OrchestrationID > rs[j].OrchestrationID
	})
}

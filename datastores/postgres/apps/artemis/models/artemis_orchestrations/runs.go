package artemis_orchestrations

import (
	"context"
	"encoding/json"
	"fmt"

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
	TotalWorkflowTokenUsage int              `db:"total_workflow_token_usage" json:"totalWorkflowTokenUsage"`
	RunCycles               int              `db:"max_run_cycle" json:"runCycles"`
	AggregatedData          []AggregatedData `db:"aggregated_data" json:"aggregatedData"`
	AggregatedEvalResults   []EvalMetric     `json:"aggregatedEvalResults"`

	artemis_autogen_bases.Orchestrations `json:"orchestration,omitempty"`
}

func SelectAiSystemOrchestrations(ctx context.Context, ou org_users.OrgUser) ([]OrchestrationsAnalysis, error) {
	var ojs []OrchestrationsAnalysis
	q := sql_query_templates.QueryParams{}

	// uses main for unique id, so type == real name for related workflow
	q.RawQuery = `WITH cte_0 AS (
						SELECT
							o.orchestration_id,
							o.orchestration_name,
							o.group_name AS orchestration_group_name,
							o.type AS orchestration_type,
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
							public.ai_workflow_analysis_results ai_res
						JOIN 
							public.orchestrations AS o ON o.orchestration_id = ai_res.orchestration_id
						WHERE 
							o.org_id = $1 AND
						 	(
								EXISTS (
									SELECT 1
									FROM ai_workflow_template
									WHERE workflow_name = o.type
								)
								OR EXISTS (
									SELECT 1
									FROM ai_workflow_template
									WHERE workflow_group = o.group_name
								)
							)
							GROUP BY 								
								o.orchestration_id,
								o.orchestration_name,
								o.group_name,
								o.type,
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
							), cte_1 AS (
								SELECT 
								ai_res.orchestration_id,
								ai_res.orchestration_type,
								ai_res.orchestration_name,
								ai_res.active,
								ai_res.orchestration_group_name,
										JSONB_AGG(
											JSON_BUILD_OBJECT(
												'orchestrationID', ai_res.orchestration_id,
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
									ai_res.orchestration_id,ai_res.orchestration_type,ai_res.orchestration_group_name,ai_res.orchestration_name,ai_res.active
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
								c1.orchestration_id,
								c1.orchestration_name,
								c1.orchestration_group_name,
								c1.orchestration_type,
								c1.active,
							  	c00.max_run_cycle,
							  	c00.total_workflow_token_usage,
								COALESCE(aggregated_data, '[]'::jsonb) AS aggregated_data,
								COALESCE(aggregated_eval_results, '[]'::jsonb) AS aggregated_eval_results
							 FROM cte_1 c1
							 LEFT JOIN cte_00 c00 ON c00.orchestration_id = c1.orchestration_id
							 LEFT JOIN cte_2 c2 ON c2.orchestration_id = c1.orchestration_id
							 ORDER BY orchestration_id DESC;`

	log.Debug().Interface("SelectSystemOrchestrationsWithInstructionsByGroup", q.LogHeader(Orchestrations))
	rows, err := apps.Pg.Query(ctx, q.RawQuery, ou.OrgID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Orchestrations)); returnErr != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		oj := OrchestrationsAnalysis{}

		var evals []EvalMetric
		var agdd []AggregatedData
		rowErr := rows.Scan(&oj.OrchestrationID, &oj.OrchestrationName, &oj.GroupName,
			&oj.Type, &oj.Active, &oj.RunCycles, &oj.TotalWorkflowTokenUsage, &agdd, &evals)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(Orchestrations))
			return nil, rowErr
		}
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
			if evals[j].EvalMetricResult.EvalResultOutcomeBool != nil {

				var resultBool bool
				if evals[j].EvalMetricResult.EvalResultOutcomeBool != nil {
					resultBool = *evals[j].EvalMetricResult.EvalResultOutcomeBool
				}

				if j == 8 {
					fmt.Println("evals[j].EvalMetricResult.EvalResultOutcomeBool", evals[j].EvalMetricResult.EvalResultOutcomeBool)
					fmt.Println("evals[j].EvalMetricResult.EvalResultOutcomeStateStr", evals[j].EvalMetricResult.EvalResultOutcomeStateStr)

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
				}
			}

			evals[j].EvalMetricResult.EvalMetricResultStrID = aws.String(fmt.Sprintf("%d", aws.ToInt(evals[j].EvalMetricResult.EvalMetricResultID)))
			if _, ok := seen[aws.ToInt(evals[j].EvalMetricResult.EvalMetricResultID)]; !ok {
				filteredResults = append(filteredResults, evals[j])
				seen[aws.ToInt(evals[j].EvalMetricResult.EvalMetricResultID)] = true
			}
		}

		oj.AggregatedEvalResults = filteredResults
		ojs = append(ojs, oj)
	}
	return ojs, err
}

package artemis_orchestrations

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type AggregatedData struct {
	WorkflowResultID      int             `json:"workflowResultId"`
	ResponseID            int             `json:"responseId"`
	SourceTaskID          int             `json:"sourceTaskId"`
	TaskName              string          `json:"taskName"`
	TaskType              string          `json:"taskType"`
	Model                 string          `json:"model"`
	RunningCycleNumber    int             `json:"runningCycleNumber"`
	SearchWindowUnixStart int             `json:"searchWindowUnixStart"`
	SearchWindowUnixEnd   int             `json:"searchWindowUnixEnd"`
	Metadata              json.RawMessage `json:"metadata,omitempty"`
	CompletionChoices     json.RawMessage `json:"completionChoices,omitempty"`
	Prompt                json.RawMessage `json:"prompt,omitempty"`
	PromptTokens          int             `json:"promptTokens"`
	CompletionTokens      int             `json:"completionTokens"`
	TotalTokens           int             `json:"totalTokens"`
}

type OrchestrationsAnalysis struct {
	TotalWorkflowTokenUsage int                 `db:"total_workflow_token_usage" json:"totalWorkflowTokenUsage"`
	RunCycles               int                 `db:"max_run_cycle" json:"runCycles"`
	AggregatedData          []AggregatedData    `db:"aggregated_data" json:"aggregatedData"`
	AggregatedEvalResults   []EvalMetricsResult `db:"eval_fn_metric_results" json:"aggregatedEvalResults"`

	artemis_autogen_bases.Orchestrations `json:"orchestration,omitempty"`
}

func SelectAiSystemOrchestrations(ctx context.Context, ou org_users.OrgUser) ([]OrchestrationsAnalysis, error) {
	var ojs []OrchestrationsAnalysis
	q := sql_query_templates.QueryParams{}

	// uses main for unique id, so type == real name for related workflow
	q.RawQuery = `SELECT 
						o.orchestration_id,
						o.orchestration_name AS orch_name,
						o.group_name AS orch_group_name,
						o.type AS orch_type,
						o.active,
						MAX(ai_res.running_cycle_number) AS max_run_cycle,
						SUM(comp_resp.total_tokens) AS total_workflow_token_usage,
						JSON_AGG(
							JSON_BUILD_OBJECT(
								'workflowResultId', ai_res.workflow_result_id,
								'responseId', ai_res.response_id,
								'sourceTaskId', ai_res.source_task_id,
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
							) ORDER BY ai_res.running_cycle_number DESC, ai_res.response_id DESC
						) AS aggregated_data,
						JSONB_AGG(
							CASE 
								WHEN eval_res.eval_metrics_result_id IS NOT NULL THEN
									JSON_BUILD_OBJECT(
									  	'evalName', ef.eval_name,
										'evalMetricName', eval_met.eval_metric_name,
										'evalMetricsResultId', eval_res.eval_metrics_result_id,
										'evalMetricResult', eval_met.eval_metric_result,
										'evalComparisonBoolean', eval_met.eval_comparison_boolean,
										'evalComparisonNumber', eval_met.eval_comparison_number,
										'evalComparisonString', eval_met.eval_comparison_string,
										'evalMetricDataType', eval_met.eval_metric_data_type,
										'evalOperator', eval_met.eval_operator,
										'evalState', eval_met.eval_state,
										'runningCycleNumber', eval_res.running_cycle_number,
										'searchWindowUnixStart', eval_res.search_window_unix_start,
										'searchWindowUnixEnd', eval_res.search_window_unix_end,	
										'evalResultOutcome', eval_res.eval_result_outcome
									) 
							END
							ORDER BY eval_res.running_cycle_number DESC, eval_res.eval_metrics_result_id DESC
						) AS aggregated_eval_results
					FROM 
						public.ai_workflow_analysis_results AS ai_res
					LEFT JOIN 
						public.ai_task_library AS task_lib ON task_lib.task_id = ai_res.source_task_id
					LEFT JOIN 
						public.completion_responses AS comp_resp ON comp_resp.response_id = ai_res.response_id
					LEFT JOIN 
						public.eval_metrics_results AS eval_res ON eval_res.orchestration_id = ai_res.orchestrations_id
					LEFT JOIN 
						public.eval_metrics AS eval_met ON eval_met.eval_metric_id = eval_res.eval_metric_id
					LEFT JOIN 
						public.eval_fns AS ef ON ef.eval_id = eval_met.eval_id
					JOIN 
    					public.orchestrations AS o ON o.orchestration_id = ai_res.orchestrations_id
									WHERE o.org_id = $1
										AND (
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
						o.orchestration_id, o.orchestration_name, o.group_name, o.type, o.active, eval_res.orchestration_id, eval_res.eval_metrics_result_id 
					ORDER BY 
						o.orchestration_id DESC;`

	log.Debug().Interface("SelectSystemOrchestrationsWithInstructionsByGroup", q.LogHeader(Orchestrations))
	rows, err := apps.Pg.Query(ctx, q.RawQuery, ou.OrgID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Orchestrations)); returnErr != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		oj := OrchestrationsAnalysis{}
		var aggregatedEvalResults json.RawMessage

		rowErr := rows.Scan(&oj.OrchestrationID, &oj.OrchestrationName, &oj.GroupName,
			&oj.Type, &oj.Active, &oj.RunCycles, &oj.TotalWorkflowTokenUsage, &oj.AggregatedData, &aggregatedEvalResults)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(Orchestrations))
			return nil, rowErr
		}

		// Unmarshal the aggregated evaluation results
		var evalMetricsResults []EvalMetricsResult
		if err = json.Unmarshal(aggregatedEvalResults, &evalMetricsResults); err != nil {
			log.Err(err).Msg("Error unmarshaling aggregated evaluation results")
			return nil, err
		}

		var filteredResults []EvalMetricsResult
		seen := make(map[int]bool)
		for _, evalMetricsResult := range evalMetricsResults {
			if evalMetricsResult.EvalMetricsResultID == 0 {
				continue
			}

			if _, ok := seen[evalMetricsResult.EvalMetricsResultID]; !ok {
				filteredResults = append(filteredResults, evalMetricsResult)
				seen[evalMetricsResult.EvalMetricsResultID] = true
			}
		}

		oj.AggregatedEvalResults = filteredResults
		ojs = append(ojs, oj)
	}
	return ojs, err
}

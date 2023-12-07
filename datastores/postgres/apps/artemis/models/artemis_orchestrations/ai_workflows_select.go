package artemis_orchestrations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type WorkflowTemplateData struct {
	WfAnalysisTaskID              int     `json:"wfAnalysisTaskID"`
	AnalysisTaskID                int     `json:"analysisTaskID"`
	AnalysisCycleCount            int     `json:"analysisCycleCount"`
	AnalysisPrompt                string  `json:"analysisPrompt"`
	AnalysisModel                 string  `json:"analysisModel"`
	AnalysisTokenOverflowStrategy string  `json:"analysisTokenOverflowStrategy"`
	AnalysisTaskName              string  `json:"analysisTaskName"`
	AnalysisTaskType              string  `json:"analysisTaskType"`
	AnalysisMaxTokensPerTask      int     `json:"analysisMaxTokensPerTask"`
	AggTaskID                     *int    `json:"aggTaskID,omitempty"`
	AggCycleCount                 *int    `json:"aggCycleCount,omitempty"`
	AggTaskName                   *string `json:"aggTaskName,omitempty"`
	AggTaskType                   *string `json:"aggTaskType,omitempty"`
	AggPrompt                     *string `json:"aggPrompt,omitempty"`
	AggModel                      *string `json:"aggModel,omitempty"`
	AggTokenOverflowStrategy      *string `json:"aggTokenOverflowStrategy,omitempty"`
	AggMaxTokensPerTask           *int    `json:"aggMaxTokensPerTask,omitempty"`
	RetrievalName                 *string `json:"retrievalName,omitempty"`
	RetrievalGroup                *string `json:"retrievalGroup,omitempty"`
	RetrievalPlatform             *string `json:"retrievalPlatform,omitempty"`
	RetrievalInstructions         []byte  `json:"retrievalInstructions,omitempty"`
}

func SelectWorkflowTemplate(ctx context.Context, ou org_users.OrgUser, workflowName string) ([]WorkflowTemplateData, error) {
	var results []WorkflowTemplateData
	q := sql_query_templates.QueryParams{}
	params := []interface{}{ou.OrgID, workflowName}
	q.RawQuery = `SELECT
                    awtat.analysis_task_id,
                    awtat.task_id,
                    awtat.cycle_count as analysis_cycle_count,
                    ait.task_name AS analysis_task_name,
                    ait.task_type AS analysis_task_type,
                    ait.prompt AS analysis_prompt, 
                    ait.model AS analysis_model,
                    ait.token_overflow_strategy AS analysis_token_overflow_strategy,
                    ait.max_tokens_per_task AS analysis_max_tokens_per_task,
                    ait1.task_id as agg_task_id,
                    ait1.task_name as agg_task_name,
                    ait1.task_type as agg_task_type,
                    ait1.prompt as agg_prompt,
                    ait1.model as agg_model,
                    ait1.token_overflow_strategy as agg_token_overflow_strategy,
                    ait1.max_tokens_per_task as agg_max_tokens_per_task,
                    agt.cycle_count as agg_cycle_count,
                    art.retrieval_name,
                    art.retrieval_group,
                    art.retrieval_platform as retrieval_platform,
                    art.instructions as retrieval_instructions
                FROM ai_workflow_template wate
                INNER JOIN public.ai_workflow_template_analysis_tasks awtat ON awtat.workflow_template_id = wate.workflow_template_id
                LEFT JOIN public.ai_retrieval_library art ON art.retrieval_id = awtat.retrieval_id
                LEFT JOIN public.ai_workflow_template_agg_tasks agt ON agt.analysis_task_id = awtat.analysis_task_id
                LEFT JOIN public.ai_task_library ait ON ait.task_id = awtat.task_id
                LEFT JOIN public.ai_task_library ait1 ON ait1.task_id = agt.agg_task_id
                WHERE wate.org_id = $1 AND wate.workflow_name = $2`
	rows, err := apps.Pg.Query(ctx, q.RawQuery, params...)
	if err != nil {
		log.Err(err).Msg("Error querying SelectWorkflowTemplate")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var data WorkflowTemplateData
		rowErr := rows.Scan(
			&data.WfAnalysisTaskID,
			&data.AnalysisTaskID,
			&data.AnalysisCycleCount,
			&data.AnalysisTaskName,
			&data.AnalysisTaskType,
			&data.AnalysisPrompt,
			&data.AnalysisModel,
			&data.AnalysisTokenOverflowStrategy,
			&data.AnalysisMaxTokensPerTask,
			&data.AggTaskID,
			&data.AggTaskName,
			&data.AggTaskType,
			&data.AggPrompt,
			&data.AggModel,
			&data.AggTokenOverflowStrategy,
			&data.AggMaxTokensPerTask,
			&data.AggCycleCount,
			&data.RetrievalName,
			&data.RetrievalGroup,
			&data.RetrievalPlatform,
			&data.RetrievalInstructions,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg("Error scanning row in SelectWorkflowTemplate")
			return nil, rowErr
		}
		results = append(results, data)
	}
	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		log.Err(err).Msg("error iterating over rows")
		return nil, err
	}

	return results, nil
}
func SelectWorkflowTemplates(ctx context.Context, ou org_users.OrgUser) ([]WorkflowTemplate, error) {
	var results []WorkflowTemplate
	q := sql_query_templates.QueryParams{}
	params := []interface{}{ou.OrgID}
	q.RawQuery = `SELECT
					wate.workflow_template_id,
					wate.workflow_name,
					wate.workflow_group,
					wate.fundamental_period,
					wate.fundamental_period_time_unit,
					JSON_AGG(
						JSON_BUILD_OBJECT(
							'taskID', ait.task_id,
							'taskName', ait.task_name,
							'taskType', ait.task_type,
							'model', ait.model,
							'prompt', ait.prompt,
							'cycleCount', awtat.cycle_count,
							'retrievalName', COALESCE(art.retrieval_name, 'none'),
							'retrievalPlatform', COALESCE(art.retrieval_platform, 'none')
						)
					) AS tasks,
					JSON_AGG(
						CASE 
							WHEN ait2.task_name IS NOT NULL AND ait2.task_type IS NOT NULL AND ait2.model IS NOT NULL AND ait2.prompt IS NOT NULL AND agt.cycle_count IS NOT NULL THEN
								JSON_BUILD_OBJECT(
									'retrievalName', COALESCE(ait.task_name, 'none'),
									'taskID', ait2.task_id,
									'taskName', ait2.task_name,
									'taskType', ait2.task_type,
									'model', ait2.model,
									'prompt', ait2.prompt,
									'cycleCount', agt.cycle_count
								)
							END
					) AS agg_tasks
				FROM 
					ai_workflow_template wate
				LEFT JOIN 
					ai_workflow_template_analysis_tasks awtat ON awtat.workflow_template_id = wate.workflow_template_id
				LEFT JOIN
					ai_retrieval_library art ON art.retrieval_id = awtat.retrieval_id
				LEFT JOIN 
					ai_task_library ait ON ait.task_id = awtat.task_id
				LEFT JOIN 
					ai_workflow_template_agg_tasks agt ON agt.analysis_task_id = awtat.analysis_task_id
				LEFT JOIN 
					ai_task_library ait2 ON ait2.task_id = agt.agg_task_id
				WHERE wate.org_id = $1
				GROUP BY 
					wate.workflow_template_id,
					wate.workflow_name,
					wate.workflow_group,
					wate.fundamental_period,
					wate.fundamental_period_time_unit
				ORDER BY 
					wate.workflow_name, wate.workflow_group`

	rows, err := apps.Pg.Query(ctx, q.RawQuery, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var data WorkflowTemplate
		var taskJSON string    // To store the JSON data
		var aggTaskJSON string // To store the JSON data

		rowErr := rows.Scan(
			&data.WorkflowTemplateID,
			&data.WorkflowName,
			&data.WorkflowGroup,
			&data.FundamentalPeriod,
			&data.FundamentalPeriodTimeUnit,
			&taskJSON, // Scan the JSON data into a string
			&aggTaskJSON,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg("Error scanning row in SelectWorkflowTemplate")
			return nil, rowErr
		}

		// Unmarshal the JSON data into the Tasks field
		jsonErr := json.Unmarshal([]byte(taskJSON), &data.Tasks)
		if jsonErr != nil {
			log.Err(jsonErr).Msg("Error unmarshalling task JSON")
			return nil, jsonErr
		}

		// Unmarshal the JSON data into the Tasks field
		var aggTasks []Task
		jsonErr = json.Unmarshal([]byte(aggTaskJSON), &aggTasks)
		if jsonErr != nil {
			log.Err(jsonErr).Msg("Error unmarshalling task JSON")
			return nil, jsonErr
		}
		for _, at := range aggTasks {
			if at.TaskName != "" && at.TaskType != "" && at.Model != "" && at.Prompt != "" && at.CycleCount != 0 {
				data.Tasks = append(data.Tasks, at)
			}
		}

		uniqueTasks := make(map[string]bool)
		for _, at := range data.Tasks {
			if at.TaskName != "" && at.TaskType != "" && at.Model != "" && at.Prompt != "" && at.CycleCount != 0 {
				// Create a unique key for each task. This could be a concatenation of the fields.
				taskKey := fmt.Sprintf("%s-%s-%s-%s-%d-%s-%s", at.TaskName, at.TaskType, at.Model, at.Prompt, at.CycleCount, at.RetrievalPlatform, at.RetrievalName)
				uniqueTasks[taskKey] = true
			}
		}
		var dt []Task
		for _, at := range data.Tasks {
			taskKey := fmt.Sprintf("%s-%s-%s-%s-%d-%s-%s", at.TaskName, at.TaskType, at.Model, at.Prompt, at.CycleCount, at.RetrievalPlatform, at.RetrievalName)
			if uniqueTasks[taskKey] {
				dt = append(dt, at)
				uniqueTasks[taskKey] = false
			}
		}
		data.Tasks = dt
		results = append(results, data)
	}
	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		log.Err(err).Msg("error iterating over rows")
		return nil, err
	}

	return results, nil
}

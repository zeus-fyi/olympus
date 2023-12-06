package artemis_orchestrations

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

/*
SELECT
	awtat.analysis_task_id,
	awtat.task_id,
	awtat.cycle_count as analysis_cycle_count,
	ait.task_name,
	ait.task_type,
	ait1.task_id as agg_task_id,
	agt.cycle_count as agg_cycle_count,
	art.retrieval_name,
	art.retrieval_group, art.instructions
FROM
    public.ai_workflow_template_analysis_tasks awtat
LEFT JOIN
    public.ai_retrieval_library art ON art.retrieval_id = awtat.retrieval_id
 LEFT JOIN
    public.ai_workflow_template_agg_tasks agt ON agt.analysis_task_id = awtat.analysis_task_id
LEFT JOIN
	 public.ai_task_library ait ON ait.task_id = awtat.task_id
LEFT JOIN
    public.ai_task_library ait1 ON ait1.task_id = agt.agg_task_id
WHERE
    awtat.workflow_template_id =1701667609583734016

*/

type WorkflowTemplateData struct {
	WfAnalysisTaskID              int
	AnalysisTaskID                int
	AnalysisCycleCount            int
	AnalysisPrompt                string
	AnalysisModel                 string
	AnalysisTokenOverflowStrategy string
	AnalysisTaskName              string
	AnalysisTaskType              string
	AnalysisMaxTokensPerTask      int
	AggTaskID                     *int
	AggCycleCount                 *int
	AggTaskName                   *string
	AggTaskType                   *string
	AggPrompt                     *string
	AggModel                      *string
	AggTokenOverflowStrategy      *string
	AggMaxTokensPerTask           *int
	RetrievalName                 *string
	RetrievalGroup                *string
	RetrievalPlatform             *string
	RetrievalInstructions         []byte
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
                    wate.fundamental_period_time_unit
                FROM ai_workflow_template wate
                WHERE wate.org_id = $1
                ORDER BY wate.workflow_name, wate.workflow_group ASC`
	rows, err := apps.Pg.Query(ctx, q.RawQuery, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var data WorkflowTemplate
		rowErr := rows.Scan(
			&data.WorkflowTemplateID,
			&data.WorkflowName,
			&data.WorkflowGroup,
			&data.FundamentalPeriod,
			&data.FundamentalPeriodTimeUnit,
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

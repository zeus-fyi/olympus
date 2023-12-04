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
	AnalysisTaskID     int
	TaskID             int
	AnalysisCycleCount int
	TaskName           string
	TaskType           string
	AggTaskID          *int
	AggCycleCount      *int
	RetrievalName      string
	RetrievalGroup     string
	Instructions       []byte // Assuming instructions are stored as binary data
}

func SelectWorkflowTemplate(ctx context.Context, ou org_users.OrgUser, workflowName string) ([]WorkflowTemplateData, error) {
	var results []WorkflowTemplateData

	q := sql_query_templates.QueryParams{}
	params := []interface{}{ou.OrgID, workflowName}
	q.RawQuery = `SELECT
                    awtat.analysis_task_id,
                    awtat.task_id,
                    awtat.cycle_count as analysis_cycle_count,
                    ait.task_name,
                    ait.task_type,
                    ait1.task_id as agg_task_id,
                    agt.cycle_count as agg_cycle_count,
                    art.retrieval_name,
                    art.retrieval_group, art.instructions
                FROM ai_workflow_template wate
                INNER JOIN public.ai_workflow_template_analysis_tasks awtat ON awtat.workflow_template_id = wate.workflow_template_id
                LEFT JOIN public.ai_retrieval_library art ON art.retrieval_id = awtat.retrieval_id
                LEFT JOIN public.ai_workflow_template_agg_tasks agt ON agt.analysis_task_id = awtat.analysis_task_id
                LEFT JOIN public.ai_task_library ait ON ait.task_id = awtat.task_id
                LEFT JOIN public.ai_task_library ait1 ON ait1.task_id = agt.agg_task_id
                WHERE wate.org_id = $1 AND wate.workflow_name = $2`

	rows, err := apps.Pg.Query(ctx, q.RawQuery, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var data WorkflowTemplateData
		rowErr := rows.Scan(
			&data.AnalysisTaskID,
			&data.TaskID,
			&data.AnalysisCycleCount,
			&data.TaskName,
			&data.TaskType,
			&data.AggTaskID,
			&data.AggCycleCount,
			&data.RetrievalName,
			&data.RetrievalGroup,
			&data.Instructions,
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

package artemis_orchestrations

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func DeleteWorkflowTemplates(ctx context.Context, ou org_users.OrgUser, wfs []WorkflowTemplate) error {
	q := sql_query_templates.QueryParams{}
	var wfIds []int
	for _, wf := range wfs {
		wfIds = append(wfIds, wf.WorkflowTemplateID)
	}
	params := []interface{}{ou.OrgID, pq.Array(wfIds)}
	q.RawQuery = `WITH verified_ids AS (
						SELECT workflow_template_id  
						FROM ai_workflow_template
						WHERE org_id = $1 AND workflow_template_id = ANY($2)
				  ), cte_delete_agg_tasks AS (
						DELETE FROM ai_workflow_template_agg_tasks 
						WHERE workflow_template_id IN (SELECT workflow_template_id FROM verified_ids)
				  ), cte_delete_analysis_tasks AS (
						DELETE FROM ai_workflow_template_analysis_tasks
						WHERE workflow_template_id IN (SELECT workflow_template_id FROM verified_ids)
				  ), cte_delete_eval_tasks AS (
						DELETE FROM ai_workflow_template_eval_task_relationships
						WHERE workflow_template_id IN (SELECT workflow_template_id FROM verified_ids)
				  ) DELETE FROM ai_workflow_template
					WHERE workflow_template_id IN (SELECT workflow_template_id FROM verified_ids)
				  `
	_, err := apps.Pg.Exec(ctx, q.RawQuery, params...)
	if err != nil && err != pgx.ErrNoRows {
		log.Err(err).Msg(q.LogHeader(Orchestrations))
		return err
	}

	return nil
}

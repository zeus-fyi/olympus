package artemis_orchestrations

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
)

type WorkflowTemplate struct {
	WorkflowTemplateID        int    `json:"workflowID,omitempty"`
	WorkflowName              string `json:"workflowName"`
	WorkflowGroup             string `json:"workflowGroup"`
	FundamentalPeriod         int    `json:"fundamentalPeriod"`
	FundamentalPeriodTimeUnit string `json:"fundamentalPeriodTimeUnit"`
}

type WorkflowComponentDependency struct {
	ComponentID           int
	ComponentDependencyID int
}

func InsertWorkflowTemplate(ctx context.Context, ou org_users.OrgUser, template *WorkflowTemplate) error {
	// SQL query to insert a new workflow template or update existing
	query := `
        INSERT INTO public.ai_workflow_template (workflow_name, workflow_group, org_id, user_id, fundamental_period, fundamental_period_time_unit)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (org_id, workflow_name) 
        DO UPDATE SET 
            user_id = EXCLUDED.user_id,
            fundamental_period = EXCLUDED.fundamental_period,
            fundamental_period_time_unit = EXCLUDED.fundamental_period_time_unit
        RETURNING workflow_template_id;`

	// Executing the query
	err := apps.Pg.QueryRowWArgs(ctx, query, template.WorkflowName, template.WorkflowGroup, ou.OrgID, ou.UserID, template.FundamentalPeriod, template.FundamentalPeriodTimeUnit).Scan(&template.WorkflowTemplateID)
	if err != nil {
		log.Err(err).Msg("failed to insert workflow template")
		return err
	}

	return nil
}

type AggTask struct {
	AggId      int
	CycleCount int
	Tasks      []AITaskLibrary
}

type WorkflowTasks struct {
	AggTasks          []AggTask
	AnalysisOnlyTasks []AITaskLibrary
}

func InsertWorkflowWithComponents(ctx context.Context, ou org_users.OrgUser, workflowTemplate *WorkflowTemplate, tasks WorkflowTasks) error {
	// Start a transaction
	tx, err := apps.Pg.Begin(ctx)
	if err != nil {
		return err
	}

	// Rollback in case of any error
	defer tx.Rollback(ctx)

	// Insert the workflow template and get its ID

	query := `
        INSERT INTO public.ai_workflow_template (workflow_name, workflow_group, org_id, user_id, fundamental_period, fundamental_period_time_unit)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (org_id, workflow_name) 
        DO UPDATE SET 
            user_id = EXCLUDED.user_id,
            fundamental_period = EXCLUDED.fundamental_period,
            fundamental_period_time_unit = EXCLUDED.fundamental_period_time_unit
        RETURNING workflow_template_id;`

	err = tx.QueryRow(ctx, query,
		workflowTemplate.WorkflowName, workflowTemplate.WorkflowGroup, ou.OrgID, ou.UserID, workflowTemplate.FundamentalPeriod, workflowTemplate.FundamentalPeriodTimeUnit).Scan(&workflowTemplate.WorkflowTemplateID)
	if err != nil {
		log.Err(err).Msg("failed to insert workflow template")
		return err
	}

	ts := chronos.Chronos{}
	for _, aggTask := range tasks.AggTasks {
		// Link component to the workflow template
		for _, at := range aggTask.Tasks {
			for _, rd := range at.RetrievalDependencies {
				aid := ts.UnixTimeStampNow()
				err = tx.QueryRow(ctx, `INSERT INTO ai_workflow_template_analysis_tasks(analysis_task_id, workflow_template_id, task_id, retrieval_id, cycle_count) VALUES ($1, $2, $3, $4, $5) RETURNING analysis_task_id`, aid, workflowTemplate.WorkflowTemplateID, at.TaskID, rd.RetrievalID, at.CycleCount).Scan(&aid)
				if err != nil {
					log.Err(err).Msg("failed to insert workflow component")
					return err
				}
				err = tx.QueryRow(ctx, `INSERT INTO ai_workflow_template_agg_tasks (agg_task_id, workflow_template_id, analysis_task_id, cycle_count) VALUES ($1, $2, $3, $4) RETURNING analysis_task_id`, aggTask.AggId, workflowTemplate.WorkflowTemplateID, aid, aggTask.CycleCount).Scan(&aid)
				if err != nil {
					log.Err(err).Msg("failed to insert workflow component")
					return err
				}
			}
		}
	}
	for _, at := range tasks.AnalysisOnlyTasks {
		// Link component to the workflow template
		for _, rd := range at.RetrievalDependencies {
			aid := ts.UnixTimeStampNow()
			err = tx.QueryRow(ctx, `INSERT INTO ai_workflow_template_analysis_tasks(analysis_task_id, workflow_template_id, task_id, retrieval_id, cycle_count) VALUES ($1, $2, $3, $4, $5) RETURNING analysis_task_id`, aid, workflowTemplate.WorkflowTemplateID, at.TaskID, rd.RetrievalID, at.CycleCount).Scan(&aid)
			if err != nil {
				log.Err(err).Msg("failed to insert workflow component")
				return err
			}
		}
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		log.Err(err).Msg("failed to commit transaction")
		return err
	}
	return nil
}
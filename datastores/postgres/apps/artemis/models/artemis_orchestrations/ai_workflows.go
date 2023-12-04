package artemis_orchestrations

import (
	"context"

	"github.com/jackc/pgx/v4"
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

func InsertWorkflowWithComponents(ctx context.Context, ou org_users.OrgUser, workflowTemplate *WorkflowTemplate, tasks []AITaskLibrary) error {
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

	for _, task := range tasks {

		mainTaskId := ts.UnixTimeStampNow()
		task.TaskDependencies = append(task.TaskDependencies, task)
		for i, td := range task.TaskDependencies {
			componentID := ts.UnixTimeStampNow()
			if i == 0 {
				componentID = mainTaskId
			}

			query = `INSERT INTO ai_workflow_component(component_id) VALUES($1) RETURNING component_id`
			err = tx.QueryRow(ctx, query, componentID).Scan(&componentID)
			if err != nil {
				log.Err(err).Msg("failed to insert workflow component")
				return err
			}
			// Link component to the workflow template
			_, err = tx.Exec(ctx, `INSERT INTO ai_workflow_template_components (workflow_template_id, component_id) VALUES ($1, $2)`, workflowTemplate.WorkflowTemplateID, componentID)
			if err == pgx.ErrNoRows {
				err = nil // Ignore if the component is already linked
			}
			if err != nil {
				log.Err(err).Msg("failed to insert workflow component")
				return err
			}
			// Link component to the workflow template
			_, err = tx.Exec(ctx, `INSERT INTO ai_workflow_template_component_task(component_id, task_id, cycle_count) VALUES ($1, $2, $3)`, componentID, td.TaskID, task.CycleCount)
			if err == pgx.ErrNoRows {
				err = nil // Ignore if the component is already linked
			}
			if err != nil {
				log.Err(err).Msg("failed to insert workflow component")
				return err
			}

			if i > 0 {
				// Link component to the workflow template
				_, err = tx.Exec(ctx, `INSERT INTO ai_workflow_component_dependency(component_id, component_dependency_id) VALUES ($1, $2)`, mainTaskId, componentID)
				if err == pgx.ErrNoRows {
					err = nil // Ignore if the component is already linked
				}
				if err != nil {
					log.Err(err).Msg("failed to insert workflow ai_workflow_component_dependency")
					return err
				}
			}
		}

		// fyi component of last is the main task id
		for _, rd := range task.RetrievalDependencies {
			retCompId := ts.UnixTimeStampNow()
			query = `INSERT INTO ai_workflow_component(component_id) VALUES($1) RETURNING component_id`
			err = tx.QueryRow(ctx, query, retCompId).Scan(&retCompId)
			if err != nil {
				log.Err(err).Msg("failed to insert workflow component")
				return err
			}
			// Link component to the workflow template
			_, err = tx.Exec(ctx, `INSERT INTO ai_workflow_template_components (workflow_template_id, component_id) VALUES ($1, $2)`, workflowTemplate.WorkflowTemplateID, retCompId)
			if err == pgx.ErrNoRows {
				err = nil // Ignore if the component is already linked
			}
			if err != nil {
				log.Err(err).Msg("failed to insert workflow component")
				return err
			}
			// Link component to the workflow template
			_, err = tx.Exec(ctx, `INSERT INTO ai_workflow_template_component_retrieval(component_id, retrieval_id) VALUES ($1, $2)`, retCompId, rd.RetrievalID)
			if err == pgx.ErrNoRows {
				err = nil // Ignore if the component is already linked
			}
			if err != nil {
				log.Err(err).Msg("failed to insert workflow component")
				return err
			}
			_, err = tx.Exec(ctx, `INSERT INTO ai_workflow_component_dependency(component_id, component_dependency_id) VALUES ($1, $2)`, mainTaskId, retCompId)
			if err == pgx.ErrNoRows {
				err = nil // Ignore if the component is already linked
			}
			if err != nil {
				log.Err(err).Msg("failed to insert workflow ai_workflow_component_dependency")
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

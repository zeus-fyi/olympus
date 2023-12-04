package artemis_orchestrations

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

type WorkflowTemplate struct {
	WorkflowTemplateID        int    `json:"workflowID,omitempty"`
	WorkflowName              string `json:"workflowName"`
	WorkflowGroup             string `json:"workflowGroup"`
	FundamentalPeriod         int    `json:"fundamentalPeriod"`
	FundamentalPeriodTimeUnit string `json:"fundamentalPeriodTimeUnit"`
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

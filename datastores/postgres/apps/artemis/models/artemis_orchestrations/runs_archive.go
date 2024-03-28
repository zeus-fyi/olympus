package artemis_orchestrations

import (
	"context"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func UpdateOrchestrationsToArchive(ctx context.Context, ou org_users.OrgUser, orchestrationName []string, isArchived bool) error {
	// Construct the query with the necessary join and update logic
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
    WITH updated AS (
        SELECT awr.orchestration_id
        FROM ai_workflow_runs awr
		JOIN orchestrations o ON awr.orchestration_id = o.orchestration_id
        WHERE o.org_id = $1 AND o.orchestration_name IN (SELECT * FROM UNNEST ($2::text[]))
    )
    UPDATE ai_workflow_runs
    SET is_archived = $3
    FROM updated
    WHERE ai_workflow_runs.orchestration_id = updated.orchestration_id;
    `
	// Execute the query
	_, err := apps.Pg.Exec(ctx, q.RawQuery, ou.OrgID, orchestrationName, isArchived)
	if err != nil {
		return misc.ReturnIfErr(err, q.LogHeader("Orchestrations"))
	}
	return nil
}

package artemis_orchestrations

import (
	"context"
	"fmt"

	"github.com/jackc/pgtype"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func InsertAiOrchestrations(ctx context.Context, ou org_users.OrgUser, action string, wfs []WorkflowTemplate) (int, error) {
	for _, wf := range wfs {
		fmt.Println(wf.WorkflowName, wf.WorkflowGroup, "=======")
		wtd, err := SelectWorkflowTemplate(ctx, ou, wf.WorkflowName)
		if err != nil {
			return 0, err
		}
		for _, wd := range wtd {
			fmt.Println(wd.TaskID)
		}
	}
	return 1, nil
}

func InsertOrchestrationRef(ctx context.Context, oj OrchestrationJob, b []byte) (int, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `INSERT INTO orchestrations(org_id, orchestration_name, group_name, type, instructions)
				  VALUES ($1, $2, $3, $4, $5)
				  ON CONFLICT (org_id, orchestration_name) 
				  DO UPDATE SET instructions = EXCLUDED.instructions
				  RETURNING orchestration_id;`

	var id int
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, oj.OrgID, oj.OrchestrationName, oj.GroupName, oj.Type, &pgtype.JSONB{Bytes: sanitizeBytesUTF8(b), Status: IsNull(b)}).Scan(&id)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Orchestrations)); returnErr != nil {
		return 0, err
	}
	return id, misc.ReturnIfErr(err, q.LogHeader(Orchestrations))
}

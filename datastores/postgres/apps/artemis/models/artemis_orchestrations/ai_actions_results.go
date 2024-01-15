package artemis_orchestrations

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type AIWorkflowTriggerResultResponse struct {
	TriggerResultID  int `db:"trigger_result_id" json:"triggerResultID"`
	WorkflowResultID int `db:"workflow_result_id" json:"workflowResultID"`
	TriggerID        int `db:"trigger_id" json:"triggerID"`
	ResponseID       int `db:"response_id" json:"responseID"`
}

func InsertOrUpdateAIWorkflowTriggerResultResponse(ctx context.Context, wtrr AIWorkflowTriggerResultResponse) (int, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `INSERT INTO ai_workflow_trigger_result_responses(trigger_result_id, workflow_result_id, trigger_id, response_id)
                  VALUES ($1, $2, $3, $4)
                  ON CONFLICT (trigger_result_id) 
                  DO UPDATE SET 
                      workflow_result_id = EXCLUDED.workflow_result_id,
                      trigger_id = EXCLUDED.trigger_id,
                      response_id = EXCLUDED.response_id
                  RETURNING trigger_result_id;`

	if wtrr.TriggerResultID <= 0 {
		ch := chronos.Chronos{}
		wtrr.TriggerResultID = ch.UnixTimeStampNow()
	}
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, wtrr.TriggerResultID, wtrr.WorkflowResultID, wtrr.TriggerID, wtrr.ResponseID).Scan(&wtrr.TriggerResultID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader("AIWorkflowTriggerResultResponse")); returnErr != nil {
		log.Err(returnErr).Interface("wtrr", wtrr).Msg(q.LogHeader("AIWorkflowTriggerResultResponse"))
		return 0, err
	}
	return wtrr.TriggerResultID, err
}

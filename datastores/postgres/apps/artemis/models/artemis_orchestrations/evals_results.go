package artemis_orchestrations

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type AIWorkflowEvalResultResponse struct {
	EvalResultID     int `db:"eval_result_id" json:"evalResultId"`
	WorkflowResultID int `db:"workflow_result_id" json:"workflowResultId"`
	EvalID           int `db:"eval_id" json:"evalId"`
	ResponseID       int `db:"response_id" json:"responseId"`
}

func InsertOrUpdateAiWorkflowEvalResultResponse(ctx context.Context, errr AIWorkflowEvalResultResponse) (int, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `INSERT INTO ai_workflow_eval_result_response(eval_result_id, workflow_result_id, eval_id, response_id)
                  VALUES ($1, $2, $3, $4)
                  ON CONFLICT (eval_result_id) 
                  DO UPDATE SET 
                      workflow_result_id = EXCLUDED.workflow_result_id,
                      eval_id = EXCLUDED.eval_id,
                      response_id = EXCLUDED.response_id
                  RETURNING eval_result_id;`

	if errr.EvalResultID <= 0 {
		ch := chronos.Chronos{}
		errr.EvalResultID = ch.UnixTimeStampNow()
	}
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, errr.EvalResultID, errr.WorkflowResultID, errr.EvalID, errr.ResponseID).Scan(&errr.EvalResultID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader("AIWorkflowEvalResultResponse")); returnErr != nil {
		log.Err(returnErr).Interface("errr", errr).Msg(q.LogHeader("AIWorkflowEvalResultResponse"))
		return errr.EvalResultID, err
	}
	return errr.EvalResultID, nil
}

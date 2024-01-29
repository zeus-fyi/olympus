package artemis_orchestrations

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgtype"
	"github.com/labstack/echo/v4"
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

type AIWorkflowTriggerResultApiReqResponse struct {
	ResponseID   int        `db:"response_id" json:"responseID"`
	ApprovalID   int        `db:"approval_id" json:"approvalID"`
	TriggerID    int        `db:"trigger_id" json:"triggerID"`
	RetrievalID  int        `db:"retrieval_id" json:"retrievalID"`
	ReqPayloads  []echo.Map `db:"req_payloads" json:"reqPayloads,omitempty"`
	RespPayloads []echo.Map `db:"resp_payloads" json:"respPayloads,omitempty"`
}

func InsertOrUpdateAIWorkflowTriggerResultApiResponse(ctx context.Context, wtrr *AIWorkflowTriggerResultApiReqResponse) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `INSERT INTO ai_trigger_actions_api_reqs_responses(response_id, approval_id, trigger_id, retrieval_id, req_payload, resp_payload)
                  VALUES ($1, $2, $3, $4, $5, $6)
                  ON CONFLICT (response_id, approval_id, trigger_id, retrieval_id)	 
                  DO UPDATE SET 
                      req_payload = EXCLUDED.req_payload,
                      resp_payload = EXCLUDED.resp_payload
                  RETURNING response_id;`
	if wtrr.ResponseID <= 0 {
		ch := chronos.Chronos{}
		wtrr.ResponseID = ch.UnixTimeStampNow()
	}
	var pgReqJsonB, pgRespJsonB pgtype.JSONB
	if len(wtrr.ReqPayloads) > 0 {
		b, err := json.Marshal(wtrr.ReqPayloads)
		if err != nil {
			log.Err(err).Msg("failed to marshal req payload")
			return err
		}
		pgReqJsonB.Bytes = sanitizeBytesUTF8(b)
		pgReqJsonB.Status = pgtype.Present
	} else {
		pgReqJsonB.Status = pgtype.Null
	}
	if len(wtrr.RespPayloads) > 0 {
		b, err := json.Marshal(wtrr.RespPayloads)
		if err != nil {
			log.Err(err).Msg("failed to marshal resp payload")
			return err
		}
		pgRespJsonB.Bytes = sanitizeBytesUTF8(b)
		pgRespJsonB.Status = pgtype.Present
	} else {
		pgRespJsonB.Status = pgtype.Null
	}
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, wtrr.ResponseID, wtrr.ApprovalID, wtrr.TriggerID, wtrr.RetrievalID, pgReqJsonB, pgRespJsonB).Scan(&wtrr.ResponseID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader("InsertOrUpdateAIWorkflowTriggerResultApiResponse")); returnErr != nil {
		log.Err(returnErr).Interface("wtrr", wtrr).Msg(q.LogHeader("InsertOrUpdateAIWorkflowTriggerResultApiResponse"))
		return err
	}
	return err
}

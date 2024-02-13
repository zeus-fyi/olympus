package artemis_orchestrations

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgtype"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func CreateOrUpdateTriggerActionApproval(ctx context.Context, ou org_users.OrgUser, approval TriggerActionsApproval) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `WITH cte_get_expiration AS (
					SELECT 
						expires_after_seconds
					FROM 
						ai_trigger_actions
					WHERE 
						trigger_id = $2
				)
				INSERT INTO ai_trigger_actions_approvals(
					eval_id, 
					trigger_id, 
					workflow_result_id, 
					approval_state, 
					request_summary, 
					expires_at
				)
				VALUES (
					$1, 
					$2, 
					$3, 
					$4, 
					$5, 
					CASE
						WHEN (SELECT expires_after_seconds FROM cte_get_expiration) = 0 THEN NULL
						ELSE NOW() + (SELECT expires_after_seconds FROM cte_get_expiration) * INTERVAL '1 second'
					END
				)
				ON CONFLICT (approval_id, eval_id, trigger_id, workflow_result_id)
				DO UPDATE SET 
					request_summary = EXCLUDED.request_summary,
					approval_state = EXCLUDED.approval_state
				RETURNING approval_id;`

	if approval.ApprovalState == "" {
		approval.ApprovalState = "pending"
	}
	if approval.RequestSummary == "" {
		approval.RequestSummary = "Requesting approval for trigger action"
	}
	// Executing the query
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, approval.EvalID, approval.TriggerID, approval.WorkflowResultID,
		approval.ApprovalState, approval.RequestSummary).Scan(&approval.ApprovalID)
	if err != nil {
		log.Err(err).Msg("failed to insert or update trigger action approval")
		return err
	}

	return nil
}

const (
	pendingState  = "pending"
	finishedState = "finished"
)

func CreateOrUpdateTriggerActionApprovalWithApiReq(ctx context.Context, ou org_users.OrgUser, approval TriggerActionsApproval, wtrr AIWorkflowTriggerResultApiReqResponse) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
			WITH cte_get_expiration AS (
			   SELECT 
					CASE
						WHEN expires_after_seconds = 0 THEN NULL
						ELSE NOW() + expires_after_seconds * INTERVAL '1 second'
					END AS expires_at
				FROM 
					ai_trigger_actions
				WHERE 
					trigger_id = $2
			), cte_create_approval AS (
				INSERT INTO ai_trigger_actions_approvals(approval_id, eval_id, trigger_id, workflow_result_id, approval_state, request_summary, expires_at)
				VALUES (
						$1, 
						$2, 
						$3, 
						$4, 
						$5, 
						$6,	
        				(SELECT expires_at FROM cte_get_expiration)
					)
				ON CONFLICT (approval_id, eval_id, trigger_id, workflow_result_id)
				DO UPDATE SET 
					request_summary = EXCLUDED.request_summary,
					approval_state = EXCLUDED.approval_state
				RETURNING approval_id, request_summary, approval_state
		) INSERT INTO ai_trigger_actions_api_reqs_responses(response_id, approval_id, trigger_id, retrieval_id, req_payload, resp_payload)
          SELECT $7, cte_create_approval.approval_id, $8, $9, COALESCE($10, '[]'::jsonb), COALESCE($11, '[]'::jsonb)
		  FROM cte_create_approval
		  ON CONFLICT (response_id, approval_id, trigger_id, retrieval_id)	 
		  DO UPDATE SET 
			  req_payload = EXCLUDED.req_payload,
			  resp_payload = EXCLUDED.resp_payload
		  RETURNING response_id, approval_id;
	`

	if approval.ApprovalID <= 0 {
		ch := chronos.Chronos{}
		approval.ApprovalID = ch.UnixTimeStampNow()
	}
	if wtrr.ResponseID <= 0 {
		ch := chronos.Chronos{}
		wtrr.ResponseID = ch.UnixTimeStampNow()
	}
	breq, err := json.MarshalIndent(wtrr.ReqPayloads, "", "  ")
	if err != nil {
		log.Err(err).Msg("failed to marshal req payload")
		return err
	}
	bresp, err := json.MarshalIndent(wtrr.RespPayloads, "", "  ")
	if err != nil {
		log.Err(err).Msg("failed to marshal resp payload")
		return err
	}

	if approval.ApprovalState == "" {
		approval.ApprovalState = pendingState
	}
	if approval.ApprovalState == pendingState {
		approval.RequestSummary = "Requesting approval for trigger action\n" + string(breq)
	}
	if approval.ApprovalState == finishedState {
		approval.RequestSummary = "Finished approval for trigger action\n" + string(bresp)
	}
	var returnedApprovalID int
	err = apps.Pg.QueryRowWArgs(ctx, q.RawQuery,
		approval.ApprovalID, approval.EvalID, approval.TriggerID,
		approval.WorkflowResultID, approval.ApprovalState,
		approval.RequestSummary,
		wtrr.ResponseID, approval.TriggerID, wtrr.RetrievalID,
		&pgtype.JSONB{Bytes: sanitizeBytesUTF8(breq), Status: IsNull(sanitizeBytesUTF8(breq))},
		&pgtype.JSONB{Bytes: sanitizeBytesUTF8(bresp), Status: IsNull(sanitizeBytesUTF8(bresp))}).Scan(&wtrr.ResponseID, &returnedApprovalID)
	if err != nil {
		log.Err(err).Interface("breq", breq).Interface("bresp", bresp).Interface("ou", ou).Int("returnedApprovalID", returnedApprovalID).Int("respID", wtrr.ResponseID).Msg("failed to insert or update trigger action approval for api")
		return err
	}

	return nil
}

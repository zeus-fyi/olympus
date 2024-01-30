package artemis_orchestrations

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func SelectTriggerActionApproval(ctx context.Context, ou org_users.OrgUser, state string, approvalID int) ([]TriggerActionsApproval, error) {
	if approvalID <= 0 {
		return nil, nil
	}
	var approvals []TriggerActionsApproval
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
        SELECT a.approval_id, a.approval_id::text,
               a.eval_id, a.eval_id::text,
               a.trigger_id, a.trigger_id::text,
               a.workflow_result_id, a.workflow_result_id::text, 
               a.approval_state, a.request_summary, a.updated_at,
               t.trigger_action
        FROM public.ai_trigger_actions_approvals a
        JOIN public.ai_trigger_actions t ON a.trigger_id = t.trigger_id
        WHERE t.org_id = $1 AND a.approval_state = $2 AND a.approval_id = $3
        ORDER BY a.approval_id DESC;`

	// Executing the query
	rows, err := apps.Pg.Query(ctx, q.RawQuery, ou.OrgID, state, approvalID)
	if err != nil {
		log.Err(err).Msg("failed to execute query for trigger action approvals")
		return nil, err
	}
	defer rows.Close()

	// Iterating through the query results
	for rows.Next() {
		var approval TriggerActionsApproval
		err = rows.Scan(&approval.ApprovalID, &approval.ApprovalStrID,
			&approval.EvalID, &approval.EvalStrID,
			&approval.TriggerID, &approval.TriggerStrID,
			&approval.WorkflowResultID, &approval.WorkflowResultStrID,
			&approval.ApprovalState, &approval.RequestSummary, &approval.UpdatedAt,
			&approval.TriggerAction)
		if err != nil {
			log.Err(err).Msg("failed to scan trigger action approval")
			return nil, err
		}
		approvals = append(approvals, approval)
	}
	// Check for any error encountered during iteration
	if err = rows.Err(); err != nil {
		log.Err(err).Msg("error encountered during rows iteration")
		return nil, err
	}
	return approvals, nil
}

type ApprovalApiReqResp struct {
	RetrievalItem                         RetrievalItem                         `json:"retrievalItem"`
	TriggerActionsApproval                TriggerActionsApproval                `json:"triggerActionsApproval"`
	AIWorkflowTriggerResultApiReqResponse AIWorkflowTriggerResultApiReqResponse `json:"aiWorkflowTriggerResultApiReqResponse"`
}

func SelectTriggerActionApprovalWithReqResponses(ctx context.Context, ou org_users.OrgUser, state string, approvalID, workflowResultID int) ([]ApprovalApiReqResp, error) {
	if approvalID <= 0 {
		return nil, nil
	}
	var approvals []ApprovalApiReqResp
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
        SELECT a.approval_id, a.approval_id::text,
               a.eval_id, a.eval_id::text,
               a.trigger_id, a.trigger_id::text,
               a.workflow_result_id, a.workflow_result_id::text, 
               a.approval_state, a.request_summary, a.updated_at,
               t.trigger_action,
               r.response_id,
               rl.retrieval_id, rl.retrieval_id::text AS retrieval_id_str, rl.retrieval_name, rl.retrieval_group, rl.retrieval_platform, rl.instructions,
               r.req_payload, r.resp_payload
        FROM public.ai_trigger_actions_approvals a
        JOIN public.ai_trigger_actions t ON a.trigger_id = t.trigger_id
		JOIN public.ai_trigger_actions_api_reqs_responses r ON r.approval_id = a.approval_id
        JOIN public.ai_retrieval_library rl ON rl.retrieval_id = r.retrieval_id
        WHERE t.org_id = $1 AND a.approval_state = $2 AND a.approval_id = $3 AND a.workflow_result_id = $4
        ORDER BY a.approval_id DESC;`

	// Executing the query
	rows, err := apps.Pg.Query(ctx, q.RawQuery, ou.OrgID, state, approvalID, workflowResultID)
	if err != nil {
		log.Err(err).Msg("failed to execute query for trigger action approvals")
		return nil, err
	}
	defer rows.Close()

	// Iterating through the query results
	for rows.Next() {
		var approval TriggerActionsApproval
		var resp AIWorkflowTriggerResultApiReqResponse
		var retrieval RetrievalItem
		err = rows.Scan(&approval.ApprovalID, &approval.ApprovalStrID,
			&approval.EvalID, &approval.EvalStrID,
			&approval.TriggerID, &approval.TriggerStrID,
			&approval.WorkflowResultID, &approval.WorkflowResultStrID,
			&approval.ApprovalState, &approval.RequestSummary, &approval.UpdatedAt,
			&approval.TriggerAction,
			&resp.ResponseID,
			&resp.RetrievalID, &retrieval.RetrievalStrID, &retrieval.RetrievalName,
			&retrieval.RetrievalGroup, &retrieval.RetrievalPlatform, &retrieval.Instructions,
			&resp.ReqPayloads, &resp.RespPayloads,
		)
		if err != nil {
			log.Err(err).Msg("failed to scan trigger action approval")
			return nil, err
		}

		retrieval.RetrievalID = aws.Int(resp.RetrievalID)
		b := retrieval.Instructions
		if b != nil {
			err = json.Unmarshal(b, &retrieval.RetrievalItemInstruction)
			if err != nil {
				log.Err(err).Msg("failed to unmarshal retrieval instructions")
				return nil, err
			}
			retrieval.Instructions = nil
		}
		resp.ApprovalID = approval.ApprovalID
		resp.TriggerID = approval.TriggerID
		tmp := ApprovalApiReqResp{
			TriggerActionsApproval:                approval,
			AIWorkflowTriggerResultApiReqResponse: resp,
			RetrievalItem:                         retrieval,
		}
		approvals = append(approvals, tmp)
	}
	// Check for any error encountered during iteration
	if err = rows.Err(); err != nil {
		log.Err(err).Msg("error encountered during rows iteration")
		return nil, err
	}
	return approvals, nil
}

func UpdateTriggerActionApproval(ctx context.Context, ou org_users.OrgUser, approval TriggerActionsApproval) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
        UPDATE public.ai_trigger_actions_approvals
		SET approval_state = $1, request_summary = $2
		WHERE approval_id = $3;`

	if approval.ApprovalState == "" {
		approval.ApprovalState = "pending"
	}
	if approval.RequestSummary == "" {
		approval.RequestSummary = "Requesting approval for trigger action"
	}
	// Executing the query
	ro, err := apps.Pg.Exec(ctx, q.RawQuery, approval.ApprovalState, approval.RequestSummary, approval.ApprovalID)
	if err != nil {
		log.Err(err).Interface("ro", ro.RowsAffected()).Msg("failed to insert or update trigger action approval")
		return err
	}
	return nil
}

package artemis_orchestrations

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type TriggerAction struct {
	TriggerID               int                      `db:"trigger_id" json:"triggerID"`
	OrgID                   int                      `db:"org_id" json:"orgID"`
	UserID                  int                      `db:"user_id" json:"userID"`
	TriggerName             string                   `db:"trigger_name" json:"triggerName"`
	TriggerGroup            string                   `db:"trigger_group" json:"triggerGroup"`
	EvalTriggerActions      []EvalTriggerActions     `db:"eval_trigger_actions" json:"evalTriggerActions"`
	TriggerActionsApprovals []TriggerActionsApproval `json:"aiTriggerActionsApproval,omitempty"`
}

type TriggerActionsApproval struct {
	ApprovalID       int       `db:"approval_id" json:"approvalID"`
	EvalID           int       `db:"eval_id" json:"evalID"`
	TriggerID        int       `db:"trigger_id" json:"triggerID"`
	WorkflowResultID int       `db:"workflow_result_id" json:"workflowResultID"`
	ApprovalState    string    `db:"approval_state" json:"approvalState"`
	RequestSummary   string    `db:"request_summary" json:"requestSummary"`
	UpdatedAt        time.Time `db:"updated_at" json:"updatedAt"`
}

type EvalTriggerActions struct {
	EvalID               int    `db:"eval_id" json:"evalID"`
	TriggerID            int    `db:"trigger_id" json:"triggerID"`
	EvalTriggerState     string `db:"eval_trigger_state" json:"evalTriggerState"`
	EvalResultsTriggerOn string `db:"eval_results_trigger_on" json:"evalResultsTriggerOn"`
}

func SelectTriggerActionApprovals(ctx context.Context, ou org_users.OrgUser, state string) ([]TriggerActionsApproval, error) {
	var approvals []TriggerActionsApproval

	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
        SELECT a.approval_id, a.eval_id, a.trigger_id, a.workflow_result_id, a.approval_state, a.request_summary, a.updated_at
        FROM public.ai_trigger_actions_approval a
        JOIN public.ai_trigger_actions t ON a.trigger_id = t.trigger_id
        WHERE t.org_id = $1 AND a.approval_state = $2
        ORDER BY a.approval_id DESC;`

	// Executing the query
	rows, err := apps.Pg.Query(ctx, q.RawQuery, ou.OrgID, state)
	if err != nil {
		log.Err(err).Msg("failed to execute query for trigger action approvals")
		return nil, err
	}
	defer rows.Close()

	// Iterating through the query results
	for rows.Next() {
		var approval TriggerActionsApproval
		err = rows.Scan(&approval.ApprovalID, &approval.EvalID, &approval.TriggerID, &approval.WorkflowResultID, &approval.ApprovalState, &approval.RequestSummary, &approval.UpdatedAt)
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

func CreateOrUpdateTriggerActionApproval(ctx context.Context, approval TriggerActionsApproval) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
        INSERT INTO public.ai_trigger_actions_approval (eval_id, trigger_id, workflow_result_id, approval_state, request_summary)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (eval_id, trigger_id, workflow_result_id)
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

func CreateOrUpdateAction(ctx context.Context, ou org_users.OrgUser, trigger *TriggerAction) error {
	if trigger == nil {
		return errors.New("trigger cannot be nil")
	}

	// Start a transaction
	tx, err := apps.Pg.Begin(ctx)
	if err != nil {
		return err
	}

	// Defer a rollback in case of failure
	defer tx.Rollback(ctx)

	// Insert or update the ai_trigger_actions
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
        INSERT INTO public.ai_trigger_actions (org_id, user_id, trigger_name, trigger_group)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (org_id, trigger_name) 
        DO UPDATE SET 
            user_id = EXCLUDED.user_id,
            trigger_group = EXCLUDED.trigger_group
        RETURNING trigger_id;`

	err = tx.QueryRow(ctx, q.RawQuery, ou.OrgID, ou.UserID, trigger.TriggerName, trigger.TriggerGroup).Scan(&trigger.TriggerID)
	if err != nil {
		log.Err(err).Msg("failed to insert ai trigger action")
		return err
	}
	for _, eta := range trigger.EvalTriggerActions {
		q.RawQuery = `
            INSERT INTO public.ai_trigger_actions_evals(eval_id, trigger_id, eval_trigger_state, eval_results_trigger_on)
            VALUES ($1, $2, $3, $4)
         	ON CONFLICT (eval_id, trigger_id)
    		DO UPDATE SET
				eval_trigger_state = EXCLUDED.eval_trigger_state,
				eval_results_trigger_on = EXCLUDED.eval_results_trigger_on;` // Adjust as needed

		_, err = tx.Exec(ctx, q.RawQuery, eta.EvalID, trigger.TriggerID, eta.EvalTriggerState, eta.EvalResultsTriggerOn)
		if err != nil {
			log.Err(err).Msg("failed to insert eval trigger action")
			return err
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Err(err).Msg("failed to commit transaction")
		return err
	}
	return nil
}

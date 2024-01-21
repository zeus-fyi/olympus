package artemis_orchestrations

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type TriggerAction struct {
	TriggerID                int                      `db:"trigger_id" json:"triggerID,omitempty"`
	OrgID                    int                      `db:"org_id" json:"orgID,omitempty"`
	UserID                   int                      `db:"user_id" json:"userID,omitempty"`
	TriggerName              string                   `db:"trigger_name" json:"triggerName"`
	TriggerGroup             string                   `db:"trigger_group" json:"triggerGroup"`
	TriggerAction            string                   `db:"trigger_action" json:"triggerAction"`
	TriggerPlatformReference TriggerPlatformReference `db:"platforms_reference" json:"platformReference,omitempty"`
	EvalTriggerAction        EvalTriggerActions       `db:"eval_trigger_actions" json:"evalTriggerAction,omitempty"`
	EvalTriggerActions       []EvalTriggerActions     `db:"eval_trigger_actions" json:"evalTriggerActions,omitempty"`
	TriggerActionsApprovals  []TriggerActionsApproval `json:"triggerActionsApprovals,omitempty"`
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

type TriggerPlatformReference struct {
	PlatformReferenceID   int    `db:"platforms_reference_id" json:"platformReferenceID"`
	PlatformReferenceName string `db:"platforms_reference_name" json:"platformReferenceName"`
}

type EvalTriggerActions struct {
	EvalID               int    `db:"eval_id" json:"evalID,omitempty"`
	TriggerID            int    `db:"trigger_id" json:"triggerID,omitempty"`
	EvalTriggerState     string `db:"eval_trigger_state" json:"evalTriggerState"` // eg. info, filter, etc
	EvalResultsTriggerOn string `db:"eval_results_trigger_on" json:"evalResultsTriggerOn"`
}

func SelectTriggerActionsByOrgAndOptParams(ctx context.Context, ou org_users.OrgUser, evalID int) ([]TriggerAction, error) {
	var triggerActions []TriggerAction
	var currentTriggerID int
	triggerActionMap := make(map[int]*TriggerAction)

	q := sql_query_templates.QueryParams{}
	params := []interface{}{ou.OrgID}

	additionalQuery := ""
	if evalID != 0 {
		additionalQuery = "AND eval_id = $2"
		params = append(params, evalID)
	}
	// Updated query to include TriggerActionsApproval
	q.RawQuery = `
			WITH TriggerActions AS (
				SELECT ta.trigger_id, ta.trigger_name, ta.trigger_group, ta.trigger_action,
					   COALESCE(tae.eval_id, 0) as eval_id, taee.eval_trigger_state, taee.eval_results_trigger_on
				FROM public.ai_trigger_actions ta
				LEFT JOIN public.ai_trigger_actions_evals tae ON ta.trigger_id = tae.trigger_id
				LEFT JOIN public.ai_trigger_eval taee ON ta.trigger_id = taee.trigger_id
				WHERE ta.org_id = $1` + additionalQuery + `
			)
			SELECT ta.trigger_id, ta.trigger_name, ta.trigger_group, ta.trigger_action,
				   COALESCE(ta.eval_id,0), ta.eval_trigger_state, ta.eval_results_trigger_on,
				   COALESCE(JSON_AGG(
					   JSON_BUILD_OBJECT(
						   'workflowResultID', ataa.workflow_result_id,
 					       'evalID', ataa.eval_id,
						   'triggerID', ataa.trigger_id,
						   'approvalID', ataa.approval_id,
						   'approvalState', ataa.approval_state,
						   'requestSummary', ataa.request_summary,
						   'updatedAt', ataa.updated_at
					   ) ORDER BY CASE WHEN ataa.approval_state = 'pending' THEN 0 ELSE 1 END, ataa.approval_id DESC
				   ) FILTER (WHERE ataa.approval_id IS NOT NULL), '[]') AS approvals
			FROM TriggerActions ta
			LEFT JOIN public.ai_trigger_actions_approval ataa ON ta.trigger_id = ataa.trigger_id
			GROUP BY ta.trigger_id, ta.trigger_name, ta.trigger_group, ta.trigger_action,
					 ta.eval_id, ta.eval_trigger_state, ta.eval_results_trigger_on;`

	// Executing the query
	rows, err := apps.Pg.Query(ctx, q.RawQuery, params...)
	if err != nil {
		log.Err(err).Msg("failed to execute query for trigger actions")
		return nil, err
	}
	defer rows.Close()

	// Iterating through the query results
	for rows.Next() {
		var triggerName, triggerGroup, triggerEnv string
		var approvalsJSON *string
		var currentEvalTriggerActions EvalTriggerActions
		var currentTriggerActionsApprovals []TriggerActionsApproval

		err = rows.Scan(
			&currentTriggerID,
			&triggerName,
			&triggerGroup,
			&triggerEnv,
			&currentEvalTriggerActions.EvalID,
			&currentEvalTriggerActions.EvalTriggerState,
			&currentEvalTriggerActions.EvalResultsTriggerOn,
			&approvalsJSON,
		)
		if err != nil {
			log.Err(err).Msg("failed to scan trigger action")
			return nil, err
		}

		if approvalsJSON != nil {
			// Parse the JSON string into TriggerActionsApproval slice
			err = json.Unmarshal([]byte(*approvalsJSON), &currentTriggerActionsApprovals)
			if err != nil {
				log.Err(err).Msg("failed to unmarshal trigger actions approvals")
				return nil, err
			}
		}
		currentEvalTriggerActions.TriggerID = currentTriggerID
		if _, exists := triggerActionMap[currentTriggerID]; exists {
			triggerActionMap[currentTriggerID].TriggerActionsApprovals = append(triggerActionMap[currentTriggerID].TriggerActionsApprovals, currentTriggerActionsApprovals...)
		} else {
			// TODO fix the slice conversion hack
			triggerActionMap[currentTriggerID] = &TriggerAction{
				TriggerID:               currentTriggerID,
				TriggerName:             triggerName,
				TriggerGroup:            triggerGroup,
				TriggerAction:           triggerEnv,
				EvalTriggerAction:       currentEvalTriggerActions,
				EvalTriggerActions:      []EvalTriggerActions{currentEvalTriggerActions},
				TriggerActionsApprovals: currentTriggerActionsApprovals,
			}
		}
	}

	// Check for any error encountered during iteration
	if err = rows.Err(); err != nil {
		log.Err(err).Msg("error encountered during rows iteration")
		return nil, err
	}

	// Convert the map to a slice
	for tid, ta := range triggerActionMap {
		ta.TriggerID = tid
		triggerActions = append(triggerActions, *ta)
	}

	return triggerActions, nil
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

// TODO, need to add orgID to the query

// request summary?

func CreateOrUpdateTriggerActionApproval(ctx context.Context, ou org_users.OrgUser, approval *TriggerActionsApproval) error {
	if approval == nil {
		return errors.New("approval cannot be nil")
	}
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
        INSERT INTO ai_trigger_actions_approval(eval_id, trigger_id, workflow_result_id, approval_state, request_summary)
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

func CreateOrUpdateTriggerAction(ctx context.Context, ou org_users.OrgUser, trigger *TriggerAction) error {
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
        INSERT INTO public.ai_trigger_actions (org_id, user_id, trigger_name, trigger_group, trigger_action)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (org_id, trigger_name) 
        DO UPDATE SET 
            user_id = EXCLUDED.user_id,
            trigger_action = EXCLUDED.trigger_action,
            trigger_group = EXCLUDED.trigger_group
        RETURNING trigger_id;`

	err = tx.QueryRow(ctx, q.RawQuery, ou.OrgID, ou.UserID, trigger.TriggerName, trigger.TriggerGroup, trigger.TriggerAction).Scan(&trigger.TriggerID)
	if err != nil {
		log.Err(err).Msg("failed to insert ai trigger action")
		return err
	}
	for _, eta := range trigger.EvalTriggerActions {
		q.RawQuery = `
            INSERT INTO public.ai_trigger_eval(trigger_id, eval_trigger_state, eval_results_trigger_on)
            VALUES ($1, $2, $3)
         	ON CONFLICT (trigger_id)
         	DO UPDATE SET
				eval_trigger_state = EXCLUDED.eval_trigger_state,
				eval_results_trigger_on = EXCLUDED.eval_results_trigger_on;`
		_, err = tx.Exec(ctx, q.RawQuery, trigger.TriggerID, eta.EvalTriggerState, eta.EvalResultsTriggerOn)
		if err != nil {
			log.Err(err).Msg("failed to insert eval trigger action")
			return err
		}
		if eta.EvalID != 0 {
			q.RawQuery = `
            INSERT INTO public.ai_trigger_actions_evals(eval_id, trigger_id)
            VALUES ($1, $2)
         	ON CONFLICT (eval_id, trigger_id)
    		DO NOTHING;`
			_, err = tx.Exec(ctx, q.RawQuery, eta.EvalID, trigger.TriggerID)
			if err != nil {
				log.Err(err).Msg("failed to insert eval trigger action")
				return err
			}
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Err(err).Msg("failed to commit transaction")
		return err
	}
	return nil
}

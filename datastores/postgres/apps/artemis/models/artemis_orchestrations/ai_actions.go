package artemis_orchestrations

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgtype"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type TriggerAction struct {
	TriggerStrID            string                   `json:"triggerStrID,omitempty"`
	TriggerID               int                      `db:"trigger_id" json:"triggerID,omitempty"`
	OrgID                   int                      `db:"org_id" json:"orgID,omitempty"`
	UserID                  int                      `db:"user_id" json:"userID,omitempty"`
	TriggerName             string                   `db:"trigger_name" json:"triggerName"`
	TriggerGroup            string                   `db:"trigger_group" json:"triggerGroup"`
	TriggerAction           string                   `db:"trigger_action" json:"triggerAction"`
	TriggerRetrievals       []RetrievalItem          `json:"triggerRetrievals,omitempty"`
	EvalTriggerAction       EvalTriggerActions       `db:"eval_trigger_action" json:"evalTriggerAction,omitempty"`
	EvalTriggerActions      []EvalTriggerActions     `db:"eval_trigger_actions" json:"evalTriggerActions,omitempty"`
	TriggerActionsApprovals []TriggerActionsApproval `json:"triggerActionsApprovals,omitempty"`
}

type TriggerActionsApproval struct {
	ApprovalStrID       string    `json:"approvalStrID,omitempty"`
	ApprovalID          int       `db:"approval_id" json:"approvalID"`
	EvalID              int       `db:"eval_id" json:"evalID"`
	TriggerID           int       `db:"trigger_id" json:"triggerID"`
	TriggerStrID        string    `db:"trigger_id" json:"triggerStrID"`
	WorkflowResultID    int       `db:"workflow_result_id" json:"workflowResultID"`
	WorkflowResultStrID string    `json:"workflowResultStrID"`
	ApprovalState       string    `db:"approval_state" json:"approvalState"`
	RequestSummary      string    `db:"request_summary" json:"requestSummary"`
	UpdatedAt           time.Time `db:"updated_at" json:"updatedAt"`
}

type EvalTriggerActions struct {
	EvalID               int    `db:"eval_id" json:"evalID,omitempty"`
	EvalStrID            string `db:"eval_str_id" json:"evalStrID,omitempty"`
	TriggerID            int    `db:"trigger_id" json:"triggerID,omitempty"`
	TriggerStrID         string `db:"trigger_str_id" json:"triggerStrID,omitempty"`
	EvalTriggerState     string `db:"eval_trigger_state" json:"evalTriggerState"` // eg. info, filter, etc
	EvalResultsTriggerOn string `db:"eval_results_trigger_on" json:"evalResultsTriggerOn"`
}

func SelectTriggerActionsByOrgAndOptParams(ctx context.Context, ou org_users.OrgUser, evalID int) ([]TriggerAction, error) {
	var triggerActions []TriggerAction
	q := sql_query_templates.QueryParams{}
	params := []interface{}{ou.OrgID}

	additionalQuery := ""
	if evalID != 0 {
		additionalQuery = "AND eval_id = $2"
		params = append(params, evalID)
	}
	// Updated query to include TriggerActionsApproval
	q.RawQuery = `
			WITH cte_trigger_acts AS (
				SELECT ta.trigger_id, ta.trigger_name, ta.trigger_group, ta.trigger_action
				FROM public.ai_trigger_actions ta
				WHERE ta.org_id = $1 ` + additionalQuery + `
			),
			cte_trigger_action_evals AS (
				SELECT ta.trigger_id, COALESCE(tae.eval_id, 0) as eval_id, taee.eval_trigger_state, taee.eval_results_trigger_on
				FROM cte_trigger_acts ta
				LEFT JOIN public.ai_trigger_actions_evals tae ON ta.trigger_id = tae.trigger_id
				LEFT JOIN public.ai_trigger_eval taee ON ta.trigger_id = taee.trigger_id
				WHERE tae.eval_id > 0
			),
			cte_trigger_api_rets AS (
				SELECT 
					tapi.trigger_id,
					COALESCE(
						JSONB_AGG(
							JSONB_BUILD_OBJECT(
								'retrievalID', tapi.retrieval_id,
								'retrievalStrID', tapi.retrieval_id::text,
								'retrievalName', re.retrieval_name,
								'retrievalGroup', re.retrieval_group,
								'retrievalItemInstruction', JSONB_BUILD_OBJECT(
									'retrievalPlatform', '',
									'retrievalPrompt', '', 
									'retrievalPlatformGroups', '', 
									'retrievalKeywords', '', 
									'retrievalNegativeKeywords', '', 
									'retrievalUsernames', '', 
									'discordFilters', JSONB_BUILD_OBJECT(
										'categoryTopic', '',  
										'categoryName', '',
										'category', ''
									),
									'webFilters', JSONB_BUILD_OBJECT(
										'routingGroup', '',  
										'lbStrategy', '',
										'maxRetries', 0,
										'backoffCoefficient', 2,
										'endpointRoutePath', '',
										'endpointREST', ''
									),
									'instructions', COALESCE(re.instructions, '{}'::jsonb)
								)
							)
						) FILTER (WHERE tapi.retrieval_id IS NOT NULL), '[]'::jsonb
					) AS retrievals
				FROM cte_trigger_acts ta
				JOIN public.ai_trigger_actions_api tapi ON tapi.trigger_id = ta.trigger_id
				LEFT JOIN public.ai_retrieval_library re ON tapi.retrieval_id = re.retrieval_id
				GROUP BY tapi.trigger_id
			),
			cte_trigger_action_approvals AS (
				SELECT 
					ta.trigger_id,
					ataa.workflow_result_id,
					ataa.eval_id,
					ataa.trigger_id AS approval_trigger_id,
					ataa.approval_id,
					ataa.approval_state,
					ataa.request_summary,
					ataa.updated_at
				FROM cte_trigger_action_evals ta
				JOIN public.ai_trigger_actions_approval ataa ON ta.trigger_id = ataa.trigger_id
					GROUP BY ta.trigger_id,
					ataa.workflow_result_id,
					ataa.eval_id,
					ataa.trigger_id,
					ataa.approval_id,
					ataa.approval_state,
					ataa.request_summary,
					ataa.updated_at
			), cte_trigger_action_approvals_agg AS (
				SELECT 
					trigger_id,
					COALESCE(JSONB_AGG(
						JSONB_BUILD_OBJECT(
							'workflowResultStrID', workflow_result_id::text,
							'workflowResultID', workflow_result_id,
							'evalID', eval_id,
							'evalStrID', eval_id::text,
							'triggerID', trigger_id,
							'triggerStrID', trigger_id::text,
							'approvalID', approval_id,
							'approvalStrID', approval_id::text,
							'approvalState', approval_state,
							'requestSummary', request_summary,
							'updatedAt', updated_at
						) ORDER BY CASE WHEN approval_state = 'pending' THEN 0 ELSE 1 END, approval_id DESC
					) FILTER (WHERE approval_id IS NOT NULL), '[]'::jsonb) AS approvals
				FROM cte_trigger_action_approvals
				GROUP BY trigger_id, approval_trigger_id
			), cte_agg_eval_trgs AS (
				SELECT ce.trigger_id, 
				 COALESCE(JSONB_AGG(
						JSONB_BUILD_OBJECT(
							'evalID', ce.eval_id,
							'triggerID', ce.trigger_id,
							'evalStrID', ce.eval_id::text,
							'triggerStrID', ce.trigger_id::text,
							'evalTriggerState', ce.eval_trigger_state,
							'evalResultsTriggerOn', ce.eval_results_trigger_on
						) 
					), '[]'::jsonb) AS eval_triggers
			  FROM cte_trigger_action_evals ce
				GROUP BY trigger_id
			)
			SELECT 
				acts.trigger_id,
				acts.trigger_id::text,
				acts.trigger_name,
				acts.trigger_group,
				acts.trigger_action,
				eval_triggers,
				api_rets.retrievals,
				agg.approvals
			FROM cte_trigger_acts acts
			LEFT JOIN cte_agg_eval_trgs evals ON acts.trigger_id = evals.trigger_id
			LEFT JOIN cte_trigger_api_rets api_rets ON acts.trigger_id = api_rets.trigger_id
			LEFT JOIN cte_trigger_action_approvals_agg agg ON acts.trigger_id = agg.trigger_id`

	// Executing the query
	rows, err := apps.Pg.Query(ctx, q.RawQuery, params...)
	if err != nil {
		log.Err(err).Msg("failed to execute query for trigger actions")
		return nil, err
	}
	defer rows.Close()

	// Iterating through the query results
	for rows.Next() {
		var triggerAction TriggerAction
		err = rows.Scan(
			&triggerAction.TriggerID,
			&triggerAction.TriggerStrID,
			&triggerAction.TriggerName,
			&triggerAction.TriggerGroup,
			&triggerAction.TriggerAction,
			&triggerAction.EvalTriggerActions,
			&triggerAction.TriggerRetrievals,
			&triggerAction.TriggerActionsApprovals,
		)
		if err != nil {
			log.Err(err).Msg("failed to scan trigger action")
			return nil, err
		}
		for ri, _ := range triggerAction.TriggerRetrievals {
			b, berr := triggerAction.TriggerRetrievals[ri].Instructions.MarshalJSON()
			if berr != nil {
				log.Err(berr).Msg("unmarshal error")
				return nil, berr
			}
			if b != nil {
				err = json.Unmarshal(b, &triggerAction.TriggerRetrievals[ri].RetrievalItemInstruction)
				if err != nil {
					log.Err(err).Msg("failed to unmarshal retrieval instructions")
					return nil, err
				}
				triggerAction.TriggerRetrievals[ri].RetrievalItemInstruction.Instructions = pgtype.JSONB{
					Bytes:  nil,
					Status: pgtype.Null,
				}
			}
		}
		triggerActions = append(triggerActions, triggerAction)
	}
	if err = rows.Err(); err != nil {
		log.Err(err).Msg("error encountered during rows iteration")
		return nil, err
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

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
	TriggerStrID               string                   `json:"triggerStrID,omitempty"`
	TriggerID                  int                      `db:"trigger_id" json:"triggerID,omitempty"`
	OrgID                      int                      `db:"org_id" json:"orgID,omitempty"`
	UserID                     int                      `db:"user_id" json:"userID,omitempty"`
	TriggerName                string                   `db:"trigger_name" json:"triggerName"`
	TriggerGroup               string                   `db:"trigger_group" json:"triggerGroup"`
	TriggerAction              string                   `db:"trigger_action" json:"triggerAction"`
	TriggerExpirationDuration  float64                  `json:"triggerExpirationDuration,omitempty"`
	TriggerExpirationTimeUnit  string                   `json:"triggerExpirationTimeUnit,omitempty"`
	TriggerExpiresAfterSeconds int                      `db:"expires_at" json:"triggerExpiresAfter,omitempty"`
	TriggerRetrievals          []RetrievalItem          `json:"triggerRetrievals,omitempty"`
	EvalTriggerAction          EvalTriggerActions       `db:"eval_trigger_action" json:"evalTriggerAction,omitempty"`
	EvalTriggerActions         []EvalTriggerActions     `db:"eval_trigger_actions" json:"evalTriggerActions,omitempty"`
	TriggerActionsApprovals    []TriggerActionsApproval `json:"triggerActionsApprovals,omitempty"`
}

type TriggerActionsApproval struct {
	TriggerAction       string    `db:"trigger_action" json:"triggerAction"`
	ApprovalStrID       string    `json:"approvalStrID,omitempty"`
	ApprovalID          int       `db:"approval_id" json:"approvalID"`
	EvalID              int       `db:"eval_id" json:"evalID"`
	EvalStrID           string    `json:"evalStrID"`
	TriggerID           int       `db:"trigger_id" json:"triggerID"`
	TriggerStrID        string    `db:"trigger_id" json:"triggerStrID"`
	WorkflowResultID    int       `db:"workflow_result_id" json:"workflowResultID"`
	WorkflowResultStrID string    `json:"workflowResultStrID"`
	ApprovalState       string    `db:"approval_state" json:"approvalState"`
	RequestSummary      string    `db:"request_summary" json:"requestSummary"`
	UpdatedAt           time.Time `db:"updated_at" json:"updatedAt"`
	//ExpiresAt           *time.Time `db:"expires_at" json:"expiresAt"`

	Requests  json.RawMessage `json:"requests,omitempty"`
	Responses json.RawMessage `json:"responses,omitempty"`
}

type EvalTriggerActions struct {
	EvalID               int    `db:"eval_id" json:"evalID,omitempty"`
	EvalStrID            string `db:"eval_str_id" json:"evalStrID,omitempty"`
	TriggerID            int    `db:"trigger_id" json:"triggerID,omitempty"`
	TriggerStrID         string `db:"trigger_str_id" json:"triggerStrID,omitempty"`
	EvalTriggerState     string `db:"eval_trigger_state" json:"evalTriggerState"` // eg. info, filter, etc
	EvalResultsTriggerOn string `db:"eval_results_trigger_on" json:"evalResultsTriggerOn"`
}

type TriggersWorkflowQueryParams struct {
	Ou                 org_users.OrgUser `json:"ou,omitempty"`
	EvalID             int               `json:"evalID,omitempty"`
	TaskID             int               `json:"taskID,omitempty"`
	WorkflowTemplateID int               `json:"workflowTemplateID,omitempty"`
}

func (tq *TriggersWorkflowQueryParams) ValidateEvalTaskQp() bool {
	if tq.Ou.OrgID == 0 || tq.EvalID == 0 || tq.TaskID == 0 || tq.WorkflowTemplateID == 0 {
		return false
	}
	return true
}

func SelectTriggerActionsByOrgAndOptParams(ctx context.Context, tq TriggersWorkflowQueryParams) ([]TriggerAction, error) {
	if tq.Ou.OrgID == 0 {
		return nil, errors.New("orgID cannot be 0")
	}
	var triggerActions []TriggerAction
	q := sql_query_templates.QueryParams{}
	params := []interface{}{tq.Ou.OrgID}

	q1 := `	WITH cte_trigger_acts AS (
				SELECT ta.trigger_id, ta.trigger_name, ta.trigger_group, ta.trigger_action, ta.expires_after_seconds
				FROM public.ai_trigger_actions ta
				WHERE ta.org_id = $1
			),`

	if tq.EvalID != 0 && tq.TaskID != 0 && tq.WorkflowTemplateID != 0 {
		q1 = `	WITH cte_trigger_acts AS (
				SELECT ta.trigger_id, ta.trigger_name, ta.trigger_group, ta.trigger_action, ta.expires_after_seconds
				FROM public.ai_trigger_actions ta
				JOIN public.ai_trigger_actions_evals tae ON ta.trigger_id = tae.trigger_id
				JOIN public.ai_workflow_template_eval_task_relationships trrr ON trrr.eval_id = tae.eval_id
				WHERE ta.org_id = $1 AND tae.eval_id = $2 AND trrr.task_id = $3 AND trrr.workflow_template_id = $4
			),`
		params = append(params, tq.EvalID, tq.TaskID, tq.WorkflowTemplateID)
	}

	// Updated query to include TriggerActionsApproval
	q.RawQuery = q1 + `
			cte_trigger_action_evals AS (
				SELECT ta.trigger_id, ta.trigger_action, ta.expires_after_seconds, COALESCE(tae.eval_id, 0) as eval_id, taee.eval_trigger_state, taee.eval_results_trigger_on
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
					ta.trigger_action,
					ta.expires_after_seconds,
					ataa.workflow_result_id,
					ataa.eval_id,
					ataa.trigger_id AS approval_trigger_id,
					ataa.approval_id,
					ataa.approval_state,
					ataa.request_summary,
					ataa.updated_at
				FROM cte_trigger_action_evals ta
				JOIN public.ai_trigger_actions_approvals ataa ON ta.trigger_id = ataa.trigger_id
			   	WHERE ataa.expires_at > NOW() OR ataa.expires_at IS NULL OR ataa.approval_state = 'finished' OR ataa.approval_state = 'rejected' OR ataa.approval_state = 'approved'
				GROUP BY ta.trigger_id, ta.trigger_action, ta.expires_after_seconds,
					ataa.workflow_result_id,
					ataa.eval_id,
					ataa.trigger_id,
					ataa.approval_id,
					ataa.approval_state,
					ataa.request_summary,
					ataa.updated_at
		), cte_trigger_action_approvals_agg AS (
			SELECT 
				a.trigger_id,
				a.approval_id,
				COALESCE(JSONB_AGG(
					JSONB_BUILD_OBJECT(
						'triggerAction', trigger_action,
						'workflowResultStrID', workflow_result_id::text,
						'workflowResultID', workflow_result_id,
						'evalID', eval_id,
						'evalStrID', eval_id::text,
						'triggerID', a.trigger_id,
						'triggerStrID', a.trigger_id::text,
						'approvalID', a.approval_id,
						'approvalStrID', a.approval_id::text,
						'approvalState', approval_state,
						'requestSummary', request_summary,
						'updatedAt', updated_at
					) ORDER BY CASE WHEN approval_state = 'pending' THEN 0 ELSE 1 END, a.approval_id DESC
				) FILTER (WHERE a.approval_id IS NOT NULL), '[]'::jsonb) AS approvals
			FROM cte_trigger_action_approvals a
			GROUP BY a.trigger_id, approval_id
		), cte_11 AS (
			SELECT 
				a.approval_id, 
				a.trigger_id, 
				a.approvals,
						JSONB_BUILD_OBJECT('requests', 
						COALESCE(JSONB_AGG(COALESCE(r.req_payload, '{}'::jsonb)), '[]'::jsonb)) AS requests,
										JSONB_BUILD_OBJECT('responses', 
				COALESCE(JSONB_AGG(COALESCE(r.resp_payload, '{}'::jsonb)), '[]'::jsonb))AS responses
			FROM cte_trigger_action_approvals_agg a
			LEFT JOIN public.ai_trigger_actions_api_reqs_responses r ON r.approval_id = a.approval_id
			GROUP BY a.approval_id, a.trigger_id, a.approvals
		), cte_1 AS (
			SELECT 
				agg.trigger_id,
				JSONB_AGG(
					JSONB_BUILD_OBJECT(
						'triggerAction', x.approval->'triggerAction',
						'workflowResultStrID', x.approval->'workflowResultStrID',
						'workflowResultID', x.approval->'workflowResultID',
						'evalID', x.approval->'evalID',
						'evalStrID', x.approval->'evalStrID',
						'triggerID', x.approval->'triggerID',
						'triggerStrID', x.approval->'triggerStrID',
						'approvalID', x.approval->'approvalID',
						'approvalStrID', x.approval->'approvalStrID',
						'approvalState', x.approval->'approvalState',
						'requestSummary', x.approval->'requestSummary',
						'updatedAt', x.approval->'updatedAt',
						'requests', c11.requests->'requests',
						'responses', c11.responses->'responses'
					)
				) AS approvals
			FROM cte_trigger_action_approvals_agg agg
			CROSS JOIN LATERAL JSONB_ARRAY_ELEMENTS(agg.approvals) AS x(approval)
			LEFT JOIN cte_11 c11 ON agg.approval_id = c11.approval_id AND agg.trigger_id = c11.trigger_id
			GROUP BY agg.trigger_id
			), cte_agg_eval_trgs AS (
				SELECT ce.trigger_id, 
				 COALESCE(JSONB_AGG(
						JSONB_BUILD_OBJECT(
							'evalID', ce.eval_id,
							'evalStrID', ce.eval_id::text,
							'triggerID', ce.trigger_id,
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
				acts.expires_after_seconds,
				eval_triggers,
				COALESCE(api_rets.retrievals, '[]'::jsonb),
				COALESCE(approvals, '[]'::jsonb)
			FROM cte_trigger_acts acts
			LEFT JOIN cte_agg_eval_trgs evals ON acts.trigger_id = evals.trigger_id
			LEFT JOIN cte_trigger_api_rets api_rets ON acts.trigger_id = api_rets.trigger_id
			LEFT JOIN cte_1 agg ON acts.trigger_id = agg.trigger_id`

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
			&triggerAction.TriggerExpiresAfterSeconds,
			&triggerAction.EvalTriggerActions,
			&triggerAction.TriggerRetrievals,
			&triggerAction.TriggerActionsApprovals,
		)
		if err != nil {
			log.Err(err).Msg("failed to scan trigger action")
			return nil, err
		}
		if len(triggerAction.EvalTriggerActions) > 0 && triggerAction.EvalTriggerActions[0].EvalID > 0 {
			triggerAction.EvalTriggerAction = triggerAction.EvalTriggerActions[0]
		}

		if triggerAction.TriggerExpiresAfterSeconds > 0 {
			expDuration, timeUnit := ConvertSecondsToLargestUnit(triggerAction.TriggerExpiresAfterSeconds)
			triggerAction.TriggerExpirationTimeUnit = timeUnit
			triggerAction.TriggerExpirationDuration = float64(expDuration)
		}

		for ri, _ := range triggerAction.TriggerRetrievals {
			b := triggerAction.TriggerRetrievals[ri].Instructions
			if b != nil {
				err = json.Unmarshal(b, &triggerAction.TriggerRetrievals[ri].RetrievalItemInstruction)
				if err != nil {
					log.Err(err).Msg("failed to unmarshal retrieval instructions")
					return nil, err
				}
			}
			triggerAction.TriggerRetrievals[ri].Instructions = nil
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
        SELECT a.approval_id, a.approval_id::text,
               a.eval_id,
               a.trigger_id, a.trigger_id::text,
               a.workflow_result_id, a.workflow_result_id::text, 
               a.approval_state, a.request_summary, a.updated_at
        FROM public.ai_trigger_actions_approvals a
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
		err = rows.Scan(&approval.ApprovalID, &approval.ApprovalStrID,
			&approval.EvalID,
			&approval.TriggerID, &approval.TriggerStrID,
			&approval.WorkflowResultID, &approval.WorkflowResultStrID,
			&approval.ApprovalState, &approval.RequestSummary, &approval.UpdatedAt)
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

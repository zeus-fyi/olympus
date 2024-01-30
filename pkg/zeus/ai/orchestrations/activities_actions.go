package ai_platform_service_orchestrations

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
)

// SendTriggerActionRequestForApproval sends the action request to the user for human in-the-loop approval
func (z *ZeusAiPlatformActivities) SendTriggerActionRequestForApproval(ctx context.Context) error {
	return nil
}

func (z *ZeusAiPlatformActivities) LookupEvalTriggerConditions(ctx context.Context, tq artemis_orchestrations.TriggersWorkflowQueryParams) ([]artemis_orchestrations.TriggerAction, error) {
	ta, err := artemis_orchestrations.SelectTriggerActionsByOrgAndOptParams(ctx, tq)
	if err != nil {
		log.Err(err).Msg("LookupEvalTriggerConditions: failed to lookup trigger actions")
		return nil, err
	}
	return ta, nil
}

func (z *ZeusAiPlatformActivities) CreateOrUpdateTriggerActionToExec(ctx context.Context, ou org_users.OrgUser, act artemis_orchestrations.TriggerActionsApproval) error {
	err := artemis_orchestrations.CreateOrUpdateTriggerActionApproval(ctx, ou, act)
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg("CreateOrUpdateTriggerActionToExec: failed to create or update trigger action approval")
		return err
	}
	return nil
}

func (z *ZeusAiPlatformActivities) SelectTriggerActionToExec(ctx context.Context, ou org_users.OrgUser, approvalID int) ([]artemis_orchestrations.TriggerActionsApproval, error) {
	ap, err := artemis_orchestrations.SelectTriggerActionApproval(ctx, ou, "pending", approvalID)
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg("CreateOrUpdateTriggerActionToExec: failed to create or update trigger action approval")
		return nil, err
	}
	return ap, nil
}

const (
	infoState     = "info"
	pendingStatus = "pending"
)

func (z *ZeusAiPlatformActivities) CheckEvalTriggerCondition(ctx context.Context, act *artemis_orchestrations.TriggerAction, emr *artemis_orchestrations.EvalMetricsResults) ([]artemis_orchestrations.TriggerActionsApproval, error) {
	if act == nil || emr == nil || emr.EvalMetricsResults == nil {
		return nil, nil
	}
	var taps []artemis_orchestrations.TriggerActionsApproval
	m := make(map[string][]bool)

	for _, er := range emr.EvalMetricsResults {
		if er.EvalState == "" {
			continue
		}
		if _, ok := m[er.EvalState]; !ok {
			m[er.EvalState] = []bool{}
		}
		if er.EvalMetricResult == nil || er.EvalMetricResult.EvalResultOutcomeBool == nil {
			continue
		}
		m[er.EvalState] = append(m[er.EvalState], *er.EvalMetricResult.EvalResultOutcomeBool)
	}

	ch := chronos.Chronos{}
	// gets the eval results by state, eg. info, trigger, etc.
	for _, tr := range act.EvalTriggerActions {
		results := m[tr.EvalTriggerState]
		if len(results) <= 0 {
			continue
		}
		// when the trigger on eval results condition is met, create a trigger action for approval
		if checkTriggerOnEvalResults(tr.EvalResultsTriggerOn, results) {
			tap := artemis_orchestrations.TriggerActionsApproval{
				ApprovalID:       ch.UnixTimeStampNow(),
				WorkflowResultID: emr.EvalContext.AIWorkflowAnalysisResult.WorkflowResultID,
				EvalID:           tr.EvalID,
				TriggerID:        tr.TriggerID,
				TriggerAction:    apiApproval,
				ApprovalState:    pendingStatus,
			}
			taps = append(taps, tap)
		}
	}
	return taps, nil
}

func (z *ZeusAiPlatformActivities) SaveTriggerResponseOutput(ctx context.Context, trrr artemis_orchestrations.AIWorkflowTriggerResultResponse) error {
	respID, err := artemis_orchestrations.InsertOrUpdateAIWorkflowTriggerResultResponse(ctx, trrr)
	if err != nil {
		log.Err(err).Interface("respID", respID).Interface("trrr", trrr).Msg("SaveTriggerResponseOutput: failed")
		return err
	}
	return nil
}

func (z *ZeusAiPlatformActivities) SaveTriggerApiRequestResp(ctx context.Context, trrr *artemis_orchestrations.AIWorkflowTriggerResultApiReqResponse) (*artemis_orchestrations.AIWorkflowTriggerResultApiReqResponse, error) {
	err := artemis_orchestrations.InsertOrUpdateAIWorkflowTriggerResultApiResponse(ctx, trrr)
	if err != nil {
		log.Err(err).Interface("responseID", trrr.ResponseID).Interface("trrr", trrr).Msg("SaveTriggerApiResponseOutput: failed")
		return trrr, err
	}
	return trrr, nil
}

func (z *ZeusAiPlatformActivities) CreateOrUpdateTriggerActionApprovalWithApiReq(ctx context.Context, ou org_users.OrgUser,
	approval artemis_orchestrations.TriggerActionsApproval, wtrr artemis_orchestrations.AIWorkflowTriggerResultApiReqResponse) error {
	err := artemis_orchestrations.CreateOrUpdateTriggerActionApprovalWithApiReq(ctx, ou, approval, wtrr)
	if err != nil {
		log.Err(err).Interface("ou", ou).Interface("approval", approval).Interface("wtrr", wtrr).Msg("CreateOrUpdateTriggerActionApprovalWithApiReq: failed")
		return err
	}
	return nil
}

func (z *ZeusAiPlatformActivities) SelectTriggerActionApiApprovalWithReqResponses(ctx context.Context, ou org_users.OrgUser, state string, approvalID, workflowResultID int) ([]artemis_orchestrations.ApprovalApiReqResp, error) {
	resp, err := artemis_orchestrations.SelectTriggerActionApprovalWithReqResponses(ctx, ou, state, approvalID, workflowResultID)
	if err != nil {
		log.Err(err).Interface("ou", ou).Interface("state", state).Interface("approvalID", approvalID).Msg("SelectTriggerActionApprovalWithReqResponses: failed")
		return nil, err
	}
	return resp, nil
}

func (z *ZeusAiPlatformActivities) UpdateTriggerActionApproval(ctx context.Context, ou org_users.OrgUser, approval artemis_orchestrations.TriggerActionsApproval) error {
	err := artemis_orchestrations.UpdateTriggerActionApproval(ctx, ou, approval)
	if err != nil {
		log.Err(err).Interface("ou", ou).Interface("approval", approval).Msg("UpdateTriggerActionApproval: failed")
		return err
	}
	return nil
}

const (
	allPass     = "all-pass"
	anyPass     = "any-pass"
	allFail     = "all-fail"
	anyFail     = "any-fail"
	mixedStatus = "mixed-status"
)

func checkTriggerOnEvalResults(value string, results []bool) bool {
	switch value {
	case allPass:
		for _, result := range results {
			if !result {
				return false
			}
		}
		return true
	case anyPass:
		for _, result := range results {
			if result {
				return true
			}
		}
		return false
	case allFail:
		// Return true if all values in results are false
		for _, result := range results {
			if result {
				return false
			}
		}
		return true
	case anyFail:
		// Return true if any value in results is false
		for _, result := range results {
			if !result {
				return true
			}
		}
		return false
	case mixedStatus:
		// Return true if there is a mix of true and false in results
		hasTrue := false
		hasFalse := false
		for _, result := range results {
			if result {
				hasTrue = true
			} else {
				hasFalse = true
			}
		}
		return hasTrue && hasFalse
	default:
		log.Warn().Str("value", value).Msg("checkTriggerOnEvalResults: unknown value")
		return false
	}
}

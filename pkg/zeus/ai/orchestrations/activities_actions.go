package ai_platform_service_orchestrations

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

// SendTriggerActionRequestForApproval sends the action request to the user for human in-the-loop approval
func (z *ZeusAiPlatformActivities) SendTriggerActionRequestForApproval(ctx context.Context) error {
	return nil
}

func (z *ZeusAiPlatformActivities) LookupEvalTriggerConditions(ctx context.Context, ou org_users.OrgUser, evalID int) ([]artemis_orchestrations.TriggerAction, error) {
	ta, err := artemis_orchestrations.SelectTriggerActionsByOrgAndOptParams(ctx, ou, evalID)
	if err != nil {
		log.Err(err).Msg("LookupEvalTriggerConditions: failed to lookup trigger actions")
		return nil, err
	}
	return ta, nil
}

func (z *ZeusAiPlatformActivities) CreateOrUpdateTriggerActionToExec(ctx context.Context, mb *MbChildSubProcessParams, act *artemis_orchestrations.TriggerAction) error {
	if act == nil || mb == nil {
		return nil
	}
	for _, tra := range act.TriggerActionsApprovals {
		err := artemis_orchestrations.CreateOrUpdateTriggerActionApproval(ctx, mb.Ou, &tra)
		if err != nil {
			log.Err(err).Msg("CreateOrUpdateTriggerActionToExec: failed to create or update trigger action approval")
			return err
		}
	}
	return nil
}

func (z *ZeusAiPlatformActivities) CheckEvalTriggerCondition(ctx context.Context, act *artemis_orchestrations.TriggerAction, emr *artemis_orchestrations.EvalMetricsResults) (*artemis_orchestrations.TriggerAction, error) {
	if act == nil || emr == nil || emr.EvalMetricsResults == nil {
		return act, nil
	}
	m := make(map[string][]bool)

	for _, er := range emr.EvalMetricsResults {
		if _, ok := m[er.EvalState]; !ok {
			m[er.EvalState] = []bool{}
		}
		if er.EvalMetricResult == nil || er.EvalMetricResult.EvalResultOutcomeBool == nil {
			continue
		}
		m[er.EvalState] = append(m[er.EvalState], aws.ToBool(er.EvalMetricResult.EvalResultOutcomeBool))
	}
	// gets the eval results by state, eg. info, trigger, etc.
	for _, tr := range act.EvalTriggerActions {
		results := m[tr.EvalTriggerState]
		// when the trigger on eval results condition is met, create a trigger action for approval
		if checkTriggerOnEvalResults(tr.EvalResultsTriggerOn, results) {
			tap := artemis_orchestrations.TriggerActionsApproval{
				WorkflowResultID: emr.EvalContext.AIWorkflowAnalysisResult.WorkflowResultID,
				EvalID:           tr.EvalID,
				TriggerID:        tr.TriggerID,
				ApprovalState:    "pending",
			}
			act.TriggerActionsApprovals = append(act.TriggerActionsApprovals, tap)
		}
	}
	return act, nil
}

func (z *ZeusAiPlatformActivities) SaveTriggerResponseOutput(ctx context.Context, trrr artemis_orchestrations.AIWorkflowTriggerResultResponse) error {
	respID, err := artemis_orchestrations.InsertOrUpdateAIWorkflowTriggerResultResponse(ctx, trrr)
	if err != nil {
		log.Err(err).Interface("respID", respID).Interface("trrr", trrr).Msg("SaveTriggerResponseOutput: failed")
		return err
	}
	return nil
}

func checkTriggerOnEvalResults(value string, results []bool) bool {
	switch value {
	case "all-pass":
		for _, result := range results {
			if !result {
				return false
			}
		}
		return true
	case "any-pass":
		for _, result := range results {
			if result {
				return true
			}
		}
		return false
	case "all-fail":
		// Return true if all values in results are false
		for _, result := range results {
			if result {
				return false
			}
		}
		return true
	case "any-fail":
		// Return true if any value in results is false
		for _, result := range results {
			if !result {
				return true
			}
		}
		return false
	case "mixed-status":
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
		return false
	}
}

package ai_platform_service_orchestrations

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (z *ZeusAiPlatformActivities) SelectRetrievalTask(ctx context.Context, ou org_users.OrgUser, retID int) ([]artemis_orchestrations.RetrievalItem, error) {
	resp, err := artemis_orchestrations.SelectRetrievals(ctx, ou, retID)
	if err != nil {
		log.Err(err).Interface("resp", resp).Int("retID", retID).Msg("SelectRetrievalTask: failed")
		return resp, err
	}
	return resp, nil
}

func (z *ZeusAiPlatformActivities) CreateWsr(ctx context.Context, cp *MbChildSubProcessParams) (*MbChildSubProcessParams, error) {
	wio := WorkflowStageIO{
		WorkflowStageReference: cp.Wsr,
		WorkflowExecParams:     cp.WfExecParams,
		WorkflowStageInfo: WorkflowStageInfo{
			PromptReduction: &PromptReduction{
				MarginBuffer:          cp.Tc.MarginBuffer,
				Model:                 cp.Tc.Model,
				TokenOverflowStrategy: cp.Tc.TokenOverflowStrategy,
			},
		},
	}
	wio.Org.OrgID = cp.Ou.OrgID
	_, err := s3ws(ctx, cp, &wio)
	if err != nil {
		log.Err(err).Msg("AiRetrievalTask: failed")
		return nil, err
	}
	return cp, nil
}

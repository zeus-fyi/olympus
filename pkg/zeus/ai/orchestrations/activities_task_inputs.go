package ai_platform_service_orchestrations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
)

type AggRetResp struct {
	AIWorkflowAnalysisResultSlice []artemis_orchestrations.AIWorkflowAnalysisResult
	InputDataAnalysisToAggSlice   []InputDataAnalysisToAgg
}

/*
	tmp := AiAggregateAnalysisRetrievalTaskInputDebug{
		SourceTaskIds: sourceTaskIds,
		Cp:            cp,
	}
	tmp.Save()
*/

func (z *ZeusAiPlatformActivities) AiAggregateAnalysisRetrievalTask(ctx context.Context, cp *MbChildSubProcessParams, sourceTaskIds []int) (*MbChildSubProcessParams, error) {
	results, err := artemis_orchestrations.SelectAiWorkflowAnalysisResults(ctx, cp.Window, []int{cp.Oj.OrchestrationID}, sourceTaskIds)
	if err != nil {
		log.Err(err).Msg("AiAggregateAnalysisRetrievalTask: SelectAiWorkflowAnalysisResults failed")
		return nil, err
	}
	var resp []InputDataAnalysisToAgg
	for _, r := range results {
		b, berr := gs3wfsCustomTaskName(ctx, cp, fmt.Sprintf("%d", r.WorkflowResultID))
		if berr != nil {
			log.Err(berr).Msg("AiAggregateAnalysisRetrievalTask: failed")
			continue
		}
		if b == nil {
			continue
		}
		tmp := InputDataAnalysisToAgg{}
		jerr := json.Unmarshal(b.Bytes(), &tmp)
		if jerr != nil {
			log.Err(jerr).Msg("AiAggregateAnalysisRetrievalTask: failed")
			continue
		}
		resp = append(resp, tmp)
	}

	log.Info().Interface("len(results)", results).Msg("AiAggregateAnalysisRetrievalTask")
	wio := WorkflowStageIO{
		WorkflowStageReference: cp.Wsr,
		WorkflowStageInfo: WorkflowStageInfo{
			PromptReduction: &PromptReduction{
				MarginBuffer:              cp.Tc.MarginBuffer,
				Model:                     cp.Tc.Model,
				TokenOverflowStrategy:     cp.Tc.TokenOverflowStrategy,
				DataInAnalysisAggregation: resp,
			},
		},
	}
	_, err = s3ws(ctx, cp, &wio)
	if err != nil {
		log.Err(err).Msg("AiRetrievalTask: failed")
		return nil, err
	}
	return cp, nil
}

package ai_platform_service_orchestrations

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
)

type AggRetResp struct {
	AIWorkflowAnalysisResultSlice []artemis_orchestrations.AIWorkflowAnalysisResult
	InputDataAnalysisToAggSlice   []InputDataAnalysisToAgg
}

func (z *ZeusAiPlatformActivities) AiAggregateAnalysisRetrievalTask(ctx context.Context, cp *MbChildSubProcessParams, sourceTaskIds []int) (*MbChildSubProcessParams, error) {
	results, err := artemis_orchestrations.SelectAiWorkflowAnalysisResults(ctx, cp.Window, []int{cp.Oj.OrchestrationID}, sourceTaskIds)
	if err != nil {
		log.Err(err).Msg("AiAggregateAnalysisRetrievalTask: SelectAiWorkflowAnalysisResults failed")
		return nil, err
	}
	var resp []InputDataAnalysisToAgg
	if cp.WfExecParams.WorkflowOverrides.IsUsingFlows {
		for _, vi := range cp.WfExecParams.WorkflowTasks {
			if vi.AnalysisTaskName != "" {
				tn := cp.Tc.TaskName
				cp.Tc.TaskName = vi.AnalysisTaskName
				wso, werr := gs3wfs(ctx, cp)
				if werr != nil {
					log.Err(err).Msg("AiAggregateAnalysisRetrievalTask: SelectAiWorkflowAnalysisResults failed")
					return nil, err
				}
				resp = append(resp, wso.InputDataAnalysisToAgg)
				cp.Tc.TaskName = tn
			}
		}
	} else {
		for _, r := range results {
			b, berr := json.Marshal(r.Metadata)
			if berr != nil {
				log.Err(berr).Msg("AiAggregateAnalysisRetrievalTask: failed")
				continue
			}
			if b == nil || string(b) == "null" {
				continue
			}
			tmp := InputDataAnalysisToAgg{}
			jerr := json.Unmarshal(b, &tmp)
			if jerr != nil {
				log.Err(jerr).Msg("AiAggregateAnalysisRetrievalTask: failed")
				continue
			}
			resp = append(resp, tmp)
		}
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

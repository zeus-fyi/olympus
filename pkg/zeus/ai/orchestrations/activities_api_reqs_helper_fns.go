package ai_platform_service_orchestrations

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	"go.temporal.io/sdk/activity"
)

const (
	iterateApiReq = "iterate"
	iterateQpOnly = "iterate-qp-only"
	bulkApiReq    = "bulk"
)

func (z *ZeusAiPlatformActivities) FanOutApiCallRequestTask(ctx context.Context, rts []iris_models.RouteInfo, cp *MbChildSubProcessParams) (*MbChildSubProcessParams, error) {
	echoReqs := getPayloads(cp)
	wio, werr := gs3wfs(ctx, cp)
	if werr != nil {
		log.Err(werr).Msg("TokenOverflowReduction: failed to select workflow io")
		return nil, werr
	}
	log.Info().Interface("inputID", cp.Wsr.InputID).Interface("wio.WorkflowStageInfo.ApiIterationCount", wio.WorkflowStageInfo.ApiIterationCount).Msg("TokenOverflowReduction: wio")
	if wio.WorkflowStageInfo.ApiIterationCount > 0 {
		cp.Tc.ApiIterationCount = wio.WorkflowStageInfo.ApiIterationCount
	}
	na := NewZeusAiPlatformActivities()
	for _, rtas := range rts {
		rt := RouteTask{
			Ou:        cp.Ou,
			Retrieval: cp.Tc.Retrieval,
			RouteInfo: rtas,
		}
		switch getRetOpt(cp, echoReqs) {
		case iterateApiReq, iterateQpOnly:
			for pi, ple := range echoReqs {
				if pi <= cp.Tc.ApiIterationCount && cp.Tc.ApiIterationCount > 0 {
					continue
				}
				log.Info().Interface("pi", pi).Msg("FanOutApiCallRequestTask: starting")
				cp.Tc.ApiIterationCount = pi
				rt.RouteInfo.Payload = ple
				_, err := na.ApiCallRequestTask(ctx, rt, cp)
				if err != nil {
					log.Err(err).Interface("pi", pi).Msg("FanOutApiCallRequestTask: failed")
					return nil, err
				}
				activity.RecordHeartbeat(ctx, fmt.Sprintf("iterate-%d", pi))
			}
		case bulkApiReq:
			rt.RouteInfo.Payloads = echoReqs
			_, err := na.ApiCallRequestTask(ctx, rt, cp)
			if err != nil {
				log.Err(err).Interface("len(pls)", len(echoReqs)).Msg("FanOutApiCallRequestTask: bulk failed")
				return nil, err
			}
			activity.RecordHeartbeat(ctx, fmt.Sprintf("bulk"))
		default:
			if len(echoReqs) > 1 {
				rt.RouteInfo.Payloads = echoReqs
			} else if len(echoReqs) == 1 {
				rt.RouteInfo.Payload = echoReqs[0]
			}
			_, err := na.ApiCallRequestTask(ctx, rt, cp)
			if err != nil {
				log.Err(err).Msg("FanOutApiCallRequestTask: default failed")
				return nil, err
			}
			activity.RecordHeartbeat(ctx, fmt.Sprintf("default"))
		}
	}
	cp.Tc.RegexSearchResults = nil
	cp.Tc.ApiResponseResults = nil
	cp.Tc.JsonResponseResults = nil
	return cp, nil
}

func SaveResult(ctx context.Context, cp *MbChildSubProcessParams, sg *hera_search.SearchResultGroup, sres hera_search.SearchResult, reqHash string) (*MbChildSubProcessParams, error) {
	if cp == nil || sg == nil {
		log.Warn().Msg("SaveResult: cp or sg is nil")
		return nil, nil
	}
	wio, werr := gs3wfs(ctx, cp)
	if werr != nil {
		log.Err(werr).Msg("TokenOverflowReduction: failed to select workflow io")
		return nil, werr
	}
	wio.ApiIterationCount = cp.Tc.ApiIterationCount
	if cp.WfExecParams.WorkflowOverrides.IsUsingFlows && wio.WorkflowStageInfo.WorkflowInCacheHash != nil && len(reqHash) > 0 {
		if _, ok := wio.WorkflowStageInfo.WorkflowInCacheHash[reqHash]; ok {
			log.Info().Interface("reqHash", reqHash).Msg("SaveResult: reqHash found in cache; skip adding again to wf result")
			return cp, nil
		}
	} else if cp.WfExecParams.WorkflowOverrides.IsUsingFlows && wio.WorkflowStageInfo.WorkflowInCacheHash == nil {
		icm := make(map[string]bool)
		if cp.WfExecParams.WorkflowOverrides.IsUsingFlows && len(reqHash) > 0 {
			icm[reqHash] = true
			wio.WorkflowStageInfo.WorkflowInCacheHash = icm
			wio.ApiIterationCount = cp.Tc.ApiIterationCount
		}
	}
	if cp.WfExecParams.WorkflowOverrides.IsUsingFlows && wio.WorkflowStageInfo.PromptReduction == nil {
		wio.ApiIterationCount = cp.Tc.ApiIterationCount
		wio.WorkflowStageInfo.PromptReduction = &PromptReduction{
			MarginBuffer:          cp.Tc.MarginBuffer,
			Model:                 cp.Tc.Model,
			TokenOverflowStrategy: cp.Tc.TokenOverflowStrategy,
			PromptReductionSearchResults: &PromptReductionSearchResults{
				InPromptBody:  cp.Tc.Prompt,
				InSearchGroup: sg,
			},
		}
	} else if wio.WorkflowStageInfo.PromptReduction.PromptReductionSearchResults == nil || wio.WorkflowStageInfo.PromptReduction.PromptReductionSearchResults.InSearchGroup == nil {
		wio.ApiIterationCount = cp.Tc.ApiIterationCount
		wio.WorkflowStageInfo.PromptReduction.PromptReductionSearchResults = &PromptReductionSearchResults{
			InPromptBody:  cp.Tc.Prompt,
			InSearchGroup: sg,
		}
	} else {
		wio.ApiIterationCount = cp.Tc.ApiIterationCount
		wio.AppendApiResponseResult(sres)
	}
	_, err := s3ws(ctx, cp, wio)
	if err != nil {
		log.Err(err).Msg("TokenOverflowReduction: failed to update workflow io")
		return nil, err
	}
	return cp, nil
}

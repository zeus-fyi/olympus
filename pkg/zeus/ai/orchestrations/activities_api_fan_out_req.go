package ai_platform_service_orchestrations

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"go.temporal.io/sdk/activity"
)

const (
	iterateApiReq = "iterate"
	iterateQpOnly = "iterate-qp-only"
	bulkApiReq    = "bulk"
)

func FanOutInitSetup(ctx context.Context, mb *MbChildSubProcessParams) (map[int]map[int]bool, error) {
	na := NewZeusAiPlatformActivities()
	if mb.Tc.Retrieval.RetrievalID == nil {
		return nil, fmt.Errorf("no retrieval found")
	}
	rev, rerr := na.SelectRetrievalTask(ctx, mb.Ou, *mb.Tc.Retrieval.RetrievalID)
	if rerr != nil {
		log.Err(rerr).Interface("rev", rev).Msg("FanOutApiCallRequestTask: SelectRetrievalTask failed")
		return nil, rerr
	}
	log.Info().Interface("*mb.Tc.Retrieval.RetrievalID", *mb.Tc.Retrieval.RetrievalID).Msg("*mb.Tc.Retrieval.RetrievalID")
	sv, serr := artemis_orchestrations.SelectRetrievalResultsIds(ctx, mb.Window, []int{mb.Oj.OrchestrationID}, []int{*mb.Tc.Retrieval.RetrievalID})
	if serr != nil {
		log.Err(serr).Msg("FanOutApiCallRequestTask: SelectRetrievalResultsIds failed")
		return nil, serr
	}
	sm := make(map[int]map[int]bool)
	for _, vi := range sv {
		if _, ok := sm[vi.ChunkOffset]; !ok {
			sm[vi.ChunkOffset] = make(map[int]bool)
		}
		sm[vi.ChunkOffset][vi.IterationCount] = true
	}
	return sm, nil
}

func getPendingRetWrAndIter(cp *MbChildSubProcessParams, iter, tc int) *artemis_orchestrations.AIWorkflowRetrievalResult {
	ch := chronos.Chronos{}
	wr := &artemis_orchestrations.AIWorkflowRetrievalResult{
		WorkflowResultID:      ch.UnixTimeStampNow(),
		OrchestrationID:       cp.Oj.OrchestrationID,
		RetrievalID:           aws.ToInt(cp.Tc.Retrieval.RetrievalID),
		IterationCount:        iter,
		ChunkOffset:           cp.Tc.ChunkIterator,
		RunningCycleNumber:    cp.Wsr.RunCycle,
		SearchWindowUnixStart: cp.Window.UnixStartTime,
		SearchWindowUnixEnd:   cp.Window.UnixEndTime,
		Status:                fmt.Sprintf("complete (%d/%d)", iter+1, tc),
		SkipRetrieval:         false,
	}
	return wr
}

func saveRetrievalResp(ctx context.Context, mb *MbChildSubProcessParams, searchResult hera_search.SearchResult) error {
	if mb == nil || mb.Tc.WorkflowRetrievalResult == nil || mb.Tc.WorkflowRetrievalResult.WorkflowResultID == 0 {
		return fmt.Errorf("saveRetrievalResp: MbChildSubProcessParams & WorkflowRetrievalResult required")
	}
	err := s3wsCustomTaskName(ctx, mb, fmt.Sprintf("%d", mb.Tc.WorkflowRetrievalResult.WorkflowResultID), searchResult)
	if err != nil {
		log.Err(err).Msg("saveRetrievalResp: s3wsCustomTaskName: saveCsvResp failed")
		return err
	}
	err = artemis_orchestrations.InsertWorkflowRetrievalResult(ctx, mb.Tc.WorkflowRetrievalResult)
	if err != nil {
		log.Err(err).Interface("wr", mb.Tc.WorkflowRetrievalResult).Interface("sr", searchResult).Msg("saveCsvResp: InsertWorkflowRetrievalResult failed")
		return err
	}
	return nil
}

func saveRetrievalRespErr(ctx context.Context, mb *MbChildSubProcessParams) error {
	sv := mb.Tc.WorkflowRetrievalResult.Status
	mb.Tc.WorkflowRetrievalResult.Status = "error"
	err := artemis_orchestrations.InsertWorkflowRetrievalResultError(ctx, mb.Tc.WorkflowRetrievalResult)
	if err != nil {
		log.Err(err).Interface("wr", mb.Tc.WorkflowRetrievalResult).Msg("saveRetrievalRespErr: InsertWorkflowRetrievalResult failed")
		return err
	}
	mb.Tc.WorkflowRetrievalResult.Status = sv
	return err
}

func (z *ZeusAiPlatformActivities) FanOutApiCallRequestTask(ctx context.Context, rts []iris_models.RouteInfo, cp *MbChildSubProcessParams) (*MbChildSubProcessParams, error) {
	sm, serr := FanOutInitSetup(ctx, cp)
	if serr != nil {
		log.Err(serr).Msg("FanOutApiCallRequestTask: SelectRetrievalResultsIds failed")
		return nil, serr
	}
	na := NewZeusAiPlatformActivities()
	echoReqs := getPayloads(cp)
	wio, werr := gs3wfs(ctx, cp)
	if werr != nil {
		log.Err(werr).Msg("FanOutApiCallRequestTask: failed to select workflow io")
		return nil, werr
	}
	log.Info().Interface("inputID", cp.Wsr.InputID).Interface("wio.WorkflowStageInfo.ApiIterationCount", wio.WorkflowStageInfo.ApiIterationCount).Msg("TokenOverflowReduction: wio")
	for _, rtas := range rts {
		rt := RouteTask{
			Ou:        cp.Ou,
			Retrieval: cp.Tc.Retrieval,
			RouteInfo: rtas,
		}
		switch getRetOpt(cp, echoReqs) {
		case iterateApiReq, iterateQpOnly:
			for pi, ple := range echoReqs {
				if tv, ok := sm[cp.Tc.ChunkIterator][pi]; ok && tv {
					continue
				}
				log.Info().Interface("pi", pi).Msg("FanOutApiCallRequestTask: starting")
				cp.Tc.WorkflowRetrievalResult = getPendingRetWrAndIter(cp, pi, len(echoReqs))
				rt.RouteInfo.Payload = ple
				_, err := na.ApiCallRequestTask(ctx, rt, cp)
				if err != nil {
					if cp.Tc.WorkflowRetrievalResult.Attempts < 7 {
						cp.Tc.WorkflowRetrievalResult.Attempts += 1
						time.Sleep(time.Duration(cp.Tc.WorkflowRetrievalResult.Attempts) * time.Minute)
					} else {
						rerrr := saveRetrievalRespErr(ctx, cp)
						if rerrr != nil {
							log.Err(rerrr).Interface("pi", pi).Msg("FanOutApiCallRequestTask: saveRetrievalRespErr failed")
							return nil, rerrr
						}
					}
					log.Err(err).Interface("pi", pi).Interface("cp.Tc.WorkflowRetrievalResult.Attempts", cp.Tc.WorkflowRetrievalResult.Attempts).Msg("FanOutApiCallRequestTask: failed")
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

	log.Info().Interface("len(cp.Tc.RegexSearchResults)", len(cp.Tc.RegexSearchResults)).Interface("len(cp.Tc.ApiResponseResults)", len(cp.Tc.ApiResponseResults)).Msg("FanOutApiCallRequestTask {")
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
	err := saveRetrievalResp(ctx, cp, sres)
	if err != nil {
		log.Err(err).Msg("SaveResult: saveRetrievalResp failed to save")
		return nil, err
	}
	wio, werr := gs3wfs(ctx, cp)
	if werr != nil {
		log.Err(werr).Msg("SaveResult: failed to select workflow io")
		return nil, werr
	}
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
		}
	}
	if cp.WfExecParams.WorkflowOverrides.IsUsingFlows && wio.WorkflowStageInfo.PromptReduction == nil {
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
		wio.WorkflowStageInfo.PromptReduction.PromptReductionSearchResults = &PromptReductionSearchResults{
			InPromptBody:  cp.Tc.Prompt,
			InSearchGroup: sg,
		}
	} else {
		wio.AppendApiResponseResult(sres)
	}
	_, err = s3ws(ctx, cp, wio)
	if err != nil {
		log.Err(err).Msg("SaveResult: failed to update workflow io")
		return nil, err
	}
	return cp, nil
}

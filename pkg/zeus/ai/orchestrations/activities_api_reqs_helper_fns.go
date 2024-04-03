package ai_platform_service_orchestrations

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
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

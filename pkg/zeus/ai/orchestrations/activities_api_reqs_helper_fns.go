package ai_platform_service_orchestrations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	iris_api_requests "github.com/zeus-fyi/olympus/pkg/iris/proxy/orchestrations/api_requests"
	"go.temporal.io/sdk/activity"
)

func (z *ZeusAiPlatformActivities) FanOutApiCallRequestTask(ctx context.Context, rts []iris_models.RouteInfo, cp *MbChildSubProcessParams) (*MbChildSubProcessParams, error) {
	var echoReqs []map[string]interface{}
	if cp.WfExecParams.WorkflowOverrides.RetrievalOverrides != nil {
		if v, ok := cp.WfExecParams.WorkflowOverrides.RetrievalOverrides[cp.Tc.Retrieval.RetrievalName]; ok {
			for _, pl := range v.Payloads {
				echoReqs = append(echoReqs, pl)
			}
		}
	}
	na := NewZeusAiPlatformActivities()
	retOpt := "default"
	if cp.Tc.Retrieval.WebFilters != nil && cp.Tc.Retrieval.WebFilters.PayloadPreProcessing != nil && len(echoReqs) > 0 {
		retOpt = *cp.Tc.Retrieval.WebFilters.PayloadPreProcessing
	}
	wio, werr := gs3wfs(ctx, cp)
	if werr != nil {
		log.Err(werr).Msg("TokenOverflowReduction: failed to select workflow io")
		return nil, werr
	}
	log.Info().Interface("inputID", cp.Wsr.InputID).Interface("wio.WorkflowStageInfo.ApiIterationCount", wio.WorkflowStageInfo.ApiIterationCount).Msg("TokenOverflowReduction: wio")
	if wio.WorkflowStageInfo.ApiIterationCount > 0 {
		cp.Tc.ApiIterationCount = wio.WorkflowStageInfo.ApiIterationCount
	}
	for _, rtas := range rts {
		rt := RouteTask{
			Ou:        cp.Ou,
			Retrieval: cp.Tc.Retrieval,
			RouteInfo: rtas,
		}
		switch retOpt {
		case "iterate", "iterate-qp-only":
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
		case "bulk":
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

func cacheIfCond(ctx context.Context, cp *MbChildSubProcessParams, r RouteTask, req *iris_api_requests.ApiProxyRequest, reqHash string, reqCached bool) {
	var rg string
	if r.Retrieval.WebFilters != nil && r.Retrieval.WebFilters.RoutingGroup != nil {
		rg = *r.Retrieval.WebFilters.RoutingGroup
	}
	if req.StatusCode >= 200 && req.StatusCode < 300 && cp.WfExecParams.WorkflowOverrides.IsUsingFlows && !reqCached {
		ht, err := artemis_entities.HashWebRequestResultsAndParams(r.Ou, r.RouteInfo)
		if err != nil {
			log.Err(err).Msg("ApiCallRequestTask: failed to hash request cache")
		}
		log.Info().Interface("hash", ht.RequestCache).Msg("ApiCallRequestTask: request cache end")
		if ht.RequestCache != "" {
			reqHash = ht.RequestCache
			uew := artemis_entities.UserEntity{
				Nickname: ht.RequestCache,
				Platform: rg,
			}
			if len(req.Response) > 0 {
				b, cerr := json.Marshal(req.Response)
				if cerr != nil {
					log.Err(cerr).Msg("ApiCallRequestTask: failed to marshal response")
				}
				if b != nil && string(b) != "null" {
					uew.MdSlice = append(uew.MdSlice, artemis_entities.UserEntityMetadata{
						JsonData: b,
					})
				}
			}
			if len(req.RawResponse) > 0 && len(uew.MdSlice) <= 0 {
				uew.MdSlice = append(uew.MdSlice, artemis_entities.UserEntityMetadata{
					TextData: aws.String(string(req.RawResponse)),
				})
			}
			if len(uew.MdSlice) > 0 {
				_, err = s3globalWf(ctx, cp, &uew)
				if err != nil {
					log.Err(err).Msg("ApiCallRequestTask: s3globalWf err")
				}
			}
		}
	}
	return
}

func checkIfCached(ctx context.Context, cp *MbChildSubProcessParams, r RouteTask, req *iris_api_requests.ApiProxyRequest) (string, bool) {
	reqCached := false
	var rg string
	if r.Retrieval.WebFilters != nil && r.Retrieval.WebFilters.RoutingGroup != nil {
		rg = *r.Retrieval.WebFilters.RoutingGroup
	}

	if !cp.WfExecParams.WorkflowOverrides.IsUsingFlows || len(rg) == 0 {
		return "", false
	}
	// cache is either
	/*
		1. err = json.Unmarshal(uew.MdSlice[0].JsonData, &req.Response)
		2. req.RawResponse = []byte(*uew.MdSlice[0].TextData)

		check date last modified; if > 30 days then replace/ set no cache
	*/
	ht, err := artemis_entities.HashWebRequestResultsAndParams(r.Ou, r.RouteInfo)
	if err != nil {
		log.Err(err).Msg("ApiCallRequestTask: failed to hash request cache")
		return "", false
	}
	log.Info().Interface("hash", ht.RequestCache).Msg("start")
	if len(ht.RequestCache) > 0 {
		uew := &artemis_entities.UserEntity{
			Nickname: ht.RequestCache,
			Platform: rg,
		}
		uew, err = gs3globalWf(ctx, cp, uew)
		if err != nil {
			log.Err(err).Msg("ApiCallRequestTask: failed to unmarshal response")
			err = nil
		} else {
			log.Info().Interface("mdslicelen", len(uew.MdSlice)).Msg("FanOutApiCallRequestTask: uew")
			if len(uew.MdSlice) > 0 && uew.MdSlice[0].JsonData != nil {
				tmp := uew.MdSlice[0].JsonData
				if string(tmp) != "null" {
					err = json.Unmarshal(uew.MdSlice[0].JsonData, &req.Response)
					if err != nil {
						log.Err(err).Msg("ApiCallRequestTask: failed to unmarshal response")
					} else {
						log.Info().Interface("hash", ht.RequestCache).Interface("len(uew.MdSlice)", uew.MdSlice[0].JsonData).Msg("FanOutApiCallRequestTask: json cache found skipping")
						reqCached = true
					}
				}
			} else if len(uew.MdSlice) > 0 && uew.MdSlice[0].TextData != nil && *uew.MdSlice[0].TextData != "" {
				reqCached = true
				req.RawResponse = []byte(*uew.MdSlice[0].TextData)
				log.Info().Interface("hash", ht.RequestCache).Interface("len(uew.MdSlice)", uew.MdSlice[0].TextData).Msg("FanOutApiCallRequestTask: text cache found skipping")
			}
		}
	}
	return ht.RequestCache, reqCached
}

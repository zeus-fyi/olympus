package ai_platform_service_orchestrations

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
	iris_api_requests "github.com/zeus-fyi/olympus/pkg/iris/proxy/orchestrations/api_requests"
)

//739ae3a2f25dde86e51f455621aa5213b8f78118793a91b9205217d5ad3e4328739ae3a2f25dde86e51f455621aa5213b8f78118793a91b9205217d5ad3e4328

func cacheIfCond(ctx context.Context, cp *MbChildSubProcessParams, r RouteTask, req *iris_api_requests.ApiProxyRequest, reqHash string, reqCached bool) {
	var rg string
	if r.Retrieval.WebFilters != nil && r.Retrieval.WebFilters.RoutingGroup != nil {
		rg = *r.Retrieval.WebFilters.RoutingGroup
	}
	var uew artemis_entities.UserEntity
	if req.StatusCode >= 200 && req.StatusCode < 300 && cp.WfExecParams.WorkflowOverrides.IsUsingFlows && !reqCached {
		ht, err := artemis_entities.HashWebRequestResultsAndParams(r.Ou, r.RouteInfo)
		if err != nil {
			log.Err(err).Msg("ApiCallRequestTask: failed to hash request cache")
		}
		log.Info().Interface("hash", ht.RequestCache).Msg("ApiCallRequestTask: request cache end")
		if ht.RequestCache != "" {
			reqHash = ht.RequestCache
			uew = artemis_entities.UserEntity{
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
		}
	}
	if len(uew.MdSlice) > 0 {
		_, err := s3globalWf(ctx, cp, &uew)
		if err != nil {
			log.Err(err).Msg("ApiCallRequestTask: s3globalWf err")
		}
	}
	return
}

func checkIfCached(ctx context.Context, cp *MbChildSubProcessParams, r RouteTask, req *iris_api_requests.ApiProxyRequest) (string, bool) {
	var reqCached bool
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
	var uew *artemis_entities.UserEntity
	ht, err := artemis_entities.HashWebRequestResultsAndParams(r.Ou, r.RouteInfo)
	if err != nil {
		log.Err(err).Msg("ApiCallRequestTask: failed to hash request cache")
		return "", false
	}
	log.Info().Interface("hash", ht.RequestCache).Msg("start")
	if len(ht.RequestCache) > 0 {
		uew = &artemis_entities.UserEntity{
			Nickname: ht.RequestCache,
			Platform: rg,
		}
		//// clean cache debug
		//p, _ := globalWfEntityStageNamePath(cp, uew)
		//deleteFromS3(ctx, p)
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
						log.Info().Interface("hash", ht.RequestCache).Msg("FanOutApiCallRequestTask: json cache found skipping")
						//log.Info().Interface("hash", ht.RequestCache).Interface("len(uew.MdSlice)", uew.MdSlice[0].JsonData).Msg("FanOutApiCallRequestTask: json cache found skipping")
						reqCached = true
					}
				}
			} else if len(uew.MdSlice) > 0 && uew.MdSlice[0].TextData != nil && *uew.MdSlice[0].TextData != "" {
				reqCached = true
				req.RawResponse = []byte(*uew.MdSlice[0].TextData)
				log.Info().Interface("hash", ht.RequestCache).Msg("FanOutApiCallRequestTask: text cache found skipping")
				//log.Info().Interface("hash", ht.RequestCache).Interface("len(uew.MdSlice)", uew.MdSlice[0].TextData).Msg("FanOutApiCallRequestTask: text cache found skipping")
			}
		}
	}
	return ht.RequestCache, reqCached
}

// to debug ^ and comment out caching
// reqCached

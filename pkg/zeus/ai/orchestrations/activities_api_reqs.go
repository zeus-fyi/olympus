package ai_platform_service_orchestrations

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	"github.com/zeus-fyi/olympus/pkg/aegis/aws_secrets"
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
				cp.Tc.ApiIterationCount = pi
				log.Info().Interface("pi", pi).Msg("FanOutApiCallRequestTask: ple")
				rt.RouteInfo.Payload = ple
				_, err := na.ApiCallRequestTask(ctx, rt, cp)
				if err != nil {
					log.Err(err).Msg("FanOutApiCallRequestTask: failed")
					return nil, err
				}
				activity.RecordHeartbeat(ctx, fmt.Sprintf("iterate-%d", pi))
			}
		case "bulk":
			rt.RouteInfo.Payloads = echoReqs
			_, err := na.ApiCallRequestTask(ctx, rt, cp)
			if err != nil {
				log.Err(err).Msg("FanOutApiCallRequestTask: bulk failed")
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

func (z *ZeusAiPlatformActivities) ApiCallRequestTask(ctx context.Context, r RouteTask, cp *MbChildSubProcessParams) (*MbChildSubProcessParams, error) {
	retInst := r.Retrieval
	if retInst.WebFilters == nil || retInst.WebFilters.RoutingGroup == nil || len(*retInst.WebFilters.RoutingGroup) <= 0 {
		return nil, nil
	}
	rg := *retInst.WebFilters.RoutingGroup
	restMethod := http.MethodGet
	if retInst.WebFilters.EndpointREST != nil {
		restMethod = strings.ToLower(*retInst.WebFilters.EndpointREST)
		switch restMethod {
		case "post", "POST":
			restMethod = http.MethodPost
		case "put", "PUT":
			restMethod = http.MethodPut
		case "delete":
			restMethod = http.MethodDelete
		case "patch":
			restMethod = http.MethodPatch
		case "get":
			restMethod = http.MethodGet
		default:
			log.Info().Str("restMethod", restMethod).Msg("ApiCallRequestTask: rest method")
		}
	}
	var routeExt string
	var orgRouteExt string
	if cp.Tc.Retrieval.WebFilters != nil && cp.Tc.Retrieval.WebFilters.EndpointRoutePath != nil {
		orgRouteExt = *cp.Tc.Retrieval.WebFilters.EndpointRoutePath
		routeExt = orgRouteExt
	}
	if r.RouteInfo.Payload != nil {
		rp, qps, err := ReplaceAndPassParams(routeExt, r.RouteInfo.Payload)
		if err != nil {
			log.Err(err).Msg("ApiCallRequestTask: failed to replace route path params")
			return nil, err
		}
		log.Info().Interface("qps", qps).Msg("ApiCallRequestTask: qps")
		r.Qps = qps
		routeExt = rp
		r.RouteInfo.RouteExt = rp
		if len(r.RouteInfo.Payload) == 0 {
			r.RouteInfo.Payload = nil
		}
	}
	if retInst.WebFilters.RequestHeaders != nil {
		if r.Headers == nil {
			r.Headers = make(http.Header)
		}
		for k, v := range retInst.WebFilters.RequestHeaders {
			r.Headers.Set(k, v)
		}
	}
	var sec []int
	if retInst.WebFilters.DontRetryStatusCodes != nil {
		sec = retInst.WebFilters.DontRetryStatusCodes
	}
	secretNameRefApi := fmt.Sprintf("api-%s", rg)
	var regexPatterns []string
	for _, rgp := range retInst.WebFilters.RegexPatterns {
		regexPatterns = append(regexPatterns, FixRegexInput(rgp))
	}
	rw := iris_api_requests.NewIrisApiRequestsActivities()
	req := &iris_api_requests.ApiProxyRequest{
		Url:                    r.RouteInfo.RoutePath,
		OrgID:                  r.Ou.OrgID,
		UserID:                 r.Ou.UserID,
		ExtRoutePath:           routeExt,
		Payload:                r.RouteInfo.Payload,
		Payloads:               r.RouteInfo.Payloads,
		PayloadTypeREST:        restMethod,
		RequestHeaders:         r.Headers,
		RegexFilters:           regexPatterns,
		SkipErrorOnStatusCodes: sec,
		SecretNameRef:          secretNameRefApi,
		IsFlowRequest:          cp.WfExecParams.WorkflowOverrides.IsUsingFlows,
	}
	log.Info().Interface("req.Url", req.Url).Msg("req value")
	reqCached := false
	if cp.WfExecParams.WorkflowOverrides.IsUsingFlows {
		// cache is either
		/*
			1. err = json.Unmarshal(uew.MdSlice[0].JsonData, &req.Response)
			2. req.RawResponse = []byte(*uew.MdSlice[0].TextData)

			check date last modified; if > 30 days then replace/ set no cache
		*/
		ht, err := artemis_entities.HashWebRequestResultsAndParams(r.Ou, r.RouteInfo)
		if err != nil {
			log.Err(err).Msg("ApiCallRequestTask: failed to hash request cache")
			return nil, err
		}
		log.Info().Interface("hash", ht.RequestCache).Msg("start")
		if len(ht.RequestCache) > 0 {
			uew := artemis_entities.UserEntity{
				Nickname: ht.RequestCache,
				Platform: rg,
			}
			_, err = gs3globalWf(ctx, cp, uew)
			if err != nil {
				log.Err(err).Msg("ApiCallRequestTask: failed to unmarshal response")
			}
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
	if !reqCached {
		rr, rrerr := rw.ExtLoadBalancerRequest(ctx, req)
		if rrerr != nil {
			if rr.StatusCode == 401 {
				// clear the cache
				log.Warn().Interface("routingTable", fmt.Sprintf("api-%s", *retInst.WebFilters.RoutingGroup)).Int("statusCode", rr.StatusCode).Msg("ApiCallRequestTask: clearing org secret cache")
				aws_secrets.ClearOrgSecretCache(r.Ou)
			}
			log.Err(rrerr).Interface("routingTable", fmt.Sprintf("api-%s", *retInst.WebFilters.RoutingGroup)).Msg("ApiCallRequestTask: failed to get response")
			return nil, rrerr
		}
		req = rr
	}
	var reqHash string
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
				_, err = s3globalWf(ctx, cp, uew)
				if err != nil {
					log.Err(err).Msg("ApiCallRequestTask: s3globalWf err")
				}
			}
		}
	}
	req.ExtRoutePath = orgRouteExt
	wr := hera_search.WebResponse{
		WebFilters: retInst.WebFilters,
		Body:       req.Response,
		RawMessage: req.RawResponse,
	}
	value := ""
	if wr.Body != nil && wr.RawMessage == nil {
		b, jer := json.Marshal(wr.Body)
		if jer != nil {
			log.Err(jer).Interface("routingTable", fmt.Sprintf("api-%s", *retInst.WebFilters.RoutingGroup)).Msg("ApiCallRequestTask: failed to get response")
			return nil, jer
		}
		value = fmt.Sprintf("%s", b)
	} else if wr.RawMessage != nil && len(req.RegexFilters) > 0 {
		value = fmt.Sprintf("%s", wr.RawMessage)
		wr.RegexFilteredBody = value
	} else if wr.Body != nil && wr.RawMessage != nil {
		value = fmt.Sprintf("%s", wr.RawMessage)
	}

	sres := hera_search.SearchResult{
		Source:      req.Url,
		Value:       value,
		QueryParams: r.Qps,
		Group:       aws.StringValue(retInst.WebFilters.RoutingGroup),
		WebResponse: wr,
	}
	sg := &hera_search.SearchResultGroup{
		PlatformName: cp.Tc.Retrieval.RetrievalPlatform,
		Window:       cp.Window,
	}
	sg.ApiResponseResults = []hera_search.SearchResult{sres}
	if req.RegexFilters != nil && len(req.RegexFilters) > 0 {
		cp.Tc.RegexSearchResults = append(cp.Tc.RegexSearchResults, sres)
	} else if len(r.Qps) > 0 {
		cp.Tc.RegexSearchResults = append(cp.Tc.RegexSearchResults, sres)
	} else {
		cp.Tc.ApiResponseResults = append(cp.Tc.ApiResponseResults, sres)
	}
	sg.SourceTaskID = cp.Tc.TaskID
	return SaveResult(ctx, cp, sg, sres, reqHash)
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
	}
	if cp.WfExecParams.WorkflowOverrides.IsUsingFlows && wio.WorkflowStageInfo.WorkflowInCacheHash == nil {
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
		if wio.WorkflowStageInfo.PromptReduction.PromptReductionSearchResults.InSearchGroup.ApiResponseResults == nil {
			wio.WorkflowStageInfo.PromptReduction.PromptReductionSearchResults.InSearchGroup.ApiResponseResults = make([]hera_search.SearchResult, 0)
		}
		wio.WorkflowStageInfo.PromptReduction.PromptReductionSearchResults.InSearchGroup.ApiResponseResults = append(wio.WorkflowStageInfo.PromptReduction.PromptReductionSearchResults.InSearchGroup.ApiResponseResults, sres)
	}
	_, err := s3ws(ctx, cp, wio)
	if err != nil {
		log.Err(err).Msg("TokenOverflowReduction: failed to update workflow io")
		return nil, err
	}
	return cp, nil
}

func FixRegexInput(input string) string {
	if len(input) > 0 {
		// Check if the first character is a backtick and replace it with a double quote
		if input[0] == '`' {
			input = "\"" + input[1:]
		}
		// Check if the last character is a backtick and replace it with a double quote
		if input[len(input)-1] == '`' {
			input = input[:len(input)-1] + "\""
		}
	}
	return input
}

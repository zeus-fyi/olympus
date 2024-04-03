package ai_platform_service_orchestrations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/rs/zerolog/log"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/pkg/aegis/aws_secrets"
	iris_api_requests "github.com/zeus-fyi/olympus/pkg/iris/proxy/orchestrations/api_requests"
)

func (z *ZeusAiPlatformActivities) ApiCallRequestTask(ctx context.Context, r RouteTask, cp *MbChildSubProcessParams) (*MbChildSubProcessParams, error) {
	retInst := r.Retrieval
	if retInst.WebFilters == nil || retInst.WebFilters.RoutingGroup == nil || len(*retInst.WebFilters.RoutingGroup) <= 0 {
		return nil, nil
	}
	rg := *retInst.WebFilters.RoutingGroup
	restMethod := getRestMethod(retInst)
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
	r = setHeaders(retInst, r)
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
		RegexFilters:           getRegexPatterns(retInst),
		SkipErrorOnStatusCodes: setDontRetryCodes(retInst),
		SecretNameRef:          fmt.Sprintf("api-%s", rg),
		IsFlowRequest:          cp.WfExecParams.WorkflowOverrides.IsUsingFlows,
	}
	log.Info().Interface("req.Url", req.Url).Msg("req value")
	reqHash, reqCached := checkIfCached(ctx, cp, r, req)
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
	cacheIfCond(ctx, cp, r, req, reqHash, reqCached)
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
		PlatformName:  cp.Tc.Retrieval.RetrievalPlatform,
		Window:        cp.Window,
		RetrievalName: aws.String(cp.Tc.Retrieval.RetrievalName),
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

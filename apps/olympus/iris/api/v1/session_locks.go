package v1_iris

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	iris_api_requests "github.com/zeus-fyi/olympus/pkg/iris/proxy/orchestrations/api_requests"
	iris_usage_meters "github.com/zeus-fyi/olympus/pkg/iris/proxy/usage_meters"
	iris_serverless "github.com/zeus-fyi/olympus/pkg/iris/serverless"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

// only used in tests

//func ProcessLockedSessionsHandler(c echo.Context) error {
//	request := new(ProxyRequest)
//	request.Body = echo.Map{}
//	if err := c.Bind(&request.Body); err != nil {
//		log.Err(err).Msg("proxy_anvil.ProcessLockedSessionsHandler: c.Bind")
//		return err
//	}
//	anvilHeader := c.Request().Header.Get(AnvilSessionLockHeader)
//	ou := org_users.OrgUser{}
//	ouc := c.Get("orgUser")
//	if ouc != nil {
//		ouser, ok := ouc.(org_users.OrgUser)
//		if ok {
//			ou = ouser
//		} else {
//			log.Warn().Interface("ou", ouser).Msg("proxy_anvil.ProcessLockedSessionsHandler: orgUser not found")
//			return c.JSON(http.StatusUnauthorized, Response{Message: "user not found"})
//		}
//	}
//	return request.ProcessLockedSessionRoute(c, ou.OrgID, anvilHeader, c.Request().Method, "", "test")
//}

/*

func (a *AnvilProxy) GetSessionLockedRoute(ctx context.Context, sessionID string) (string, error) {
	if sessionID == "Zeus-Test" {
		return "http://anvil.eeb335ad-78da-458f-9cfb-9928514d65d0.svc.cluster.local:8545", nil
	}
*/

/*
1. needs to have a registry of all anvil routes
  - this should be auto inserted when new anvil services are deployed

2a. needs to have a registry of all locked sessions
2b. needs to have a registry of all locked routes
3. needs to be able to dynamically add/remove anvil services
*/
const (
	anvilServerlessRoutesTableName = "anvil"
)

var cctx = zeus_common_types.CloudCtxNs{
	CloudProvider: "ovh",
	Region:        "us-west-or-1",
	Context:       "kubernetes-admin@zeusfyi",
	Namespace:     "anvil-serverless-4d383226",
}

const (
	internalOrgID = 7138983863666903883
)

func GetSessionLockedRoute(ctx context.Context, orgID int, sessionID, serverlessTableName, plan string) (string, bool, error) {
	if sessionID == "Zeus-Test" && orgID == internalOrgID {
		return "http://anvil.eeb335ad-78da-458f-9cfb-9928514d65d0.svc.cluster.local:8545", false, nil
	}
	if sessionID == "Zeus-Service-Test" && orgID == internalOrgID {
		return "http://anvil-49.anvil.anvil-serverless-aa3ffbd0.svc.cluster.local:8888", false, nil
	}
	route, isNewSession, err := iris_redis.IrisRedisClient.GetNextServerlessRoute(context.Background(), orgID, sessionID, serverlessTableName, plan)
	if err != nil {
		log.Err(err).Msg("proxy_anvil.SessionLocker.GetNextServerlessRoute")
		return route, isNewSession, err
	}

	if isNewSession {
		podName, perr := extractPodName(route)
		if perr != nil {
			log.Err(perr).Str("route", route).Msg("GetSessionLockedRoute: extractPodName")
			return "", isNewSession, perr
		}
		err = iris_serverless.IrisPlatformServicesWorker.ExecuteIrisServerlessPodRestartWorkflow(context.Background(), orgID, cctx, podName, serverlessTableName, sessionID, iris_redis.ServerlessSessionMaxRunTime)
		if err != nil {
			log.Err(err).Str("podName", podName).Msg("GetSessionLockedRoute: iris_serverless.IrisPlatformServicesWorker.ExecuteIrisServerlessPodRestartWorkflow")
			return "", isNewSession, err
		}
	}
	return route, isNewSession, err
}

func extractPodName(s string) (string, error) {
	// Parse the URL
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}

	// Split the host by "." to get the first part
	parts := strings.Split(u.Hostname(), ".")
	if len(parts) > 0 {
		return parts[0], nil
	}

	return "", fmt.Errorf("cannot extract pod name from url: %s", s)
}

var Env = "production"

func (p *ProxyRequest) ProcessLockedSessionRoute(c echo.Context, orgID int, sessionID, method, tempToken, plan string) error {
	endLockedSessionLease := c.Request().Header.Get(EndSessionLockHeader)
	if endLockedSessionLease == sessionID {
		// todo remove hardcoded table name
		return p.ProcessEndSessionLock(c, orgID, sessionID, anvilServerlessRoutesTableName)
	}
	routeURL, isNewSession, err := GetSessionLockedRoute(context.Background(), orgID, sessionID, anvilServerlessRoutesTableName, plan) // TODO remove hardcoded table name
	if len(routeURL) == 0 && strings.Contains(err.Error(), " max active sessions reached") {
		return c.JSON(http.StatusTooManyRequests, err)
	}
	if err != nil {
		log.Err(err).Msg("ProcessLockedSessionRoute: GetSessionLockedRoute")
		return c.JSON(http.StatusInternalServerError, err)
	}
	headers := make(http.Header)
	headers.Set(AnvilSessionLockHeader, tempToken)
	routeGroup := c.Request().Header.Get(RouteGroupHeader)
	if routeGroup != "" {
		headers.Set(RouteGroupHeader, routeGroup)
	}
	// for local testing:
	if Env == "production-local" || Env == "local" {
		routeURL = "http://localhost:8888"
	}

	if isNewSession {
		log.Info().Str("routeGroup", routeGroup).Msg("ProcessLockedSessionRoute: isNewSession")
		go func(orgID int) {
			err = iris_redis.IrisRedisClient.RecordNewServerlessRequestUsage(context.Background(), orgID)
			if err != nil {
				log.Err(err).Interface("orgID", orgID).Msg("ProcessRpcLoadBalancerRequest: iris_round_robin.RecordNewServerlessRequestUsage")
			}
		}(orgID)
	}

	if isNewSession && routeGroup != "" {
		// todo, just for anvil
		wa := web3_client.NewWeb3ClientFakeSigner(routeURL)
		wa.AddSessionLockHeader(tempToken)
		wa.IsAnvilNode = true
		wa.Dial()
		defer wa.Close()
		rpcNew := "http://localhost:8888/node"
		forkNumberHeader := c.Request().Header.Get(AnvilForkBlockNumberHeader)
		bn := 0
		if forkNumberHeader != "" {
			bn, err = strconv.Atoi(forkNumberHeader)
			if err != nil {
				bn = 0
				log.Err(err).Msg("ProcessLockedSessionRoute: strconv.Atoi")
				err = nil
			}
		}

		err = wa.ResetNetwork(context.Background(), rpcNew, bn)
		if err != nil {
			log.Err(err).Int("orgID", orgID).Str("sessionID", sessionID).Msg("ProxyRequest: ProcessLockedSessionRoute: wa.ResetNetwork")
			return c.JSON(http.StatusInternalServerError, nil)
		}
	}
	p.Body = GetSanitizedForkPayload(p.Body)
	req := &iris_api_requests.ApiProxyRequest{
		Url:             routeURL,
		OrgID:           orgID,
		PayloadTypeREST: method,
		Payload:         p.Body,
		IsInternal:      true,
		RequestHeaders:  headers,
		Timeout:         60 * time.Second,
	}
	rw := iris_api_requests.NewIrisApiRequestsActivities()
	resp, err := rw.ExtToAnvilInternalSimForkRequest(c.Request().Context(), req)
	if err != nil {
		log.Err(err).Str("routeURL", routeURL).Msg("rw.InternalSvcRelayRequest")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	go func(orgID int, usage *iris_usage_meters.PayloadSizeMeter) {
		if usage == nil {
			log.Warn().Msg("ProcessLockedSessionRoute: usage is nil")
			return
		}
		err = iris_redis.IrisRedisClient.RecordRequestUsage(context.Background(), orgID, usage)
		if err != nil {
			log.Err(err).Interface("orgID", orgID).Interface("usage", usage).Msg("ProcessRpcLoadBalancerRequest: iris_round_robin.IncrementResponseUsageRateMeter")
		}
	}(orgID, resp.PayloadSizeMeter)
	return c.JSON(http.StatusOK, resp.Response)
}

func (p *ProxyRequest) ProcessEndSessionLock(c echo.Context, orgID int, sessionID, serverlessRoutesTable string) error {
	path, perr := iris_redis.IrisRedisClient.ReleaseServerlessRoute(context.Background(), orgID, sessionID, serverlessRoutesTable)
	if perr == redis.Nil {
		return c.JSON(http.StatusOK, "ok")
	}
	if perr != nil {
		log.Err(perr).Msg("ProxyRequest: ProcessEndSessionLock")
		return c.JSON(http.StatusInternalServerError, nil)
	}

	if path != "" {
		podName, err := extractPodName(path)
		if err != nil {
			log.Err(err).Msg("ReleaseServerlessRoute: failed to extract pod name")
			return c.JSON(http.StatusInternalServerError, nil)
		}
		if len(podName) == 0 {
			log.Err(err).Msg("ReleaseServerlessRoute: pod name is empty")
			return c.JSON(http.StatusInternalServerError, nil)
		}
		err = iris_serverless.IrisPlatformServicesWorker.EarlyStart(context.Background(), orgID, podName, serverlessRoutesTable, sessionID)
		if err != nil {
			log.Err(err).Msg("ProxyRequest: ProcessEndSessionLock: iris_serverless.IrisPlatformServicesWorker.EarlyStart")
			err = nil
		}
		return c.JSON(http.StatusOK, fmt.Sprintf("released session lock-id %s", sessionID))
	}

	return c.JSON(http.StatusOK, nil)
}

func (p *ProxyRequest) Process(c echo.Context, r *iris_api_requests.ApiProxyRequest) error {
	if r == nil {
		return c.JSON(http.StatusBadRequest, errors.New("request is nil"))
	}
	resp, err := iris_api_requests.IrisProxyWorker.ExecuteIrisProxyWorkflow(c.Request().Context(), r)
	if err != nil {
		log.Err(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	if resp == nil {
		if r != nil {
			log.Warn().Str("route", r.Url).Msg("rw.InternalSvcRelayRequest")
		}
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, resp.Response)
}
func toForkingArg(jsonRpcURL string, blockNumber int) interface{} {
	arg := map[string]map[string]any{
		"forking": {
			"jsonRpcUrl":  jsonRpcURL,
			"blockNumber": blockNumber,
		},
	}
	return arg
}

func toForkingArgResetToLatest(jsonRpcURL string) interface{} {
	arg := map[string]map[string]any{
		"forking": {
			"jsonRpcUrl": jsonRpcURL,
		},
	}
	return arg
}

const (
	NodeProxy = "http://localhost:8888/node"
)

func GetSanitizedForkPayload(b echo.Map) echo.Map {
	method, ok := b["method"].(string)
	if ok && (method == "anvil_reset" || method == "hardhat_reset") {
		var np []interface{}
		params, pok := b["params"].([]interface{})
		if pok && len(params) > 0 {
			bp, bok := params[0].(map[string]interface{})
			if bok {
				for k, v := range bp {
					if k == "forking" {
						nestedMap, mok := v.(map[string]interface{})
						if mok {
							np = []interface{}{toForkingArgResetToLatest(NodeProxy)}
							for nk, nv := range nestedMap {
								if nk == "blockNumber" {
									np = []interface{}{toForkingArg(NodeProxy, int(nv.(float64)))}
								}
							}
						}
					}

				}
			}
		}
		b["params"] = np
	}
	if ok && strings.HasPrefix(method, "hardhat_") {
		b["method"] = replacePrefix(method, "hardhat_", "anvil_")
	}
	return b
}

func replacePrefix(input string, prefix string, replacement string) string {
	if strings.HasPrefix(input, prefix) {
		return replacement + input[len(prefix):]
	}
	return input
}

func DeleteSessionRequestHandler(c echo.Context) error {
	ouc, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusBadRequest, "org not found")
	}
	if ouc.OrgID <= 0 {
		return c.JSON(http.StatusBadRequest, "org not found")
	}
	endLockedSessionLease := c.Request().Header.Get(EndSessionLockHeader)
	sessionID := c.Param("sessionID")
	if len(sessionID) == 0 && len(endLockedSessionLease) == 0 {
		return c.JSON(http.StatusBadRequest, "sessionID is required via header or route param")
	}
	if len(endLockedSessionLease) > 0 && len(sessionID) == 0 {
		sessionID = endLockedSessionLease
	}
	p := ProxyRequest{}
	return p.ProcessEndSessionLock(c, ouc.OrgID, sessionID, anvilServerlessRoutesTableName)

}

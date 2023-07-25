package v1_iris

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	artemis_api_requests "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/api_requests"
	iris_round_robin "github.com/zeus-fyi/olympus/pkg/iris/proxy/round_robin"
)

const (
	QuickNodeTestHeader = "X-QN-TESTING"
	QuickNodeIDHeader   = "x-quicknode-id"
	QuickNodeEndpointID = "x-instance-id"
	QuickNodeChain      = "x-qn-chain"
	QuickNodeNetwork    = "x-qn-network"
)

func RpcLoadBalancerRequestHandler(c echo.Context) error {
	request := new(ProxyRequest)
	request.Body = echo.Map{}
	if err := c.Bind(&request.Body); err != nil {
		log.Err(err)
		return err
	}
	return request.ProcessRpcLoadBalancerRequest(c)
}

func (p *ProxyRequest) ProcessRpcLoadBalancerRequest(c echo.Context) error {
	routeGroup := c.QueryParam("routeGroup")
	if routeGroup == "" {
		return c.JSON(http.StatusBadRequest, Response{Message: "routeGroup is required"})
	}
	ou := c.Get("orgUser").(org_users.OrgUser)
	routeInfo, err := iris_round_robin.GetNextRoute(ou.OrgID, routeGroup)
	if err != nil {
		log.Err(err).Interface("ou", ou).Str("routeGroup", routeGroup).Msg("ProcessRpcLoadBalancerRequest: iris_round_robin.GetNextRoute")
		errResp := Response{Message: "routeGroup not found"}
		return c.JSON(http.StatusBadRequest, errResp)
	}
	req := &artemis_api_requests.ApiProxyRequest{
		Url:        routeInfo,
		Payload:    p.Body,
		IsInternal: false,
		Timeout:    2 * time.Minute,
	}
	rw := artemis_api_requests.NewArtemisApiRequestsActivities()
	resp, err := rw.ExtLoadBalancerRequest(c.Request().Context(), req)
	if err != nil {
		log.Err(err).Interface("ou", ou).Str("route", routeInfo).Msg("ProcessRpcLoadBalancerRequest: rw.ExtLoadBalancerRequest")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, resp.Response)
}

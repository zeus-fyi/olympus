package v1Beta_iris

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	iris_api_requests "github.com/zeus-fyi/olympus/pkg/iris/proxy/orchestrations/api_requests"
	iris_round_robin "github.com/zeus-fyi/olympus/pkg/iris/proxy/round_robin"
)

func InternalRoundRobinRequestHandler(c echo.Context) error {
	request := new(BetaProxyRequest)
	request.Body = echo.Map{}
	if err := c.Bind(&request.Body); err != nil {
		log.Err(err)
		return err
	}
	return request.ProcessRoundRobin(c, true)
}
func (p *BetaProxyRequest) ProcessRoundRobin(c echo.Context, isInternal bool) error {
	rw := iris_api_requests.NewIrisApiRequestsActivities()
	routeGroup := c.QueryParam("routeGroup")
	ou := c.Get("orgUser").(org_users.OrgUser)
	routeInfo, err := iris_round_robin.GetNextRoute(ou.OrgID, routeGroup)
	if err != nil {
		log.Err(err).Msg("iris_round_robin.GetNextRoute")
		return c.JSON(http.StatusBadRequest, err)
	}
	req := &iris_api_requests.ApiProxyRequest{
		Url:        routeInfo,
		Payload:    p.Body,
		IsInternal: isInternal,
		Timeout:    2 * time.Minute,
	}
	resp, err := rw.InternalSvcRelayRequest(c.Request().Context(), req)
	if err != nil {
		log.Err(err).Str("route", routeInfo).Msg("rw.InternalSvcRelayRequest")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, resp.Response)
}

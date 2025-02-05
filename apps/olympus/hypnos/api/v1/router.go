package v1_hypnos

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	v1_iris "github.com/zeus-fyi/olympus/iris/api/v1"
	iris_api_requests "github.com/zeus-fyi/olympus/pkg/iris/proxy/orchestrations/api_requests"
	iris_usage_meters "github.com/zeus-fyi/olympus/pkg/iris/proxy/usage_meters"
)

func Routes(e *echo.Echo) *echo.Echo {
	// Routes
	e.GET("/health", Health)
	e.POST("/node", RpcLoadBalancerRequestHandlerNode("POST"))
	e.POST("/", RpcLoadBalancerRequestHandler("POST"))
	return e
}

var (
	sessionID  = ""
	routeTable = ""
)

const (
	AnvilSessionLockHeader = "X-Anvil-Session-Lock-ID"
	RouteGroupHeader       = "X-Route-Group"
	NodeURL                = "http://localhost:8545"
	InternalRouter         = "http://iris.iris.svc.cluster.local/v3/internal/router"
	LocalProxiedRouter     = "http://localhost:8888/node"
)

func RpcLoadBalancerRequestHandler(method string) func(c echo.Context) error {
	return func(c echo.Context) error {
		bodyBytes, err := io.ReadAll(c.Request().Body)
		if err != nil {
			log.Err(err)
			return err
		}

		tableHeader := c.Request().Header.Get(RouteGroupHeader)
		anvilHeader := c.Request().Header.Get(AnvilSessionLockHeader)
		if len(anvilHeader) <= 0 {
			log.Info().Interface("anvilHeader", anvilHeader).Msgf("Hypnos: RpcLoadBalancerRequestHandler: anvil session lock id is required")
			return c.JSON(http.StatusBadRequest, v1_iris.Response{Message: "anvil session lock id is required"})
		}
		// todo revist after pods deletion is confirmed
		sessionID = anvilHeader

		payloadSizingMeter := iris_usage_meters.NewPayloadSizeMeter(bodyBytes)
		request := new(v1_iris.ProxyRequest)
		request.Body = echo.Map{}

		if err = json.NewDecoder(payloadSizingMeter).Decode(&request.Body); err != nil {
			log.Err(err).Msgf("Hypnos: RpcLoadBalancerRequestHandler: json.NewDecoder.Decode")
			return c.JSON(http.StatusBadRequest, nil)
		}
		reqHeaders := http.Header{}

		if tableHeader != "" {
			routeTable = tableHeader
		}

		reqHeaders.Add("Authorization", "Bearer "+sessionID)
		rw := iris_api_requests.NewIrisApiRequestsActivities()
		req := &iris_api_requests.ApiProxyRequest{
			Url:              NodeURL,
			Payload:          request.Body,
			PayloadTypeREST:  method,
			RequestHeaders:   reqHeaders,
			IsInternal:       true,
			Timeout:          1 * time.Minute,
			StatusCode:       http.StatusOK, // default
			PayloadSizeMeter: payloadSizingMeter,
		}
		resp, err := rw.ExtLoadBalancerRequest(context.Background(), req)
		if err != nil {
			log.Err(err).Interface("body", request.Body).Msgf("Hypnos: RpcLoadBalancerRequestHandler: rw.ExtLoadBalancerRequest")
			return c.JSON(resp.StatusCode, resp.Response)
		}
		return c.JSON(resp.StatusCode, resp.Response)
	}
}

func RpcLoadBalancerRequestHandlerNode(method string) func(c echo.Context) error {
	return func(c echo.Context) error {
		bodyBytes, err := io.ReadAll(c.Request().Body)
		if err != nil {
			log.Err(err)
			return err
		}

		payloadSizingMeter := iris_usage_meters.NewPayloadSizeMeter(bodyBytes)
		request := new(v1_iris.ProxyRequest)
		request.Body = echo.Map{}

		if err = json.NewDecoder(payloadSizingMeter).Decode(&request.Body); err != nil {
			log.Err(err).Interface("req", request.Body).Msgf("Hypnos: RpcLoadBalancerRequestHandler: json.NewDecoder.Decode")
			return c.JSON(http.StatusBadRequest, nil)
		}
		rw := iris_api_requests.NewIrisApiRequestsActivities()
		req := &iris_api_requests.ApiProxyRequest{
			Url:              InternalRouter,
			Payload:          request.Body,
			PayloadTypeREST:  method,
			IsInternal:       true,
			Timeout:          1 * time.Minute,
			StatusCode:       http.StatusOK, // default
			PayloadSizeMeter: payloadSizingMeter,
		}
		req.RequestHeaders = http.Header{}
		if routeTable != "" {
			req.RequestHeaders.Add("X-Route-Group", routeTable)
		}
		if sessionID != "" {
			req.RequestHeaders.Add("Authorization", "Bearer "+sessionID)
		}
		log.Info().Msgf("Hypnos: RpcLoadBalancerRequestHandler: req: %+v", sessionID)
		resp, err := rw.ExtLoadBalancerRequest(context.Background(), req)
		if err != nil {
			log.Err(err).Interface("body", resp.Response).Msgf("Hypnos: RpcLoadBalancerRequestHandler: rw.ExtLoadBalancerRequest")
			return c.JSON(resp.StatusCode, resp.Response)
		}
		log.Info().Msgf("Hypnos: RpcLoadBalancerRequestHandler: rw.ExtLoadBalancerRequest: resp: %+v", resp)
		return c.JSON(resp.StatusCode, resp.Response)
	}
}

func Health(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy")
}

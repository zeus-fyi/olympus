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
	web3_actions "github.com/zeus-fyi/zeus/pkg/artemis/web3/client"
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
)

func RpcLoadBalancerRequestHandler(method string) func(c echo.Context) error {
	return func(c echo.Context) error {
		bodyBytes, err := io.ReadAll(c.Request().Body)
		if err != nil {
			log.Err(err)
			return err
		}

		tableHeader := c.Request().Header.Get(RouteGroupHeader)
		routeTable = tableHeader

		anvilHeader := c.Request().Header.Get(AnvilSessionLockHeader)
		if len(anvilHeader) <= 0 {
			return c.JSON(http.StatusBadRequest, v1_iris.Response{Message: "anvil session lock id is required"})
		}

		if sessionID == "" {
			wa := web3_actions.NewWeb3ActionsClient(NodeURL)
			wa.IsAnvilNode = true
			err = wa.ResetNetwork(context.Background(), NodeURL, 0)
			if err != nil {
				log.Err(err).Msgf("Hypnos: RpcLoadBalancerRequestHandler: wa.ResetNetwork")
				return err
			}
		}
		sessionID = anvilHeader
		payloadSizingMeter := iris_usage_meters.NewPayloadSizeMeter(bodyBytes)
		request := new(v1_iris.ProxyRequest)
		request.Body = echo.Map{}

		if err = json.NewDecoder(payloadSizingMeter).Decode(&request.Body); err != nil {
			log.Err(err).Msgf("Hypnos: RpcLoadBalancerRequestHandler: json.NewDecoder.Decode")
			return err
		}
		rw := iris_api_requests.NewIrisApiRequestsActivities()
		req := &iris_api_requests.ApiProxyRequest{
			Url:              NodeURL,
			Payload:          request.Body,
			PayloadTypeREST:  method,
			IsInternal:       true,
			Timeout:          1 * time.Minute,
			StatusCode:       http.StatusOK, // default
			PayloadSizeMeter: payloadSizingMeter,
		}

		resp, err := rw.ExtLoadBalancerRequest(context.Background(), req)
		if err != nil {
			log.Err(err).Msgf("Hypnos: RpcLoadBalancerRequestHandler: rw.ExtLoadBalancerRequest")
			return err
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
			log.Err(err).Msgf("Hypnos: RpcLoadBalancerRequestHandler: json.NewDecoder.Decode")
			return err
		}
		irisBetaSvc := "http://iris.iris.svc.cluster.local/v3/internal/router"

		rw := iris_api_requests.NewIrisApiRequestsActivities()
		req := &iris_api_requests.ApiProxyRequest{
			Url:              irisBetaSvc,
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
		resp, err := rw.ExtLoadBalancerRequest(context.Background(), req)
		if err != nil {
			log.Err(err).Msgf("Hypnos: RpcLoadBalancerRequestHandler: rw.ExtLoadBalancerRequest")
			return err
		}
		log.Info().Msgf("Hypnos: RpcLoadBalancerRequestHandler: rw.ExtLoadBalancerRequest: resp: %+v", resp)
		return c.JSON(resp.StatusCode, resp.Response)
	}
}

func Health(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy")
}

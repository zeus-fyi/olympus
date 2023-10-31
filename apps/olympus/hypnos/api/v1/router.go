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

	e.POST("/resp", RpcLoadBalancerRequestHandler("POST"))

	e.POST("/", RpcLoadBalancerRequestHandler2("POST"))
	return e
}

func RpcLoadBalancerRequestHandler(method string) func(c echo.Context) error {
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
		rw := iris_api_requests.NewIrisApiRequestsActivities()
		req := &iris_api_requests.ApiProxyRequest{
			Url:              "http://localhost:8545",
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

func RpcLoadBalancerRequestHandler2(method string) func(c echo.Context) error {
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
		irisBetaSvc := "http://iris.iris.svc.cluster.local/v1/internal/router"
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
		resp, err := rw.ExtLoadBalancerRequest(context.Background(), req)
		if err != nil {
			log.Err(err).Msgf("Hypnos: RpcLoadBalancerRequestHandler: rw.ExtLoadBalancerRequest")
			return err
		}
		return c.JSON(resp.StatusCode, resp.Response)
	}
}
func Health(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy")
}

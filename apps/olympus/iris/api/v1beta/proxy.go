package v1Beta_iris

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	iris_api_requests "github.com/zeus-fyi/olympus/pkg/iris/proxy/orchestrations/api_requests"
)

type ProxyRequest struct {
	Body echo.Map
}

func InternalProxyRequestHandler(c echo.Context) error {
	request := new(ProxyRequest)
	request.Body = echo.Map{}
	if err := c.Bind(&request.Body); err != nil {
		log.Err(err)
		return err
	}
	return request.ProcessInternal(c)
}

func ExternalProxyRequestHandler(c echo.Context) error {
	request := new(ProxyRequest)
	request.Body = echo.Map{}
	if err := c.Bind(&request.Body); err != nil {
		log.Err(err)
		return err
	}
	return request.ProcessExternal(c)
}

func (p *ProxyRequest) ProcessExternal(c echo.Context) error {
	return p.Process(c, false)
}

func (p *ProxyRequest) ProcessInternal(c echo.Context) error {
	return p.Process(c, true)
}

func (p *ProxyRequest) Process(c echo.Context, isInternal bool) error {
	// TODO set queue
	// use ttl
	// tx count per second
	// service-context-id, service-url,
	// FIFO -> process requests until service-context-id ttl expires from last request
	// other requests are queued by service-context-id
	// servCtx := c.Request().Header.Get("Service-Context-ID")
	// servCtx := c.Request().Header.Get("Service-Host")
	// probably needs a unique id to serve back to the client
	// start of session
	// queue -> dynamodb
	// end of session
	relayTo := c.Request().Header.Get("Proxy-Relay-To")
	if relayTo == "" {
		return c.JSON(http.StatusBadRequest, errors.New("Proxy-Relay-To header is required"))
	}
	urlVal, err := url.ParseRequestURI(relayTo)
	if err != nil {
		return c.JSON(http.StatusNotAcceptable, err)
	}
	apiPr := &iris_api_requests.ApiProxyRequest{
		Url:        urlVal.String(),
		Payload:    p.Body,
		IsInternal: isInternal,
		Timeout:    1 * time.Minute,
	}
	ctx := context.Background()
	resp, err := iris_api_requests.IrisProxyWorker.ExecuteIrisProxyWorkflow(ctx, apiPr)
	if err != nil {
		log.Err(err)
		return err
	}
	return c.JSON(http.StatusOK, resp.Response)
}

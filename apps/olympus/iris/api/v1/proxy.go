package v1_iris

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type ProxyRequest struct {
	Body echo.Map
}

func ProxyRequestHandler(c echo.Context) error {
	request := new(ProxyRequest)
	request.Body = echo.Map{}
	if err := c.Bind(&request.Body); err != nil {
		log.Err(err)
		return err
	}
	return request.Process(c)
}

func (p *ProxyRequest) Process(c echo.Context) error {
	// TODO set queue
	// use ttl
	// tx count per second
	// service-context-id, service-url,
	// FIFO -> process requests until service-context-id ttl expires from last request
	// other requests are queued by service-context-id
	// servCtx := c.Request().Header.Get("Service-Context-ID")
	// probably needs a unique id to serve back to the client
	// start of session
	// queue -> dynamodb
	// end of session
	return nil
}

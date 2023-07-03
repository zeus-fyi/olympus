package iris_proxy

import "github.com/labstack/echo/v4"

type ProxyRequest struct {
	Body echo.Map
}

func NewProxyRequest() ProxyRequest {
	return ProxyRequest{}
}

func (p *ProxyRequest) SetSessionLock() {
}

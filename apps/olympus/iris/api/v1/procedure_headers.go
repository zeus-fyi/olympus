package v1_iris

import (
	"github.com/labstack/echo/v4"
	iris_programmable_proxy_v1_beta "github.com/zeus-fyi/zeus/zeus/iris_programmable_proxy/v1beta"
)

func (p *ProxyRequest) ExtractProcedureIfExists(c echo.Context) string {
	procName := c.Request().Header.Get(iris_programmable_proxy_v1_beta.RequestHeaderRoutingProcedureHeader)
	if procName != "" {
		return procName
	}
	procedure := p.Body["procedure"]
	if procedure != nil {
		procNameStr, ok := procedure.(string)
		if ok {
			return procNameStr
		}
	}
	return ""
}

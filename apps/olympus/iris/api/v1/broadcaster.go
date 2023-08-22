package v1_iris

import (
	"github.com/labstack/echo/v4"
	iris_usage_meters "github.com/zeus-fyi/olympus/pkg/iris/proxy/usage_meters"
)

func (p *ProxyRequest) ProcessBroadcastETLRequest(c echo.Context, payloadSizingMeter *iris_usage_meters.PayloadSizeMeter, restType, metricName string) error {
	return nil
}

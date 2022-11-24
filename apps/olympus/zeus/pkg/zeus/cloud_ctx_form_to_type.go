package zeus

import (
	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/pkg/zeus/core/zeus_common_types"
)

func ConvertCloudCtxNsFormToType(c echo.Context) zeus_common_types.CloudCtxNs {
	cloudCtx := zeus_common_types.NewCloudCtxNs()
	cloudCtx.CloudProvider = c.FormValue("cloudProvider")
	cloudCtx.Region = c.FormValue("region")
	cloudCtx.Context = c.FormValue("context")
	cloudCtx.Namespace = c.FormValue("namespace")
	cloudCtx.Env = c.FormValue("env")
	return cloudCtx
}

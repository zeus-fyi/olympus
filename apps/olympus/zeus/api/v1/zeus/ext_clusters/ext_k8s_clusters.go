package zeus_v1_clusters_api

import (
	"github.com/labstack/echo/v4"
)

func ExternalDeployRoutes(e *echo.Group) *echo.Group {
	e.POST("/kubeconfig", CreateOrUpdateKubeConfigsHandler)
	return nil
}

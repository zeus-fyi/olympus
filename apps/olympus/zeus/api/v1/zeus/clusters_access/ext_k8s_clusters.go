package zeus_v1_clusters_api

import (
	"github.com/labstack/echo/v4"
)

func ClusterAccessRoutes(e *echo.Group) *echo.Group {
	e.GET("/clusters/all", ReadAuthorizedClustersRequestHandler)
	e.GET("/clusters/private", ReadExtKubeConfigsHandler)
	e.PUT("/clusters/private", UpdateExtClustersRequestHandler)
	e.POST("/kubeconfig", CreateOrUpdateKubeConfigsHandler)
	return nil
}

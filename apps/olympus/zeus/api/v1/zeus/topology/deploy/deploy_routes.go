package deploy_routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/create_or_update"
	delete_deploy "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/delete"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus/core"
)

type DeploymentActionRequest struct {
	Action string
	create_or_update_deploy.TopologyDeployCreateActionDeployRequest
	delete_deploy.TopologyDeployActionDeleteDeploymentRequest
}

func HandleDeploymentActionRequest(c echo.Context) error {
	request := new(DeploymentActionRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	if request.Action == "create" {
		return request.DeployTopology(c)
	}
	if request.Action == "delete" {
		return request.DeleteDeployedTopology(c)
	}
	return c.JSON(http.StatusBadRequest, nil)
}

func Routes(e *echo.Echo, k8Cfg autok8s_core.K8Util) *echo.Echo {
	core.K8Util = k8Cfg
	e.POST("/deploy", HandleDeploymentActionRequest)

	return e
}

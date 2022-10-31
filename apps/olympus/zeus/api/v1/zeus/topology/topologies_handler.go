package topology

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/base"
)

type TopologyRequest struct {
}

func HandleTopologyActionRequest(c echo.Context) error {
	request := new(base.TopologyActionRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	//if request.Action == "create" {
	//	return request.CreateTopology(c, request)
	//}
	//if request.Action == "read" {
	//	return request.ReadTopology(c, request)
	//}
	//if request.Action == "update" {
	//	return request.UpdateTopology(c, request)
	//}
	//if request.Action == "deploy" {
	//	return request.DeployTopology(c, request)
	//}
	//if request.Action == "delete-deploy" {
	//	return request.DeleteDeployedTopology(c, request)
	//}
	return c.JSON(http.StatusBadRequest, nil)
}

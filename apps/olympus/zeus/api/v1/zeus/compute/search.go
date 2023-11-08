package zeus_v1_compute_api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	hestia_compute_resources "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/resources"
	"github.com/zeus-fyi/zeus/zeus/cluster_resources/nodes"
)

type NodeSearchRequest struct {
	nodes.NodeSearchParams `json:"nodeSearchParams"`
}

func NodeSearchHandler(c echo.Context) error {
	request := new(NodeSearchRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.SearchNodes(c)
}

func (r *NodeSearchRequest) SearchNodes(c echo.Context) error {
	nodesSlice, err := hestia_compute_resources.SearchAndSelectNodes(c.Request().Context(), r.NodeSearchParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, nodesSlice)
}

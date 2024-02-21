package v1_zeus_clusters

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/cloud_ctx_logs"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

type ClusterLogsRequest struct {
	CloudCtxNs zeus_common_types.CloudCtxNs `json:"cloudCtxNs,omitempty"`
}

func ClusterLogsRequestHandler(c echo.Context) error {
	request := new(ClusterLogsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.GetLogsRequest(c)
}

func (r *ClusterLogsRequest) GetLogsRequest(c echo.Context) error {
	request := c.Get("CloudCtxNsLogs").(cloud_ctx_logs.CloudCtxNsLogs)
	logs, err := cloud_ctx_logs.SelectCloudCtxNsLogs(c.Request().Context(), request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, strings.Join(logs, "\n"))
}

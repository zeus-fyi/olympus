package external_api_workloads

import "github.com/labstack/echo/v4"

func ExternalApiWorkloadQueryActionRequestHandler(c echo.Context) error {
	request := new(TopologyCloudCtxNsQueryRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.ReadDeployedWorkloads(c)
}

package read_infra

import "github.com/labstack/echo/v4"

func ReadTopologyChartContentsHandler(c echo.Context) error {
	request := new(TopologyReadRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.ReadTopologyChart(c)
}

func ReadTopologiesOrgCloudCtxNsHandler(c echo.Context) error {
	request := new(TopologyReadRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.ReadTopologiesOrgCloudCtxNs(c)
}

func ReadOrgAppsHandler(c echo.Context) error {
	request := new(TopologyReadPrivateAppsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.ListPrivateAppsRequest(c)
}

func ReadOrgAppDetailsHandler(c echo.Context) error {
	request := new(TopologyReadPrivateAppsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.GetAppDetailsRequest(c)
}

func ReadClusterAppViewOrgCloudCtxNsHandler(c echo.Context) error {
	request := new(TopologyReadRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.ReadClusterAppViewOrgCloudCtxNs(c)
}

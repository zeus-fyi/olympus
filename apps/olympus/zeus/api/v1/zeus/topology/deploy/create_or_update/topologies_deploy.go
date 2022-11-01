package create_or_update_deploy

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	read_charts "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/charts"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/base"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
)

type TopologyDeployCreateActionDeployRequest struct {
	base.TopologyActionRequest
}

func (t *TopologyDeployCreateActionDeployRequest) DeployTopology(c echo.Context) error {
	chartReader := read_charts.Chart{}
	chartReader.ChartPackageID = 6831980425944305799

	ctx := context.Background()
	q := sql_query_templates.QueryParams{}
	q.CTEQuery.Params = append(q.CTEQuery.Params, chartReader.ChartPackageID)
	err := chartReader.SelectSingleChartsResources(ctx, q)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	err = DeployChartPackage(ctx, test.Kns, chartReader)

	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	return err
}

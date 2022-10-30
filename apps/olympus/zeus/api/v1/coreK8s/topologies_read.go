package coreK8s

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	read_charts "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/charts"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

// TODO should read the topology id
func (t *TopologyActionRequest) ReadTopology(c echo.Context, request *TopologyActionRequest) error {
	//chart := t.GetInfraChartPackage()

	chartReader := read_charts.Chart{}
	chartReader.ChartPackageID = 6759098198199705229
	//chart.ChartPackageID

	ctx := context.Background()
	q := sql_query_templates.QueryParams{}
	err := chartReader.SelectSingleChartsResources(ctx, q)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	return c.JSON(http.StatusOK, chartReader)
}

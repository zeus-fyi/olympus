package update_infra

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	read_charts "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/charts"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	base_infra "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/base"
)

type TopologyActionUpdateRequest struct {
	base_infra.TopologyInfraActionRequest
}

func (t *TopologyActionUpdateRequest) UpdateTopology(c echo.Context) error {
	//chart := t.GetInfraChartPackage()

	chartReader := read_charts.Chart{}
	chartReader.ChartPackageID = 6831980425944305799
	//chart.ChartPackageID

	ctx := context.Background()
	q := sql_query_templates.QueryParams{}
	err := chartReader.SelectSingleChartsResources(ctx, q)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	return c.JSON(http.StatusOK, chartReader)
}

package create

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	read_charts "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/charts"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/base"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy"
)

type TopologyActionCreateRequest struct {
	base.TopologyActionRequest
}

// TODO add components to package
// TODO should read the topology id
func (t *TopologyActionCreateRequest) CreateTopology(c echo.Context, request *base.TopologyActionRequest) error {
	//chart := t.GetInfraChartPackage()

	chartReader := read_charts.Chart{}
	chartReader.ChartPackageID = 6828704980826292343
	//chart.ChartPackageID

	ctx := context.Background()
	q := sql_query_templates.QueryParams{}
	err := chartReader.SelectSingleChartsResources(ctx, q)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	err = deploy.DeployChartPackage(ctx, t.Kns, chartReader)

	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	return err
}

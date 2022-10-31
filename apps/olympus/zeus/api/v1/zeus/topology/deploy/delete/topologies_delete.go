package delete_deploy

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	read_charts "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/charts"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/base"
)

type TopologyDeployActionDeleteDeploymentRequest struct {
	base.TopologyActionRequest
}

func (t *TopologyDeployActionDeleteDeploymentRequest) DeleteDeployedTopology(c echo.Context) error {
	// TODO add components to package
	// TODO should read the topology id

	chartReader := read_charts.Chart{}
	chartReader.ChartPackageID = 6831980425944305799
	//chart.ChartPackageID

	ctx := context.Background()
	q := sql_query_templates.QueryParams{}
	err := chartReader.SelectSingleChartsResources(ctx, q)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	err = DeleteK8sWorkload(ctx, t.Kns, chartReader)

	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	return err
}

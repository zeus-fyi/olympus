package delete

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	read_charts "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/charts"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/base"
	delete2 "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy_delete"
)

type TopologyActionDeleteRequest struct {
	base.TopologyActionRequest
}

func (t *TopologyActionDeleteRequest) DeleteDeployedTopology(c echo.Context, request *base.TopologyActionRequest) error {
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
	err = delete2.DeleteK8sWorkload(ctx, t.Kns, chartReader)

	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	return err
}

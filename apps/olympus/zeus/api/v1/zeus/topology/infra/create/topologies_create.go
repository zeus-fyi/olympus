package create_infra

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/packages"
	create_infra "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/definitions/classes/bases/infra"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/base"
)

type TopologyActionCreateRequest struct {
	base.TopologyActionRequest
	TopologyCreateRequest
}

type TopologyCreateRequest struct {
	Name string `json:"name"`

	charts.Chart
	chart_workload.NativeK8s
}

func (t *TopologyActionCreateRequest) CreateTopology(c echo.Context) error {
	pkg := packages.NewPackageInsert()
	cw, err := t.NativeK8s.CreateChartWorkloadFromNativeK8s()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.New("unable to parse chart"))
	}
	pkg.ChartWorkload = cw
	pkg.Chart = t.Chart
	topCreate := create_infra.NewOrgUserCreateInfraFromNativeK8s(t.Name, t.OrgUser, pkg)

	ctx := context.Background()
	err = topCreate.InsertInfraBase(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, nil)
}

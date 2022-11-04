package create_infra

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	create_infra "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/definitions/classes/bases/infra"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/base"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

type TopologyActionCreateRequest struct {
	base.TopologyActionRequest
	TopologyCreateRequest
}

type TopologyCreateRequest struct {
	Name string `json:"name"`
	charts.Chart
}

type TopologyCreateResponse struct {
	TopologyName     string `json:"topologyName"`
	ID               int    `json:"id"`
	ChartName        string `json:"chartName"`
	ChartDescription string `json:"chartDescription"`
	Version          string `json:"version"`
}

func (t *TopologyActionCreateRequest) CreateTopology(c echo.Context) error {
	file, err := c.FormFile("chart")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	in := bytes.Buffer{}
	if _, err = io.Copy(&in, src); err != nil {
		return err
	}
	nk, err := zeus.UnGzipK8sChart(&in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	cw, err := nk.CreateChartWorkloadFromNativeK8s()
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	inf := create_infra.NewCreateInfrastructure()
	ctx := context.Background()
	inf.Name = t.Name
	inf.ChartWorkload = cw
	inf.Chart = t.Chart
	inf.OrgUser = t.OrgUser
	err = inf.InsertInfraBase(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	resp := TopologyCreateResponse{
		TopologyName:     t.Name,
		ID:               inf.TopologyID,
		ChartName:        inf.ChartName,
		ChartDescription: inf.ChartDescription.String,
		Version:          inf.ChartVersion,
	}
	return c.JSON(http.StatusOK, resp)
}

package delete_deploy

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/base"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
)

type TopologyDeployActionDeleteDeploymentRequest struct {
	base.TopologyActionRequest
	TopologyID int
}

func (t *TopologyDeployActionDeleteDeploymentRequest) DeleteDeployedTopology(c echo.Context) error {
	tr := read_topology.NewInfraTopologyReader()
	tr.TopologyID = t.TopologyID
	tr.OrgID = t.OrgID
	tr.UserID = t.UserID

	ctx := context.Background()
	err := tr.SelectTopology(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	nk := tr.GetNativeK8s()

	// TODO should also fetch deployed kns, then update it
	kns := test.Kns
	err = DeleteK8sWorkload(ctx, kns, nk)

	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	return err
}

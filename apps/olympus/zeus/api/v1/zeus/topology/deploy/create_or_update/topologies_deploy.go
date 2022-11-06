package create_or_update_deploy

import (
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/base"
)

type TopologyDeployCreateActionDeployRequest struct {
	// TODO gets kns
	base.TopologyActionRequest
	TopologyID int
}

//func (t *TopologyDeployCreateActionDeployRequest) DeployTopology(c echo.Context) error {
//	tr := read_topology.NewInfraTopologyReader()
//	tr.TopologyID = t.TopologyID
//	tr.OrgID = t.OrgID
//	tr.UserID = t.UserID
//
//	ctx := context.Background()
//	err := tr.SelectTopology(ctx)
//	if err != nil {
//		return c.JSON(http.StatusInternalServerError, err)
//	}
//	nk := tr.GetNativeK8s()
//
//	// TODO should also write kns & deployed status status
//	kns := test.Kns
//	kns.Namespace = "zeus"
//	kns.CtxType = "dev-do-sfo3-zeus"
//	err = DeployChartPackage(ctx, kns, nk)
//
//	if err != nil {
//		return c.JSON(http.StatusBadRequest, nil)
//	}
//	return err
//}

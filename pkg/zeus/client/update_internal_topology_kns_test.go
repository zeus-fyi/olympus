package zeus_client

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/state"
)

func (t *ZeusClientTestSuite) TestUpdateInternalTopologyKnsStatus() {
	status := topology_deployment_status.Status{
		TopologyKubeCtxNs: kns.TopologyKubeCtxNs{
			TopologyID: deployKnsReq.TopologyID,
			CloudCtxNs: deployKnsReq.CloudCtxNs,
		},
		DeployStatus: topology_deployment_status.DeployStatus{},
	}
	resp, err := t.ZeusTestClient.UpdateTopologyKnsStatus(ctx, status)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}

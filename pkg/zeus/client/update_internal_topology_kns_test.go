package zeus_client

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/state"
)

func (t *ZeusClientTestSuite) TestUpdateInternalTopologyKnsStatus() {
	k := kns.TopologyKubeCtxNs{
		TopologyID: deployKnsReq.TopologyID,
		CloudCtxNs: deployKnsReq.CloudCtxNs,
	}
	status := topology_deployment_status.NewPopulatedTopologyStatus(k, "Pending")
	resp, err := t.ZeusTestClient.UpdateTopologyKnsStatus(ctx, status)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}

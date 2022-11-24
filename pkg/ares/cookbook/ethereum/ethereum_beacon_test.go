package ethereum_beacon_cookbook

import "github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"

// set your own topologyID here after uploading a chart workload
var deployConsensusKnsReq = zeus_req_types.TopologyDeployRequest{
	TopologyID: 1669168938743775000,
	CloudCtxNs: topCloudCtxNs,
}

// set your own topologyID here after uploading a chart workload
var deployExecKnsReq = zeus_req_types.TopologyDeployRequest{
	TopologyID: 1669168971777999000,
	CloudCtxNs: topCloudCtxNs,
}

func (t *ZeusEthereumBeaconTestSuite) TestDeployBeacon() {
	resp, err := t.ZeusTestClient.Deploy(ctx, deployConsensusKnsReq)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)

	resp, err = t.ZeusTestClient.Deploy(ctx, deployExecKnsReq)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}

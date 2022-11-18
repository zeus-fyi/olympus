package zeus_client

import "github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"

func (t *ZeusClientTestSuite) TestReadLiveNamespaceWorkload() {
	tar := zeus_req_types.TopologyCloudCtxNsQueryRequest{
		CloudCtxNs: topCloudCtxNs,
	}
	resp, err := t.ZeusTestClient.ReadNamespaceWorkload(ctx, tar)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}

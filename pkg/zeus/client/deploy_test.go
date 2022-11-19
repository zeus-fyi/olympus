package zeus_client

import (
	"errors"
)

// TestChartUpload will return a topology id associated with this workload
func (t *ZeusClientTestSuite) TestDeploy() {
	resp, err := t.ZeusTestClient.Deploy(ctx, deployKnsReq)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}

func (t *ZeusClientTestSuite) TestDeployWithoutTopologyIDFails() {
	deployKnsReq.TopologyID = 0
	resp, err := t.ZeusTestClient.Deploy(ctx, deployKnsReq)
	t.Require().NotNil(err)
	t.Assert().Equal(errors.New("bad request"), err)
	t.Assert().Empty(resp)
}

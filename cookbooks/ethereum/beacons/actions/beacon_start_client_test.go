package beacon_actions

import (
	client_consts "github.com/zeus-fyi/olympus/cookbooks/ethereum/beacons/constants"
)

func (t *BeaconActionsTestSuite) TestStartConsensusClient() {
	t.ConsensusClient = client_consts.Lighthouse
	resp, err := t.StartConsensusClient(ctx)
	t.Assert().Nil(err)
	t.Assert().NotEmpty(resp)
}

func (t *BeaconActionsTestSuite) TestStartExecClient() {
	t.ExecClient = client_consts.Geth
	resp, err := t.StartExecClient(ctx)
	t.Assert().Nil(err)
	t.Assert().NotEmpty(resp)
}

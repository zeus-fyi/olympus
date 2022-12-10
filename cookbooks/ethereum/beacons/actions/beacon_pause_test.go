package beacon_actions

func (t *BeaconActionsTestSuite) TestPauseConsensusClient() {
	t.ConsensusClient = "lighthouse"
	_, err := t.PauseClient(ctx, "cm-lighthouse", t.ConsensusClient)
	t.Assert().Nil(err)
}

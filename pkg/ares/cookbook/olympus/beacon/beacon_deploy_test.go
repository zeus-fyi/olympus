package olympus_beacon

func (t *ZeusAppsTestSuite) TestDeployBeacon() {
	resp, err := t.ZeusTestClient.Deploy(ctx, deployConsensusKnsReq)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)

	resp, err = t.ZeusTestClient.Deploy(ctx, deployExecKnsReq)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)

	t.TestSecretsCopy()
}

func (t *ZeusAppsTestSuite) TestDestroyBeacon() {
	resp, err := t.ZeusTestClient.DestroyDeploy(ctx, deployConsensusKnsReq)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}

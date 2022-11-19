package zeus_client

// TestDeployReplace will replace the components at this location, but does not change the underlying topology
// definitions. In other words, this is a localized change.
func (t *ZeusClientTestSuite) TestDeployReplace() {
	resp, err := t.ZeusTestClient.DeployReplace(ctx, replaceTopologyComponentPath, deployKnsReq)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}

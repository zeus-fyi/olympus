package zeus_client

// TestDeployReplace will replace the components at this location, but does not change the underlying topology
// definitions. In other words, this is a localized change.
func (t *ZeusClientTestSuite) TestCleanDeployedNamespace() {
	resp, err := t.ZeusTestClient.CleanDeployedNamespace(ctx, deployKnsReq)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}

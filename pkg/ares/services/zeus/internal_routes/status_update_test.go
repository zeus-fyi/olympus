package internal_routes

func (t *AresZeusInternalRoutesTestSuite) TestUpdateStatus() {
	resp, err := t.ZeusTestClient.UpdateTopologyStatus(ctx, status)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}

func (t *AresZeusInternalRoutesTestSuite) TestUpdateKnsStatus() {
	resp, err := t.ZeusTestClient.UpdateTopologyKnsStatus(ctx, status)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}

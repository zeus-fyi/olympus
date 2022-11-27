package zeus_actions

func (t *ZeusCookbookActionsTestSuite) TestRestart() {
	r, err := t.ZeusActionsClient.RestartZeusPods(ctx)
	t.Require().Nil(err)
	t.Assert().NotEmpty(r)
}

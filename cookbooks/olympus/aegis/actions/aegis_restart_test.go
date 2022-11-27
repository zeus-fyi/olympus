package aegis_actions

func (t *AegisCookbookActionsTestSuite) TestRestart() {
	r, err := t.AegisActionsClient.RestartAegisPods(ctx)
	t.Require().Nil(err)
	t.Assert().NotEmpty(r)
}

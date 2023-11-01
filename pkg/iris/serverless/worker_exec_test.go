package iris_serverless

func (t *IrisOrchestrationsTestSuite) TestUpdateResetTimer() {
	a := NewIrisPlatformActivities()
	err := a.ResyncServerlessRoutes(ctx, nil)
	t.Require().NoError(err)
}

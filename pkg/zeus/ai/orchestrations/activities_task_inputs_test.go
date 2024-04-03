package ai_platform_service_orchestrations

func (t *ZeusWorkerTestSuite) TestAiAggregateAnalysisRetrievalTask() {
	na := NewZeusAiPlatformActivities()
	_, err := na.AiAggregateAnalysisRetrievalTask(ctx, nil, nil)
	t.Require().Nil(err)
}

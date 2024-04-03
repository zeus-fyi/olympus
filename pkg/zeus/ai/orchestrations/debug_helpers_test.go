package ai_platform_service_orchestrations

func (t *ZeusWorkerTestSuite) TestAiAggregateAnalysisRetrievalTaskInputDebug() {
	db := AiAggregateAnalysisRetrievalTaskInputDebug{}
	db.Cp = &MbChildSubProcessParams{}
	db.Cp.Wsr.RunCycle = 0
	db.Cp.Wsr.ChunkOffset = 0
	db.Open()
}

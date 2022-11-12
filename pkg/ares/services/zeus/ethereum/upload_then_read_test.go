package ethereum

func (t *AresZeusTestSuite) TestIntegratedUploadThenRead() {
	resp := t.TestCreateAndUploadConsensusClientChart()
	t.readUploadedConsensusChart(resp.ID)
}

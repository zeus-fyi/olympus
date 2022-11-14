package zeus_ethereum

func (t *AresZeusEthereumTestSuite) TestConsensusClientIntegratedUploadThenRead() {
	resp := t.TestCreateAndUploadConsensusClientChart()
	t.readUploadedConsensusChart(resp.ID)
}

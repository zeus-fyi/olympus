package zeus_ethereum

func (t *AresZeusEthereumTestSuite) TestConsensusClientIntegratedUploadThenRead() {
	resp := t.TestCreateAndUploadConsensusClientChart()
	t.readUploadedConsensusChart(resp.ID)
}

func (t *AresZeusEthereumTestSuite) TestExecClientIntegratedUploadThenRead() {
	resp := t.TestCreateAndUploadExecClientChart()
	t.readUploadedExecChart(resp.ID)
}

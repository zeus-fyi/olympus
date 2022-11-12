package zeus_ethereum

func (t *AresZeusEthereumTestSuite) TestIntegratedUploadThenRead() {
	resp := t.TestCreateAndUploadConsensusClientChart()
	t.readUploadedConsensusChart(resp.ID)
}

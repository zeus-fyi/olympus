package artemis_reporting

import "fmt"

func (s *ReportingTestSuite) TestCalculateGasCosts() {
	bg, err := GetBundleSubmissionHistory(ctx, 0, 1)
	s.Assert().Nil(err)
	s.Assert().NotNil(bg)

	for eid, b := range bg.Map {
		fmt.Println("eventID", eid)
		for _, bundleTx := range b {
			fmt.Println("bundleTx.EthTx.TxHash", bundleTx.EthTx.TxHash, "bundle", bundleTx.EthTxGas.GasTipCap, "bundle", bundleTx.EthTxGas.GasFeeCap, "bundleTx.EthTxGas.GasPrice", bundleTx.EthTxGas.GasPrice, "bundleTx.EthTxGas.GasLimit", bundleTx.EthTxGas.GasLimit)
		}
	}
}

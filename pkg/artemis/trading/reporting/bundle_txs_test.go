package artemis_reporting

import "fmt"

// TODO:
// 1. needs to get tx receipts and calculate gas costs for each tx in the bundle
// 2. needs to add schema & save rx receipts to the db
// 3. needs to compare total gas costs to the profit of the bundle
// 4. needs to save bundle call responses to the db

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

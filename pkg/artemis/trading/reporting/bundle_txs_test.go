package artemis_reporting

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// TODO:
// 1. needs to get tx receipts and calculate gas costs for each tx in the bundle
// 2. needs to add schema & save rx receipts to the db
// 3. needs to compare total gas costs to the profit of the bundle
// 4. needs to save bundle call responses to the db
// 5. needs to setup fb eth_bundle rpc stat lookup
const AccountAddr = "0x000000641e80A183c8B736141cbE313E136bc8c6"

func (s *ReportingTestSuite) TestCalculateGasCosts() {
	bg, err := GetBundleSubmissionHistory(ctx, 0, 1)
	s.Assert().Nil(err)
	s.Assert().NotNil(bg)

	for eid, b := range bg.Map {
		fmt.Println("eventID", eid)
		for _, bundleTx := range b {
			if bundleTx.From == AccountAddr {
				fmt.Println("bundleTx.EthTx.TxHash", bundleTx.EthTx.TxHash, "bundle", bundleTx.EthTxGas.GasTipCap, "bundle", bundleTx.EthTxGas.GasFeeCap, "bundleTx.EthTxGas.GasPrice", bundleTx.EthTxGas.GasPrice, "bundleTx.EthTxGas.GasLimit", bundleTx.EthTxGas.GasLimit)
			}
			rx, found, rerr := s.w3c.GetTxReceipt(ctx, common.HexToHash(bundleTx.EthTx.TxHash))
			s.Assert().Nil(rerr)
			if !found {
				continue
			}
			fmt.Println(rx.BlockNumber)
			fmt.Println(rx.GasUsed)
			fmt.Println(rx.CumulativeGasUsed)
			fmt.Println(rx.Status)
			fmt.Println(rx.Logs)
			fmt.Println(rx.ContractAddress)
			fmt.Println(rx.TxHash)
			fmt.Println(rx.TransactionIndex)
			fmt.Println(rx.BlockHash)
			fmt.Println(rx.PostState)
		}
	}
}

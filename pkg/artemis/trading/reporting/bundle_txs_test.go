package artemis_reporting

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	artemis_eth_rxs "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/txs/eth_rxs"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (s *ReportingTestSuite) TestBundleHistoryFetch() {
	bg, err := GetBundleSubmissionHistory(ctx, 0, 1)
	s.Assert().Nil(err)
	s.Assert().NotNil(bg)

	for bundleHash, b := range bg.Map {
		s.Assert().NotEmpty(bundleHash)
		for _, bundleTx := range b {
			if bundleTx.EthTx.From == AccountAddr {

			}
			s.Assert().NotEmpty(bundleTx.EthTx.TxHash)
		}
	}
}

// effective_gas_price = priority_fee_per_gas + block.base_fee_per_gas

func (s *ReportingTestSuite) TestCalculateGasFees() {
	bg, err := GetBundleSubmissionHistory(ctx, 0, 1)
	s.Assert().Nil(err)
	s.Assert().NotNil(bg)
	// 18422122804364250 vs 0.01254829 + 0.00587382
	for bundleHash, b := range bg.Map {
		fees := 0
		rxBlockNumber := 0
		predictedRevenue := 0
		for _, bundleTx := range b {
			if bundleTx.EthTx.From == AccountAddr {
				fmt.Println("bundleTx.EthTx.TxHash", bundleTx.EthTx.TxHash, "EthTxReceipts.EffectiveGasPrice", bundleTx.EffectiveGasPrice,
					"bundleTx.EthTxGas.GasTipCap", bundleTx.EthTxGas.GasTipCap, "bundleTx.EthTxGas.GasFeeCap", bundleTx.EthTxGas.GasFeeCap,
					"bundleTx.EthTxGas.GasLimit", bundleTx.EthTxGas.GasLimit)
				fees += bundleTx.EthTxReceipts.GasUsed * bundleTx.EthTxReceipts.EffectiveGasPrice
			} else {
				mevMempoolTx, merr := artemis_mev_models.SelectEthMevMempoolTxByTxHash(ctx, bundleTx.EthTx.TxHash)
				s.Require().NoError(merr)
				s.Require().Len(mevMempoolTx, 1)
				tx := mevMempoolTx[0]
				j, merr := web3_client.UnmarshalTradeExecutionFlow(tx.TxFlowPrediction)
				s.Require().NoError(merr)
				fmt.Println(j.SandwichTrade.AmountOut)
				predictedRevenue = int(artemis_eth_units.NewBigIntFromStr(j.SandwichPrediction.ExpectedProfit).Int64())
			}
			rxBlockNumber = bundleTx.EthTxReceipts.BlockNumber
		}
		if fees == 0 {
			continue
		}
		fmt.Println("rxBlockNumber", rxBlockNumber)
		//0.117129433
		wethBalChange, werr := s.w3c.GetMainnetBalanceDiffWETH(AccountAddr, rxBlockNumber)
		s.Assert().Nil(werr)
		fmt.Println("wethBalChange", wethBalChange)
		err = InsertBundleProfit(ctx, artemis_autogen_bases.EthMevBundleProfit{
			BundleHash:        bundleHash,
			Revenue:           int(wethBalChange.Int64()),
			RevenuePrediction: predictedRevenue,
			Costs:             fees,
		})
		s.Assert().Nil(err)
	}
}

func (s *ReportingTestSuite) TestInsertRxsForEthTxs() {
	bg, err := GetBundleSubmissionHistory(ctx, 0, 1)
	s.Assert().Nil(err)
	s.Assert().NotNil(bg)

	for eid, b := range bg.Map {
		fmt.Println("eventID", eid)
		for _, bundleTx := range b {
			if bundleTx.EthTx.From == AccountAddr {
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
			fmt.Println(rx.EffectiveGasPrice.String())

			status := "unknown"
			if rx.Status == types.ReceiptStatusSuccessful {
				status = "success"
			}
			if rx.Status == types.ReceiptStatusFailed {
				status = "failed"
			}
			rxEthTx := artemis_autogen_bases.EthTxReceipts{
				Status:            status,
				GasUsed:           int(rx.GasUsed),
				CumulativeGasUsed: int(rx.CumulativeGasUsed),
				BlockHash:         rx.BlockHash.String(),
				TransactionIndex:  int(rx.TransactionIndex),
				TxHash:            rx.TxHash.String(),
				EventID:           bundleTx.EthTx.EventID,
				EffectiveGasPrice: int(rx.EffectiveGasPrice.Int64()),
				BlockNumber:       int(rx.BlockNumber.Int64()),
			}
			err = artemis_eth_rxs.InsertTxReceipt(ctx, rxEthTx)
			s.Assert().Nil(err)
		}
	}
}

/*
	assert signer.balance >= transaction.gas_limit * transaction.max_fee_per_gas

	# ensure that the user was willing to at least pay the base fee
	assert transaction.max_fee_per_gas >= block.base_fee_per_gas

	# Prevent impossibly large numbers
	assert transaction.max_fee_per_gas < 2**256
	# Prevent impossibly large numbers
	assert transaction.max_priority_fee_per_gas < 2**256
	# The total must be the larger of the two
	assert transaction.max_fee_per_gas >= transaction.max_priority_fee_per_gas

	# priority fee is capped because the base fee is filled first
	priority_fee_per_gas = min(transaction.max_priority_fee_per_gas, transaction.max_fee_per_gas - block.base_fee_per_gas)

	# signer pays both the priority fee and the base fee
	effective_gas_price = priority_fee_per_gas + block.base_fee_per_gas
	signer.balance -= transaction.gas_limit * effective_gas_price
	assert signer.balance >= 0, 'invalid transaction: signer does not have enough ETH to cover gas'
	gas_used = self.execute_transaction(transaction, effective_gas_price)
	gas_refund = transaction.gas_limit - gas_used
	cumulative_transaction_gas_used += gas_used

	# signer gets refunded for unused gas
	signer.balance += gas_refund * effective_gas_price

	# miner only receives the priority fee; note that the base fee is not given to anyone (it is burned)
	self.account(block.author).balance += gas_used * priority_fee_per_gas
*/

package artemis_realtime_trading

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	metrics_trading "github.com/zeus-fyi/olympus/pkg/apollo/ethereum/mev/trading"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	artemis_oly_contract_abis "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/contract_abis"
)

var (
	OldUniversalRouterAbi = artemis_oly_contract_abis.MustLoadOldUniversalRouterAbi()
	NewUniversalRouterAbi = artemis_oly_contract_abis.MustLoadNewUniversalRouterAbi()
	UniswapV2Router02Abi  = artemis_oly_contract_abis.MustLoadUniswapV2Router02ABI()
	UniswapV2Router01Abi  = artemis_oly_contract_abis.MustLoadUniswapV2Router01ABI()
	UniswapV3Router01Abi  = artemis_oly_contract_abis.MustLoadUniswapV3Swap1RouterAbi()
	UniswapV3Router02Abi  = artemis_oly_contract_abis.MustLoadUniswapV3Swap2RouterAbi()
)

func DecodeTx(ctx context.Context, tx *types.Transaction, m *metrics_trading.TradingMetrics) (*web3_client.MevTx, error) {
	var mevTx web3_client.MevTx
	switch tx.To().String() {
	case web3_client.UniswapUniversalRouterAddressOld:
		methodName, args, err := web3_client.DecodeTxArgDataFromAbi(ctx, tx, OldUniversalRouterAbi)
		if err != nil {
			log.Err(err).Msg("UniswapUniversalRouterAddressOld: error decoding tx arg data from abi")
			return nil, err
		}
		singleTx := web3_client.MevTx{
			MethodName: methodName,
			Args:       args,
			Tx:         tx,
		}
		mevTx = singleTx
		if m != nil {
			m.TxFetcherMetrics.TransactionGroup(web3_client.UniswapUniversalRouterAddressOld, methodName)
		}
	case web3_client.UniswapUniversalRouterAddressNew:
		methodName, args, err := web3_client.DecodeTxArgDataFromAbi(ctx, tx, NewUniversalRouterAbi)
		if err != nil {
			log.Err(err).Msg("UniswapUniversalRouterAddressNew: error decoding tx arg data from abi")
			return nil, err
		}
		singleTx := web3_client.MevTx{
			MethodName: methodName,
			Args:       args,
			Tx:         tx,
		}
		mevTx = singleTx
		if m != nil {
			m.TxFetcherMetrics.TransactionGroup(web3_client.UniswapUniversalRouterAddressNew, methodName)
		}
	case web3_client.UniswapV2Router02Address:
		methodName, args, err := web3_client.DecodeTxArgDataFromAbi(ctx, tx, UniswapV2Router02Abi)
		if err != nil {
			log.Err(err).Msg("UniswapV2Router02Address: error decoding tx arg data from abi")
			return nil, err
		}
		singleTx := web3_client.MevTx{
			MethodName: methodName,
			Args:       args,
			Tx:         tx,
		}
		mevTx = singleTx
		if m != nil {
			m.TxFetcherMetrics.TransactionGroup(web3_client.UniswapV2Router02Address, methodName)
		}
	case web3_client.UniswapV2Router01Address:
		methodName, args, err := web3_client.DecodeTxArgDataFromAbi(ctx, tx, UniswapV2Router01Abi)
		if err != nil {
			log.Err(err).Msg("UniswapV2Router01Address: error decoding tx arg data from abi")
			return nil, err
		}
		singleTx := web3_client.MevTx{
			MethodName: methodName,
			Args:       args,
			Tx:         tx,
		}
		mevTx = singleTx
		if m != nil {
			m.TxFetcherMetrics.TransactionGroup(web3_client.UniswapV2Router01Address, methodName)
		}
	case web3_client.UniswapV3Router01Address:
		methodName, args, err := web3_client.DecodeTxArgDataFromAbi(ctx, tx, UniswapV3Router01Abi)
		if err != nil {
			return nil, err
		}
		singleTx := web3_client.MevTx{
			MethodName: methodName,
			Args:       args,
			Tx:         tx,
		}
		mevTx = singleTx
		if m != nil {
			m.TxFetcherMetrics.TransactionGroup(web3_client.UniswapV3Router01Address, methodName)
		}
	case web3_client.UniswapV3Router02Address:
		methodName, args, err := web3_client.DecodeTxArgDataFromAbi(ctx, tx, UniswapV3Router02Abi)
		if err != nil {
			log.Err(err).Msg("UniswapV3Router02Address: error decoding tx arg data from abi")
			return nil, err
		}
		singleTx := web3_client.MevTx{
			MethodName: methodName,
			Args:       args,
			Tx:         tx,
		}
		mevTx = singleTx
		if m != nil {
			m.TxFetcherMetrics.TransactionGroup(web3_client.UniswapV3Router02Address, methodName)
		}
	}

	return &mevTx, nil
}

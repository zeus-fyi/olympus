package artemis_realtime_trading

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	metrics_trading "github.com/zeus-fyi/olympus/pkg/apollo/ethereum/mev/trading"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	artemis_oly_contract_abis "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/contract_abis"
)

func ProcessTxs(ctx context.Context, mevTx web3_client.MevTx, m *metrics_trading.TradingMetrics, w3a web3_actions.Web3Actions) ([]web3_client.TradeExecutionFlow, error) {
	switch mevTx.Tx.To().String() {
	case artemis_trading_constants.UniswapUniversalRouterAddressOld:
		tf, err := RealTimeProcessUniversalRouterTx(ctx, mevTx, m, w3a, artemis_oly_contract_abis.UniversalRouterDecoder)
		if err != nil {
			log.Err(err).Msg("UniswapUniversalRouterAddressOld: error processing universal router tx")
			return nil, err
		}
		return tf, nil
	case artemis_trading_constants.UniswapUniversalRouterAddressNew:
		tf, err := RealTimeProcessUniversalRouterTx(ctx, mevTx, m, w3a, artemis_oly_contract_abis.UniversalRouterDecoder)
		if err != nil {
			log.Err(err).Msg("UniswapUniversalRouterAddressNew: error processing universal router tx")
			return nil, err
		}
		return tf, nil
	case artemis_trading_constants.UniswapV2Router01Address:
		tf, err := RealTimeProcessUniswapV2RouterTx(ctx, mevTx, m, w3a, artemis_oly_contract_abis.UniswapV2Router01)
		if err != nil {
			log.Err(err).Msg("UniswapV2Router01Address: error processing v2_01 router tx")
			return nil, err
		}
		return tf, nil
	case artemis_trading_constants.UniswapV2Router02Address:
		tf, err := RealTimeProcessUniswapV2RouterTx(ctx, mevTx, m, w3a, artemis_oly_contract_abis.UniswapV2Router02)
		if err != nil {
			log.Err(err).Msg("UniswapV2Router02Address: error processing v2_02 router tx")
			return nil, err
		}
		return tf, nil
	case artemis_trading_constants.UniswapV3Router01Address:
		tf, err := RealTimeProcessUniswapV3RouterTx(ctx, mevTx, UniswapV3Router01Abi, nil, m, w3a)
		if err != nil {
			log.Err(err).Msg("UniswapV3Router01Address: error processing v3_01 router tx")
			return nil, err
		}
		return tf, nil
	case artemis_trading_constants.UniswapV3Router02Address:
		tf, err := RealTimeProcessUniswapV3RouterTx(ctx, mevTx, UniswapV3Router02Abi, nil, m, w3a)
		if err != nil {
			log.Err(err).Msg("UniswapV3Router02Address: error processing v3_02 router tx")
			return nil, err
		}
		return tf, nil
	}
	log.Warn().Msgf("ProcessTxs: tx.To() not recognized: %s", mevTx.Tx.To().String())
	return nil, errors.New("ProcessTxs: tx.To() not recognized")
}

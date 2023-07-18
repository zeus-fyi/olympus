package artemis_trade_debugger

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

func (t *TradeDebugger) FindSlippage(ctx context.Context, to *artemis_trading_types.TradeOutcome) error {
	ac := t.dat.GetSimAuxClient()
	//tf.FrontRunTrade.AmountOut = tf.FrontRunTrade.SimulatedAmountOut //  new(big.Int).SetInt64(0)
	ur, err := ac.GenerateTradeV2SwapFromTokenToToken(ctx, nil, to)
	if err != nil {
		return err
	}
	start := to.AmountOut
	num := 0
	denom := 1000
	for i := 1; i < 7; i++ {
		switch i {
		case 0:
			num = 1
			denom = 1
		case 1:
			num = 10
			denom = 1000
		case 2:
			num = 20
			denom = 1000
		case 3:
			num = 75
			denom = 1000
		case 4:
			num = 110
			denom = 1000
		case 5:
			num = 210
			denom = 1000
		default:
			return errors.New("failed to find a valid transfer tax")
		}
		to.AmountOut = artemis_eth_units.ApplyTransferTax(start, num, denom)
		fmt.Println("amount out", to.AmountOut.String())
		ur, err = ac.GenerateTradeV2SwapFromTokenToToken(ctx, nil, to)
		if err != nil {
			return err
		}
		err = t.dat.GetSimUniswapClient().InjectExecTradeV2SwapFromTokenToToken(ctx, ur, to)
		if err == nil {
			log.Info().Interface("num", num).Msgf("Injected trade with amount out: %s", to.AmountOut.String())
			break
		}
	}
	if num == 1000 {
		num = 1
		denom = 1
	}
	if to.AmountOutAddr.String() != artemis_trading_constants.WETH9ContractAddress && to.AmountInAddr.String() == artemis_trading_constants.WETH9ContractAddress {
		err = artemis_mev_models.UpdateERC20TokenTransferTaxInfo(ctx, artemis_autogen_bases.Erc20TokenInfo{
			Address:                to.AmountOutAddr.String(),
			ProtocolNetworkID:      hestia_req_types.EthereumMainnetProtocolNetworkID,
			TransferTaxNumerator:   &num,
			TransferTaxDenominator: &denom,
		})
		if err != nil {
			return err
		}
	}
	if to.AmountOutAddr.String() == artemis_trading_constants.WETH9ContractAddress && to.AmountInAddr.String() != artemis_trading_constants.WETH9ContractAddress {
		err = artemis_mev_models.UpdateERC20TokenTransferTaxInfo(ctx, artemis_autogen_bases.Erc20TokenInfo{
			Address:                to.AmountInAddr.String(),
			ProtocolNetworkID:      hestia_req_types.EthereumMainnetProtocolNetworkID,
			TransferTaxNumerator:   &num,
			TransferTaxDenominator: &denom,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

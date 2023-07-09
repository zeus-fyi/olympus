package artemis_trading_auxiliary

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (a *AuxiliaryTradingUtils) GenerateTradeV2SwapFromTokenToToken(ctx context.Context, ur *web3_client.UniversalRouterExecCmd, to *artemis_trading_types.TradeOutcome) (*web3_client.UniversalRouterExecCmd, error) {
	ur = a.checkIfCmdEmpty(ur)
	sc1 := web3_client.UniversalRouterExecSubCmd{
		Command:   artemis_trading_constants.Permit2Permit,
		CanRevert: false,
		Inputs:    nil,
	}
	psp, err := a.generatePermit2Approval(ctx, to)
	if err != nil {
		log.Err(err).Msg("failed to generate permit2 approval")
		return nil, err
	}
	sc1.DecodedInputs = psp
	ur.Commands = append(ur.Commands, sc1)
	sc2 := web3_client.UniversalRouterExecSubCmd{
		Command:   artemis_trading_constants.V2SwapExactIn,
		CanRevert: false,
		Inputs:    nil,
		DecodedInputs: web3_client.V2SwapExactInParams{
			AmountIn:      to.AmountIn,
			AmountOutMin:  to.AmountOut,
			Path:          []accounts.Address{to.AmountInAddr, to.AmountOutAddr},
			To:            artemis_trading_constants.UniversalRouterSenderAddress,
			PayerIsSender: true,
		},
	}
	ur.Commands = append(ur.Commands, sc2)
	return ur, err
}

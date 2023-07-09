package artemis_trading_auxiliary

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (a *AuxiliaryTradingUtils) GenerateTradeV2SwapFromTokenToToken(ctx context.Context, ur *web3_client.UniversalRouterExecCmd, to *artemis_trading_types.TradeOutcome) (*web3_client.UniversalRouterExecCmd, error) {
	ur = a.checkIfCmdEmpty(ur)
	deadline := a.GetDeadline()
	sc1 := web3_client.UniversalRouterExecSubCmd{
		Command:   artemis_trading_constants.Permit2Permit,
		CanRevert: false,
		Inputs:    nil,
	}
	psp := web3_client.Permit2PermitParams{
		PermitSingle: web3_client.PermitSingle{
			PermitDetails: web3_client.PermitDetails{
				Token:      to.AmountInAddr,
				Amount:     to.AmountIn,
				Expiration: deadline,
				Nonce:      a.GetPermit2Nonce(),
			},
			Spender:     artemis_trading_constants.UniswapUniversalRouterNewAddressAccount,
			SigDeadline: deadline,
		},
	}
	err := psp.SignPermit2Mainnet(a.Account)
	if err != nil {
		log.Warn().Err(err).Msg("error signing permit")
		return nil, err
	}
	if psp.Signature == nil {
		log.Warn().Msg("signature is nil")
		return nil, errors.New("signature is nil")
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

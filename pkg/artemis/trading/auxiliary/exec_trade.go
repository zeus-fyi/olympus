package artemis_trading_auxiliary

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_eth_txs "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/txs/eth_txs"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func GenerateTradeV2SwapFromTokenToToken(ctx context.Context, w3c web3_client.Web3Client, ur *web3_client.UniversalRouterExecCmd, to *artemis_trading_types.TradeOutcome) (*web3_client.UniversalRouterExecCmd, *artemis_eth_txs.Permit2Tx, error) {
	ur = checkIfCmdEmpty(ur)
	if w3c.Account == nil {
		return nil, nil, errors.New("GenerateTradeV2SwapFromTokenToToken: account is nil")
	}
	sc1 := web3_client.UniversalRouterExecSubCmd{
		Command:   artemis_trading_constants.Permit2Permit,
		CanRevert: false,
		Inputs:    nil,
	}
	psp, pt, err := generatePermit2Approval(ctx, w3c, to)
	if err != nil {
		log.Warn().Str("permit2Owner", w3c.Account.PublicKey()).Msg("GenerateTradeV2SwapFromTokenToToken: generatePermit2Approval failed")
		log.Err(err).Msg("failed to generate permit2 approval")
		return nil, nil, err
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
	return ur, pt, err
}

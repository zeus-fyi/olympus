package artemis_trading_auxiliary

import (
	"context"
	"errors"
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (a *AuxiliaryTradingUtils) GenerateCmdToExchangeETHtoWETH(ctx context.Context, ur *web3_client.UniversalRouterExecCmd, amountIn *big.Int, user *accounts.Address) (*web3_client.UniversalRouterExecCmd, error) {
	ur = a.checkIfCmdEmpty(ur)
	if a.Account == nil && user == nil {
		return nil, errors.New("no account or user address provided")
	}
	if user == nil {
		addr := a.Account.Address()
		user = &addr
	}
	wethParams := web3_client.WrapETHParams{
		Recipient: *user,
		AmountMin: amountIn,
	}
	payable := &web3_actions.SendEtherPayload{
		TransferArgs: web3_actions.TransferArgs{
			Amount:    amountIn,
			ToAddress: wethParams.Recipient,
		},
		GasPriceLimits: web3_actions.GasPriceLimits{},
	}
	ur.Commands = append(ur.Commands, web3_client.UniversalRouterExecSubCmd{
		Command:       artemis_trading_constants.WrapETH,
		DecodedInputs: wethParams,
		CanRevert:     false,
	})
	if ur.Payable == nil {
		ur.Payable = payable
	} else {
		ur.Payable.Amount = artemis_eth_units.AddBigInt(ur.Payable.Amount, amountIn)
		ur.Payable.ToAddress = wethParams.Recipient
	}
	return ur, nil
}

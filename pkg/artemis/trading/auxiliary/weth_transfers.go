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
		//addr := artemis_trading_constants.UniversalRouterSenderAddress
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
			ToAddress: artemis_trading_constants.UniswapUniversalRouterNewAddressAccount,
		},
		GasPriceLimits: web3_actions.GasPriceLimits{},
	}
	ur.Commands = append(ur.Commands, web3_client.UniversalRouterExecSubCmd{
		Command:       artemis_trading_constants.WrapETH,
		DecodedInputs: wethParams,
		CanRevert:     false,
	})
	if ur.Payable.ToAddress.String() != artemis_trading_constants.ZeroAddress && ur.Payable.ToAddress.String() != wethParams.Recipient.String() {
		return nil, errors.New("payable amount and address mismatch")
	}
	ur.Payable.ToAddress = wethParams.Recipient
	if ur.Payable != nil && ur.Payable.Amount != nil {
		ur.Payable.Amount = artemis_eth_units.AddBigInt(ur.Payable.Amount, amountIn)
	} else {
		ur.Payable = payable
	}
	return ur, nil
}

func (a *AuxiliaryTradingUtils) GenerateCmdToExchangeWETHtoETH(ctx context.Context, ur *web3_client.UniversalRouterExecCmd, amountIn *big.Int, user *accounts.Address) (*web3_client.UniversalRouterExecCmd, error) {
	ur = a.checkIfCmdEmpty(ur)
	if a.Account == nil && user == nil {
		return nil, errors.New("no account or user address provided")
	}
	if user == nil {
		addr := a.Account.Address()
		user = &addr
	}
	unwrapParams := web3_client.UnwrapWETHParams{
		Recipient: *user,
		AmountMin: amountIn,
	}
	ur.Commands = append(ur.Commands, web3_client.UniversalRouterExecSubCmd{
		Command:       artemis_trading_constants.UnwrapWETH,
		DecodedInputs: unwrapParams,
		CanRevert:     false,
	})
	return ur, nil
}

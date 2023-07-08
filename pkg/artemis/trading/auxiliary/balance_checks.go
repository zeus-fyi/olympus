package artemis_trading_auxiliary

import "C"
import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	core_entities "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/uniswap_libs/uniswap_core/entities"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

func (a *AuxiliaryTradingUtils) universalRouterCmdVerifier(ctx context.Context, ur *web3_client.UniversalRouterExecCmd, signedTx *types.Transaction, scInfo *web3_actions.SendContractTxPayload) error {
	ethRequirements := artemis_eth_units.NewBigInt(0)
	for _, sc := range ur.Commands {
		switch sc.Command {
		case artemis_trading_constants.WrapETH:
			wp := sc.DecodedInputs.(web3_client.WrapETHParams)
			ethRequirements = artemis_eth_units.AddBigInt(ethRequirements, wp.AmountMin)
			if ur.Payable == nil {
				return errors.New("payable is nil")
			}
		case artemis_trading_constants.UnwrapWETH:
		}
	}
	// todo: add gas cost using scInfo
	// scInfo
	gasCost := artemis_eth_units.NewBigInt(0)
	ethRequirements = artemis_eth_units.AddBigInt(ethRequirements, gasCost)
	hasEnough, err := a.checkAuxEthBalanceGreaterThan(ctx, ethRequirements)
	if err != nil {
		return err
	}
	if !hasEnough {
		return errors.New("user does not have enough ETH to exchange to WETH")
	}

	est, err := a.C.SuggestGasPrice(ctx)
	if err != nil {
		return err
	}
	gt, err := a.C.SuggestGasTipCap(ctx)
	if err != nil {
		return err
	}

	toAddr := common.HexToAddress(artemis_trading_constants.UniswapUniversalRouterAddressNew)
	msg := ethereum.CallMsg{From: common.HexToAddress(a.Address().Hex()), To: &toAddr, Data: signedTx.Data()}
	gasLimit, err := a.C.EstimateGas(ctx, msg)
	if err != nil {
		return err
	}
	fmt.Println("estimated gas limit", gasLimit)
	fmt.Println("estimated gas price", est.String())
	fmt.Println("estimated gas tip cap", gt.String())
	return nil
}

func (a *AuxiliaryTradingUtils) checkAuxEthBalance(ctx context.Context) (*big.Int, error) {
	bal, err := a.GetCurrentBalance(ctx)
	if err != nil {
		return bal, err
	}
	return bal, err
}

func (a *AuxiliaryTradingUtils) checkAuxERC20Balance(ctx context.Context, token core_entities.Token) (*big.Int, error) {
	bal, err := a.ReadERC20TokenBalance(ctx, token.Address.String(), a.Account.Address().String())
	if err != nil {
		return bal, err
	}
	return bal, err
}

func (a *AuxiliaryTradingUtils) checkAuxERC20BalanceGreaterThan(ctx context.Context, token core_entities.Token, amount *big.Int) (bool, error) {
	bal, err := a.checkAuxERC20Balance(ctx, token)
	if err != nil {
		return false, err
	}
	return artemis_eth_units.IsXGreaterThanY(bal, amount), err
}

func (a *AuxiliaryTradingUtils) checkAuxEthBalanceGreaterThan(ctx context.Context, amount *big.Int) (bool, error) {
	bal, err := a.GetCurrentBalance(ctx)
	if err != nil {
		return false, err
	}
	return artemis_eth_units.IsXGreaterThanY(bal, amount), err
}

func (a *AuxiliaryTradingUtils) checkAuxWETHBalance(ctx context.Context) (*big.Int, error) {
	token := core_entities.NewToken(a.getChainID(), artemis_trading_constants.WETH9ContractAddressAccount, 18, "WETH", "Wrapped Ether")
	bal, err := a.ReadERC20TokenBalance(ctx, token.Address.String(), a.Account.Address().String())
	if err != nil {
		return bal, err
	}
	return bal, err
}

func (a *AuxiliaryTradingUtils) checkAuxWETHBalanceGreaterThan(ctx context.Context, amount *big.Int) (bool, error) {
	bal, err := a.checkAuxWETHBalance(ctx)
	if err != nil {
		return false, err
	}
	return artemis_eth_units.IsXGreaterThanY(bal, amount), err
}

func (a *AuxiliaryTradingUtils) getChainID() uint {
	switch a.Network {
	case hestia_req_types.Mainnet:
		return hestia_req_types.EthereumMainnetProtocolNetworkID
	case hestia_req_types.Goerli:
		return hestia_req_types.EthereumGoerliProtocolNetworkID
	case hestia_req_types.Ephemery:
		return hestia_req_types.EthereumEphemeryProtocolNetworkID
	default:
		return hestia_req_types.EthereumMainnetProtocolNetworkID
	}
}

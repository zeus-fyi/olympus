package async_analysis

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	uniswap_core_entities "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/uniswap_libs/uniswap_core/entities"
)

func (c *ContractAnalysis) SimEthTransferFeeTaxTrade(ctx context.Context, amount *big.Int) (*uniswap_core_entities.Percent, error) {
	err := c.UserA.SetERC20BalanceBruteForce(ctx, artemis_trading_constants.WETH9ContractAddressAccount.String(), c.UserA.Address().String(), amount)
	if err != nil {
		return nil, err
	}
	pd, err := c.u.GetV2PricingData(ctx, []accounts.Address{artemis_trading_constants.WETH9ContractAddressAccount, accounts.HexToAddress(c.SmartContractAddr)})
	if err != nil {
		return nil, err
	}
	startBalUserA, err := c.UserA.ReadERC20TokenBalance(ctx, c.SmartContractAddr, c.UserA.Address().String())
	if err != nil {
		return nil, err
	}
	calculatedOut, err := pd.V2Pair.PriceImpact(artemis_trading_constants.WETH9ContractAddressAccount, amount)
	if err != nil {
		return nil, err
	}

	// test amount artemis_eth_units.NewBigInt(0)
	// calculatedOut
	trade := &artemis_trading_types.TradeOutcome{
		AmountIn:      amount,
		AmountInAddr:  artemis_trading_constants.WETH9ContractAddressAccount,
		AmountOut:     calculatedOut.AmountOut,
		AmountOutAddr: accounts.HexToAddress(c.SmartContractAddr),
	}
	err = c.u.ExecTradeV2SwapFromTokenToToken(ctx, trade)
	if err != nil {
		return nil, err
	}

	time.Sleep(3 * time.Second)
	txHash := trade.OrderedTxs[0]
	rx, err := c.UserA.GetTxReceipt(ctx, common.Hash(txHash))
	if err != nil {
		return nil, err
	}
	if rx.Status == types.ReceiptStatusFailed {
		log.Err(err).Interface("token", c.SendContractTxPayload.SmartContractAddr).Msg("tx failed, amount is 0")
		return nil, errors.New("tx failed, amount is 0")
	}
	endBalUserA, err := c.UserA.ReadERC20TokenBalance(ctx, c.SmartContractAddr, c.UserA.Address().String())
	if err != nil {
		return nil, err
	}
	fmt.Println("calculated out", calculatedOut.AmountOut.String())
	actualOut := artemis_eth_units.SubBigInt(endBalUserA, startBalUserA)
	fmt.Println("actual out", actualOut.String())
	drift := artemis_eth_units.SubBigInt(actualOut, calculatedOut.AmountOut)
	fmt.Println("drift", drift.String())

	//fmt.Println("rate", rate.Quotient())
	per := artemis_eth_units.NewPercentFromInts(1, 100)
	tokenFee := artemis_eth_units.FractionalAmount(calculatedOut.AmountOut, per)
	fmt.Println("tokenFee", tokenFee.String())

	per2 := artemis_eth_units.NewPercentFromInts(1, 1000000000000)
	tokenFee2 := artemis_eth_units.FractionalAmount(calculatedOut.AmountOut, per2)
	fmt.Println("tokenFee2", tokenFee2.String())

	fmt.Println("drift after fee", artemis_eth_units.AddBigInt(tokenFee, drift).String())
	rateOut := artemis_eth_units.PercentDiff(calculatedOut.AmountOut, actualOut)
	fmt.Println("rateOut", rateOut.String())
	return nil, nil
}

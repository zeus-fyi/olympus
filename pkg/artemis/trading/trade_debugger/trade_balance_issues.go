package artemis_trade_debugger

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	"github.com/zeus-fyi/olympus/pkg/artemis/trading/async_analysis"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	uniswap_pricing "github.com/zeus-fyi/olympus/pkg/artemis/trading/pricing/uniswap"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	artemis_oly_contract_abis "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/contract_abis"
)

func CheckExpectedReserves(ctx context.Context, w3c web3_client.Web3Client, tf *web3_client.TradeExecutionFlow) error {
	if tf.InitialPair == nil {
		return nil
	}
	// todo, do v3 pairs
	simPair := tf.InitialPair
	err := uniswap_pricing.GetPairContractPrices(ctx, tf.CurrentBlockNumber.Uint64(), w3c.Web3Actions, simPair)
	if err != nil {
		log.Err(err).Msg("error getting pair contract prices")
		return err
	}
	if tf.InitialPair.Reserve1.String() != simPair.Reserve1.String() && tf.InitialPair.Reserve0.String() != simPair.Reserve0.String() {
		fmt.Println("tf.InitialPair.Reserve0", tf.InitialPair.Reserve0.String(), simPair.Reserve0.String(), "simPair.Reserve0")
		fmt.Println("tf.InitialPair.Reserve1", tf.InitialPair.Reserve1.String(), simPair.Reserve1.String(), "simPair.Reserve1")
		return fmt.Errorf("reserve mismatch")
	}
	if tf.InitialPair.Reserve0.String() != simPair.Reserve0.String() {
		fmt.Println("tf.InitialPair.Reserve0", tf.InitialPair.Reserve0.String(), simPair.Reserve0.String(), "simPair.Reserve0")
		return fmt.Errorf("reserve0 mismatch")
	}
	if tf.InitialPair.Reserve1.String() != simPair.Reserve1.String() {
		fmt.Println("tf.InitialPair.Reserve1", tf.InitialPair.Reserve1.String(), simPair.Reserve1.String(), "simPair.Reserve1")
		return fmt.Errorf("reserve1 mismatch")
	}
	return nil
}

func (t *TradeDebugger) analyzeToken(ctx context.Context, address accounts.Address, amountTraded *big.Int) error {
	if address == artemis_trading_constants.WETH9ContractAddressAccount {
		return nil
	}
	token := address.String()
	if _, ok := artemis_trading_cache.TokenMap[token]; !ok {
		return errors.New("token not found")
	}

	info := artemis_trading_cache.TokenMap[token]
	if info.Name != nil {
		fmt.Println("ANALYZING token: ", token, "name: ", *info.Name)
	}
	den := info.TransferTaxDenominator
	num := info.TransferTaxNumerator
	if den != nil && num != nil {
		fmt.Println("token: ", token, "tradingTax: num: ", *num, "den: ", *den)
	} else {
		fmt.Println("token not found in cache")
	}

	ca := async_analysis.NewContractAnalysis(t.dat.GetSimUniswapClient(), address.String(), artemis_oly_contract_abis.MustLoadERC20Abi())
	ca.UserA = t.dat.GetSimUniswapClient().Web3Client
	ca.UserB = t.LiveNetworkClient
	//  -1 means balanceOf it wasn't cracked within 100 slots, -2 means cracking hasn't been attempted yet
	if info.BalanceOfSlotNum == -2 {
		err := ca.FindERC20BalanceOfSlotNumber(ctx)
		if err != nil {
			return err
		}
	}
	if num == nil || den == nil {
		feePerc, err := ca.CalculateTransferFeeTax(ctx, artemis_eth_units.Ether)
		if err != nil {
			return err
		}

		fmt.Println("trade feePerc: ", feePerc.Numerator, feePerc.Denominator)
		calcNum := int(feePerc.Numerator.Int64())
		calcDen := int(feePerc.Denominator.Int64())
		info.TransferTaxNumerator = &calcNum
		info.TransferTaxDenominator = &calcDen
		artemis_trading_cache.TokenMap[token] = info
		err = artemis_mev_models.UpdateERC20TokenTransferTaxInfo(ctx, info)
		if err != nil {
			return err
		}
		//err = ca.CalculateTransferFeeTaxRange(ctx)
		//if err != nil {
		//	return err
		//}
		//amountToTest := artemis_eth_units.EtherMultiple(1)
		//feePerc, err = ca.CalculateTransferFeeTax(ctx, amountToTest)
		//if err != nil {
		//	return err
		//}
		//fmt.Println("feePercCalc: ", feePerc.Numerator, feePerc.Denominator)
		//feePerc, err = ca.SimEthTransferFeeTaxTrade(ctx, amountToTest)
		//if err != nil {
		//	return err
		//}
		//fmt.Println("feePercSim: ", feePerc.Numerator, feePerc.Denominator)
	}

	t.ContractAnalysis = ca
	return nil
}

func (t *TradeDebugger) analyzeDrift(ctx context.Context, trade artemis_trading_types.TradeOutcome) error {
	fmt.Println("ANALYZING DRIFT")
	fmt.Println("trade skew: ", trade.AmountOutDrift.String())
	trade.PrintDebug()
	return errors.New("trade skewed")
}

/*

 */

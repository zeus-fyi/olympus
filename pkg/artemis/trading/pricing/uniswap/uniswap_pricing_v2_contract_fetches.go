package uniswap_pricing

import (
	"context"

	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_pricing_utils "github.com/zeus-fyi/olympus/pkg/artemis/trading/pricing/utils"
	artemis_oly_contract_abis "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/contract_abis"
)

const (
	getReserves = "getReserves"
)

var v2ABI = artemis_oly_contract_abis.MustLoadUniswapV2PairAbi()

func GetPairContractPrices(ctx context.Context, wc web3_actions.Web3Actions, p *UniswapV2Pair) error {
	scInfo := &web3_actions.SendContractTxPayload{
		SmartContractAddr: p.PairContractAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       v2ABI,
	}
	scInfo.MethodName = getReserves
	resp, err := wc.CallConstantFunction(ctx, scInfo)
	if err != nil {
		return err
	}
	if len(resp) <= 2 {
		return err
	}
	reserve0, err := artemis_pricing_utils.ParseBigInt(resp[0])
	if err != nil {
		return err
	}
	p.Reserve0 = reserve0
	reserve1, err := artemis_pricing_utils.ParseBigInt(resp[1])
	if err != nil {
		return err
	}
	p.Reserve1 = reserve1
	blockTimestampLast, err := artemis_pricing_utils.ParseBigInt(resp[2])
	if err != nil {
		return err
	}
	p.BlockTimestampLast = blockTimestampLast
	return nil
}

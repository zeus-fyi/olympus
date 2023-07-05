package uniswap_pricing

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
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

func V2PairToPrices(ctx context.Context, wc web3_actions.Web3Actions, pairAddr []accounts.Address) (UniswapV2Pair, error) {
	p := UniswapV2Pair{}
	if len(pairAddr) == 2 {
		err := p.PairForV2(pairAddr[0].String(), pairAddr[1].String())
		if err != nil {
			log.Err(err).Msg("V2PairToPrices: PairForV2")
			return p, err
		}
		err = GetPairContractPrices(ctx, wc, &p)
		if err != nil {
			log.Err(err).Msg("V2PairToPrices: GetPairContractPrices")
			return p, err
		}
		return p, err
	}
	return UniswapV2Pair{}, errors.New("pair address length is not 2, multi-hops not implemented yet")
}

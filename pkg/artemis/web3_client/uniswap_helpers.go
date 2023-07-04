package web3_client

import (
	"context"
	"errors"
	"math/big"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
	uniswap_pricing "github.com/zeus-fyi/olympus/pkg/artemis/trading/pricing/uniswap"
)

func (u *UniswapClient) SingleReadMethodBigInt(ctx context.Context, methodName string, scInfo *web3_actions.SendContractTxPayload) (*big.Int, error) {
	scInfo.MethodName = methodName
	resp, err := u.Web3Client.CallConstantFunction(ctx, scInfo)
	if err != nil {
		return &big.Int{}, err
	}
	if len(resp) == 0 {
		return &big.Int{}, errors.New("empty response")
	}
	bi, err := ParseBigInt(resp[0])
	if err != nil {
		return &big.Int{}, err
	}
	return bi, nil
}

func (u *UniswapClient) SingleReadMethodAddr(ctx context.Context, methodName string, scInfo *web3_actions.SendContractTxPayload) (accounts.Address, error) {
	scInfo.MethodName = methodName
	resp, err := u.Web3Client.CallConstantFunction(ctx, scInfo)
	if err != nil {
		return accounts.Address{}, err
	}
	if len(resp) == 0 {
		return accounts.Address{}, errors.New("empty response")
	}
	addr, err := ConvertToAddress(resp[0])
	if err != nil {
		return accounts.Address{}, err
	}
	return addr, nil
}

func (u *UniswapClient) V2PairToPrices(ctx context.Context, pairAddr []accounts.Address) (uniswap_pricing.UniswapV2Pair, error) {
	p := uniswap_pricing.UniswapV2Pair{}
	if len(pairAddr) == 2 {
		err := p.PairForV2(pairAddr[0].String(), pairAddr[1].String())
		if err != nil {
			log.Err(err).Msg("V2PairToPrices: PairForV2")
			return p, err
		}
		err = u.GetPairContractPrices(ctx, &p)
		if err != nil {
			log.Err(err).Msg("V2PairToPrices: GetPairContractPrices")
			return p, err
		}
		return p, err
	}
	return uniswap_pricing.UniswapV2Pair{}, errors.New("pair address length is not 2, multi-hops not implemented yet")
}

func (u *UniswapClient) GetPairContractPrices(ctx context.Context, p *uniswap_pricing.UniswapV2Pair) error {
	scInfo := &web3_actions.SendContractTxPayload{
		SmartContractAddr: p.PairContractAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       u.PairAbi,
	}
	scInfo.MethodName = "getReserves"

	wc := u.Web3Client
	if artemis_network_cfgs.ArtemisEthereumMainnetQuiknodeHistoricalPrimary.NodeURL != "" && u.SimMode == false {
		wc = NewWeb3Client(artemis_network_cfgs.ArtemisEthereumMainnetQuiknodeHistoricalPrimary.NodeURL, u.Web3Client.Account)
	}
	resp, err := wc.CallConstantFunction(ctx, scInfo)
	if err != nil {
		return err
	}
	if len(resp) <= 2 {
		return err
	}
	reserve0, err := ParseBigInt(resp[0])
	if err != nil {
		return err
	}
	p.Reserve0 = reserve0
	reserve1, err := ParseBigInt(resp[1])
	if err != nil {
		return err
	}
	p.Reserve1 = reserve1
	blockTimestampLast, err := ParseBigInt(resp[2])
	if err != nil {
		return err
	}
	p.BlockTimestampLast = blockTimestampLast
	return nil
}

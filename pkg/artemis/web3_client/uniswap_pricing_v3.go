package web3_client

import (
	"math/big"

	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_v3/entities"
)

const (
	ticks = "ticks"
	slot0 = "slot0"
)

type UniswapPoolV3 struct {
	web3_actions.Web3Actions
	PoolAddress string
	Slot0       Slot0
}

type Slot0 struct {
	SqrtPriceX96 *big.Int
	Tick         int
	FeeProtocol  int
}

func (p *UniswapPoolV3) GetSlot0() error {
	scInfo := &web3_actions.SendContractTxPayload{
		SmartContractAddr: p.PoolAddress,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       MustLoadPoolV3Abi(),
		MethodName:        slot0,
		Params:            []interface{}{},
	}

	resp, err := p.CallConstantFunction(ctx, scInfo)
	if err != nil {
		return err
	}

	for i, val := range resp {
		switch i {
		case 0:
			p.Slot0.SqrtPriceX96 = val.(*big.Int)
		case 1:
			tmp := val.(*big.Int)
			p.Slot0.Tick = int(tmp.Int64())
		case 7:
			tmp := val.(uint8)
			p.Slot0.FeeProtocol = int(tmp)
		}
	}
	return nil
}
func (p *UniswapPoolV3) GetTick(tickNum int) (entities.Tick, error) {
	scInfo := &web3_actions.SendContractTxPayload{
		SmartContractAddr: p.PoolAddress,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       MustLoadPoolV3Abi(),
		MethodName:        ticks,
		Params:            []interface{}{new(big.Int).SetInt64(int64(tickNum))},
	}
	tick := entities.Tick{
		Index: tickNum,
	}
	resp, err := p.CallConstantFunction(ctx, scInfo)
	if err != nil {
		return tick, err
	}
	for i, val := range resp {
		switch i {
		case 0:
			tick.LiquidityGross = val.(*big.Int)
		case 1:
			tick.LiquidityNet = val.(*big.Int)
		}
	}
	return tick, nil
}

func (p *UniswapPoolV3) NextInitializedTickWithinOneWord(tick int, lte bool, tickSpacing int) (int, bool) {
	return 0, false
}

package web3_client

import (
	"errors"
	"math/big"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_v3/constants"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_v3/entities"
)

const (
	ticks                   = "ticks"
	slot0                   = "slot0"
	liquidity               = "liquidity"
	tickBitmap              = "tickBitmap"
	getPopulatedTicksInWord = "getPopulatedTicksInWord"

	TickLensAddress = "0xbfd8137f7d1516D3ea5cA83523914859ec47F573"
)

type UniswapPoolV3 struct {
	web3_actions.Web3Actions
	PoolAddress          string
	Fee                  constants.FeeAmount
	Slot0                Slot0
	Liquidity            *big.Int
	TickListDataProvider *entities.TickListDataProvider
}

type Slot0 struct {
	SqrtPriceX96 *big.Int
	Tick         int
	FeeProtocol  int
}

func (p *UniswapPoolV3) GetLiquidity() error {
	scInfo := &web3_actions.SendContractTxPayload{
		SmartContractAddr: p.PoolAddress,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       MustLoadPoolV3Abi(),
		MethodName:        liquidity,
		Params:            []interface{}{},
	}
	resp, err := p.CallConstantFunction(ctx, scInfo)
	if err != nil {
		return err
	}
	for i, val := range resp {
		switch i {
		case 0:
			p.Liquidity = val.(*big.Int)
		}
	}
	return err
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
		case 5:
			tmp := val.(uint8)
			p.Slot0.FeeProtocol = int(tmp)
		}
	}
	return nil
}

func (p *UniswapPoolV3) GetTickMappingValue(tickNum int16) *big.Int {
	tick, err := p.GetTickMappingValueFromContract(tickNum)
	if err != nil {
		log.Err(err).Msg("GetTickMappingValueFromContract")
		return nil
	}
	return tick
}

func (p *UniswapPoolV3) GetTickMappingValueFromContract(tickNum int16) (*big.Int, error) {
	scInfo := &web3_actions.SendContractTxPayload{
		SmartContractAddr: p.PoolAddress,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       MustLoadPoolV3Abi(),
		MethodName:        tickBitmap,
		Params:            []interface{}{tickNum},
	}

	resp, err := p.CallConstantFunction(ctx, scInfo)
	if err != nil {
		return nil, err
	}
	for i, val := range resp {
		switch i {
		case 0:
			return val.(*big.Int), nil
		}
	}
	return nil, errors.New("tick mapping value not found")
}

func (p *UniswapPoolV3) GetTick(tickNum int) entities.Tick {
	tick, err := p.GetTickFromContract(tickNum)
	if err != nil {
		return tick
	}
	return tick
}

func (p *UniswapPoolV3) GetTickFromContract(tickNum int) (entities.Tick, error) {
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
			//case 2:
			//	fmt.Println("2", val.(*big.Int))
			//case 3:
			//	fmt.Println("3", val.(*big.Int))
			//case 4:
			//	fmt.Println("4", val.(*big.Int))
			//case 5:
			//	fmt.Println("5", val.(*big.Int))
			//case 6:
			//	fmt.Println("6", val.(uint32))
			//case 7:
			//	fmt.Println("7", val.(bool))
			//case 8:
			//	fmt.Println("8", val.(*big.Int))
		}
	}
	return tick, nil
}

func (p *UniswapPoolV3) GetPopulatedTicksMap() ([]entities.Tick, error) {
	scInfo := &web3_actions.SendContractTxPayload{
		SmartContractAddr: TickLensAddress,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       MustLoadTickLensAbi(),
		MethodName:        getPopulatedTicksInWord,
		Params:            []interface{}{accounts.HexToAddress(p.PoolAddress), GetTickBitmapIndex(new(big.Int).SetInt64(int64(p.Slot0.Tick)), int64(constants.TickSpacings[p.Fee]))},
	}
	var ticksSlice []entities.Tick
	resp, err := p.CallConstantFunction(ctx, scInfo)
	if err != nil {
		return ticksSlice, err
	}
	tmp := resp[0].([]struct {
		Tick           *big.Int "json:\"tick\""
		LiquidityNet   *big.Int "json:\"liquidityNet\""
		LiquidityGross *big.Int "json:\"liquidityGross\""
	})

	ticksSlice = make([]entities.Tick, len(tmp))
	for i, val := range tmp {
		i = len(tmp) - i - 1
		ticksSlice[i].Index = int(val.Tick.Int64())
		ticksSlice[i].LiquidityNet = val.LiquidityNet
		ticksSlice[i].LiquidityGross = val.LiquidityGross
	}
	return ticksSlice, nil
}

func GetTickBitmapIndex(tick *big.Int, tickSpacing int64) int64 {
	tickSpacingBig := big.NewInt(tickSpacing)
	intermediate := new(big.Int).Div(tick, tickSpacingBig)

	two := big.NewInt(2)
	eight := big.NewInt(8)

	if intermediate.Sign() < 0 {
		intermediate.Add(intermediate, big.NewInt(1))
		intermediate.Div(intermediate, new(big.Int).Exp(two, eight, nil))
		intermediate.Sub(intermediate, big.NewInt(1))
	} else {
		intermediate.Rsh(intermediate, uint(eight.Uint64()))
	}

	return intermediate.Int64()
}

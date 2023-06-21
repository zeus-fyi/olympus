package web3_client

import (
	"math"
	"math/big"

	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_v3/entities"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_v3/utils"
)

const (
	ticks     = "ticks"
	slot0     = "slot0"
	liquidity = "liquidity"
)

type UniswapPoolV3 struct {
	web3_actions.Web3Actions
	PoolAddress string
	Slot0       Slot0
	Liquidity   *big.Int
	TickBitmap  TickBitmap
}

type TickBitmap struct {
	Data map[int16]*big.Int
}

/* todo
NextInitializedTickWithinOneWord(tick int, lte bool, tickSpacing int) (int, bool)

	 * Return the next tick that is initialized within a single word
	 * @param tick The current tick
	 * @param lte Whether the next tick should be lte the current tick
	 * @param tickSpacing The tick spacing of the pool

		// because each iteration of the while loop rounds, we can't optimize this code (relative to the smart contract)
		// by simply traversing to the next available tick, we instead need to exactly replicate
		// tickBitmap.nextInitializedTickWithinOneWord
		step.tickNext, step.initialized = p.TickDataProvider.NextInitializedTickWithinOneWord(state.tick, zeroForOne, p.tickSpacing())

  function nextInitializedTickWithinOneWord(
    mapping(int16 => uint256) self,
    int24 tick,
    int24 tickSpacing,
    bool lte
  ) internal view returns (int24 next, bool initialized)

Parameters:
	Name	Type	Description
	self	mapping(int16 => uint256)	The mapping in which to compute the next initialized tick
	tick	int24	The starting tick
	tickSpacing	int24	The spacing between usable ticks
	lte	bool	Whether to search for the next initialized tick to the left (less than or equal to the starting tick)

Return Values:
	Name	Type	Description
	next	int24	The next initialized or uninitialized tick up to 256 ticks away from the current tick
	initialized	bool	Whether the next tick is initialized, as the function only searches within up to 256 ticks

*/

type Slot0 struct {
	SqrtPriceX96 *big.Int
	Tick         int
	FeeProtocol  int
}

func (p *UniswapPoolV3) NextInitializedTickWithinOneWord(tick int, lte bool, tickSpacing int) (int, bool) {
	compressed := tick / tickSpacing
	if tick < 0 && tick%tickSpacing != 0 {
		compressed-- // round towards negative infinity
	}

	if lte {
		wordPos, bitPos := position(int32(compressed))
		// all the 1s at or to the right of the current bitPos
		mask := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), uint(bitPos)), big.NewInt(1))
		mask.Add(mask, new(big.Int).Lsh(big.NewInt(1), uint(bitPos)))
		masked := new(big.Int).And(p.TickBitmap.Data[wordPos], mask)

		// if there are no initialized ticks to the right of or at the current tick, return rightmost in the word
		initialized := masked.Cmp(big.NewInt(0)) != 0
		// overflow/underflow is possible, but prevented externally by limiting both tickSpacing and tick
		next := 0
		if initialized {
			mostSigBit, err := utils.MostSignificantBit(masked)
			if err != nil {
				panic(err)
			}
			next = compressed - (int(bitPos)-int(mostSigBit))*tickSpacing
		} else {
			next = (compressed - int(bitPos)) * tickSpacing
		}

		return next, initialized
	} else {
		// start from the word of the next tick, since the current tick state doesn't matter
		wordPos, bitPos := position(int32(compressed + 1))
		// all the 1s at or to the left of the bitPos
		mask := new(big.Int).Sub(big.NewInt(0), new(big.Int).Lsh(big.NewInt(1), uint(bitPos)))
		masked := new(big.Int).And(p.TickBitmap.Data[wordPos], mask)
		// if there are no initialized ticks to the left of the current tick, return leftmost in the word
		initialized := masked.Cmp(big.NewInt(0)) != 0
		// overflow/underflow is possible, but prevented externally by limiting both tickSpacing and tick
		next := 0
		if initialized {
			leastSigBit, err := utils.MostSignificantBit(masked)
			if err != nil {
				panic(err)
			}
			next = (compressed + 1 + (int(leastSigBit) - int(bitPos))) * tickSpacing
		} else {
			next = (compressed + 1 + (math.MaxUint8 - int(bitPos))) * tickSpacing
		}
		return next, initialized
	}
}

func position(tick int32) (int16, uint8) {
	wordPos := int16(tick >> 8)
	bitPos := uint8(tick % 256)
	return wordPos, bitPos
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
		}
	}
	return tick, nil
}

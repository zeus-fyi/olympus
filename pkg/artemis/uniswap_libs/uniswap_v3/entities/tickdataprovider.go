package entities

import "math/big"

type Tick struct {
	Index          int      `json:"tick"`
	LiquidityGross *big.Int `json:"liquidityGross"`
	LiquidityNet   *big.Int `json:"liquidityNet"`
}

// Provides information about ticks

type TickDataProvider interface {

	/**
	 * Return information corresponding to a specific tick
	 * @param tick the tick to load
	 */

	GetTick(tick int) Tick

	/**
	 * Return the next tick that is initialized within a single word
	 * @param tick The current tick
	 * @param lte Whether the next tick should be lte the current tick
	 * @param tickSpacing The tick spacing of the pool
	 */

	NextInitializedTickWithinOneWord(tick int, lte bool, tickSpacing int) (int, bool)
}

type JSONTick struct {
	Index          int    `json:"tick"`
	LiquidityGross string `json:"liquidityGross"`
	LiquidityNet   string `json:"liquidityNet"`
}

func (t *Tick) ConvertToJSONType() JSONTick {
	return JSONTick{
		Index:          t.Index,
		LiquidityGross: t.LiquidityGross.String(),
		LiquidityNet:   t.LiquidityNet.String(),
	}
}

func (jt *JSONTick) ConvertToBigIntType() Tick {
	lg, _ := new(big.Int).SetString(jt.LiquidityGross, 10)
	ln, _ := new(big.Int).SetString(jt.LiquidityNet, 10)
	return Tick{
		Index:          jt.Index,
		LiquidityGross: lg,
		LiquidityNet:   ln,
	}
}

package entities

// A data provider for ticks that is backed by an in-memory array of ticks.

type TickListDataProvider struct {
	ticks []Tick `abi:"populatedTicks"`
}

func NewTickListDataProvider(ticks []Tick, tickSpacing int) (*TickListDataProvider, error) {
	if err := ValidateList(ticks, tickSpacing); err != nil {
		return nil, err
	}
	return &TickListDataProvider{ticks: ticks}, nil
}

func (p *TickListDataProvider) GetTick(tick int) Tick {
	return GetTick(p.ticks, tick)
}

func (p *TickListDataProvider) NextInitializedTickWithinOneWord(tick int, lte bool, tickSpacing int) (int, bool) {
	return NextInitializedTickWithinOneWord(p.ticks, tick, lte, tickSpacing)
}

type JSONTickListDataProvider struct {
	ticks []JSONTick `abi:"populatedTicks"`
}

func (p *TickListDataProvider) ConvertToJSONType() JSONTickListDataProvider {
	jticks := make([]JSONTick, len(p.ticks))
	for i, t := range p.ticks {
		jticks[i] = t.ConvertToJSONType()
	}
	return JSONTickListDataProvider{
		ticks: jticks,
	}
}

func (p *JSONTickListDataProvider) ConvertToBigIntType() *TickListDataProvider {
	ticks := make([]Tick, len(p.ticks))
	for i, t := range p.ticks {
		ticks[i] = t.ConvertToBigIntType()
	}
	return &TickListDataProvider{
		ticks: ticks,
	}
}

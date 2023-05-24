package v0

import "time"

const (
	mainnetGenesis = 1606824000
	goerliGenesis  = 1614588812

	mainnetBlockOffset = 10820069
)

func (c *LibV0) GetPendingMainnetSlotNum() int {
	t := time.Now().Unix()
	secondsSinceGenesis := t - mainnetGenesis
	return int(secondsSinceGenesis) / 12
}

func (c *LibV0) GetLatestMainnetBlockNumber() int {
	return c.GetPendingMainnetSlotNum() + mainnetBlockOffset - 1
}

func (c *LibV0) GetPendingMainnetBlockNumber() int {
	return c.GetPendingMainnetSlotNum() + mainnetBlockOffset
}

func (c *LibV0) GetSecsSinceLastMainnetSlot() int {
	t := time.Now().Unix()
	secondsSinceGenesis := t - mainnetGenesis
	return int(secondsSinceGenesis) % 12
}

func (c *LibV0) GetSecsSinceLastGoerliSlot() int {
	t := time.Now().Unix()
	secondsSinceGenesis := t - goerliGenesis
	return int(secondsSinceGenesis) % 12
}

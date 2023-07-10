package v0

import (
	"time"
)

const (
	mainnetGenesis = 1606824000
	goerliGenesis  = 1614588812
)

func (c *LibV0) GetPendingMainnetSlotNum() int {
	t := time.Now().Unix()
	secondsSinceGenesis := t - mainnetGenesis
	return int(secondsSinceGenesis) / 12
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

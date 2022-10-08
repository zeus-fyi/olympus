package v0

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

func (c *LibV0) ConvertEpochToSlot(epoch int64) string {
	return fmt.Sprintf("%d", epoch*int64(32))
}

func (c *LibV0) SlotToEpoch(slot int64) string {
	epochSlotMod := slot % 32
	if epochSlotMod != 0 {
		err := fmt.Errorf(fmt.Sprintf("slot %d was not at first slot in epoch, or in other words mod 32 != 0, but instead was: %d", slot, epochSlotMod))
		log.Err(err)
		return ""
	}
	return fmt.Sprintf("%d", slot/32)
}

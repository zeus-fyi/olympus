package misc

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

func ConvertEpochToSlot(epoch int64) string {
	return fmt.Sprintf("%d", epoch*int64(32))
}

func SlotToEpoch(slot int64) string {
	epochSlotMod := slot % 32
	if epochSlotMod != 0 {
		err := fmt.Errorf(fmt.Sprintf("slot %d was not at first slot in epoch, or in other words mod 32 != 0, but instead was: %d", slot, epochSlotMod))
		log.Err(err)
		return ""
	}
	return fmt.Sprintf("%d", slot/32)
}

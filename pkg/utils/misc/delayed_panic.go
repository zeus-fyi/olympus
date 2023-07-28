package misc

import (
	"time"

	"github.com/rs/zerolog/log"
)

func DelayedPanic(err error) {
	log.Err(err).Msg("DelayedPanic")
	time.Sleep(10 * time.Second)
	panic(err)
}

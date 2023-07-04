package v0

import (
	"math/rand"
	"time"
)

func (c *LibV0) Jitter(multiplier int) time.Duration {
	minDuration := 1 * time.Millisecond
	maxDuration := 10 * time.Millisecond
	jitter := time.Duration(multiplier) * (time.Duration(rand.Int63n(int64(maxDuration-minDuration))) + minDuration)
	return jitter
}

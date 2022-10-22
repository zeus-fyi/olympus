package v0

import (
	"math/rand"
	"time"
)

func (c *LibV0) UnixTimeStampNow() int {
	t := time.Now().UnixNano()
	rangeLower := 0
	rangeUpper := 999
	randomNum := rangeLower + rand.Intn(rangeUpper-rangeLower+1)
	return int(t) + randomNum
}

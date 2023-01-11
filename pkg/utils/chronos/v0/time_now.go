package v0

import (
	"math/rand"
	"time"
)

// UnixTimeStampNow sleeps for 1-10 nanoseconds after generation, to help prevent duplicate timestamps
func (c *LibV0) UnixTimeStampNow() int {
	t := time.Now().UnixNano()
	rangeLower := 1
	rangeUpper := 10
	randomNum := rand.Intn(rangeUpper - rangeLower)
	time.Sleep(time.Duration(randomNum) * time.Nanosecond)
	return int(t)
}

// UnixTimeStampNowRaw does not wait and can possibly provide duplicates if called in parallel
func (c *LibV0) UnixTimeStampNowRaw() int {
	t := time.Now().UnixNano()
	return int(t)
}

package v0

import (
	"strconv"
	"time"

	"github.com/patrickmn/go-cache"
)

var tsCache = cache.New(30*time.Second, 1*time.Minute)

// UnixTimeStampNow sleeps for 1-10 nanoseconds after generation, to help prevent duplicate timestamps
func (c *LibV0) UnixTimeStampNow() int {
	var t int64
	for {
		t = time.Now().UnixNano()
		intTime := int(t)
		key := strconv.Itoa(intTime)
		_, found := tsCache.Get(key)
		if !found {
			tsCache.Set(key, intTime, cache.DefaultExpiration)
			return intTime
		}
		time.Sleep(1 * time.Nanosecond)
	}
}

// UnixTimeStampNowRaw does not wait and can possibly provide duplicates when called in parallel
func (c *LibV0) UnixTimeStampNowRaw() int {
	t := time.Now().UnixNano()
	return int(t)
}

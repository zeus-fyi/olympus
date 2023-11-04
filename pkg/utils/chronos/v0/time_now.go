package v0

import (
	"strconv"
	"time"

	"github.com/patrickmn/go-cache"
)

var tsCache = cache.New(1*time.Second, 2*time.Second)

// UnixTimeStampNow uses cache to help prevent duplicate timestamps
func (c *LibV0) UnixTimeStampNow() int {
	var t int64
	for {
		t = time.Now().UnixNano()
		intTime := int(t)
		key := strconv.Itoa(intTime)
		_, found := tsCache.Get(key)
		if !found {
			tsCache.Set(key, intTime, 1*time.Second)
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

// UnixTimeStampNowSec does not wait and can possibly provide duplicates when called in parallel
func (c *LibV0) UnixTimeStampNowSec() int {
	t := time.Now().Unix()
	return int(t)
}

func (c *LibV0) ConvertUnixTimeStampToDate(uts int) time.Time {
	seconds := uts / 1e9
	return time.Unix(int64(seconds), 0)
}

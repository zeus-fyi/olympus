package v0

import "time"

func (c *LibV0) UnixTimeStampNow() int64 {
	return time.Now().Unix()
}

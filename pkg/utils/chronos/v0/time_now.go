package v0

import (
	"fmt"
	"math/big"
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

// ConvertTweetIDToUnixTimestamp converts a Twitter tweet ID to a Unix timestamp.
func (c *LibV0) ConvertTweetIDToUnixTimestamp(tweetID int) int {
	// The Twitter epoch time.
	const twitterEpoch int = 1288834974657

	// Right shift the tweet ID by 22 bits and add the Twitter epoch time.
	timestamp := (tweetID >> 22) + twitterEpoch

	return timestamp / 1000
}

// ConvertUnixTimestampToTweetID converts a Unix timestamp to a Twitter tweet ID.
func (c *LibV0) ConvertUnixTimestampToTweetID(unixTimestamp int) int {
	// The Twitter epoch time in milliseconds.
	const twitterEpoch int = 1288834974657

	// Convert the Unix timestamp from seconds to milliseconds.
	timestampMs := int(unixTimestamp) * 1000

	// Left shift the corrected timestamp by 22 bits.
	tweetID := (timestampMs - twitterEpoch) << 22

	return tweetID
}

// ConvertRedditIDToUnixTimestamp converts a Reddit post ID to a Unix timestamp.
func (c *LibV0) ConvertRedditIDToUnixTimestamp(redditID string) (int, error) {
	// Remove the prefix "t3_" if present.
	if len(redditID) > 3 && redditID[:3] == "t3_" {
		redditID = redditID[3:]
	}

	// Decode the Base36 encoded string.
	decodedID := new(big.Int)
	_, success := decodedID.SetString(redditID, 36)
	if !success {
		return 0, fmt.Errorf("failed to decode Base36 Reddit ID")
	}

	// Convert to Unix timestamp.
	return int(decodedID.Int64()), nil
}

func (c *LibV0) AdjustedUnixTimestampNowRaw(sinceUnixTimestamp int) int {
	// Convert sinceUnixTimestamp from seconds to duration in nanoseconds
	adjustmentDuration := time.Duration(sinceUnixTimestamp) * time.Second

	// Get the current time and apply the adjustment
	adjustedTime := time.Now().Add(adjustmentDuration)

	// Return the adjusted time in Unix nanoseconds
	return int(adjustedTime.UnixNano())
}

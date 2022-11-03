package misc

import "time"

func DelayedPanic(err error) {
	time.Sleep(30 * time.Second)
	panic(err)
}

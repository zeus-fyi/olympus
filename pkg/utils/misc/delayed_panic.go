package misc

import "time"

func DelayedPanic(err error) {
	time.Sleep(10 * time.Second)
	panic(err)
}

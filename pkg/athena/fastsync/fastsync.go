package athena_fastsync

import "context"

func Resync(ctx context.Context) {

}

func GitPull(ctx context.Context) {
	// can be a statefulset that scales up only when needed, then scales to zero
	// needs to git pull latest each time, including workspace stuff
	// needs to build binary
	// needs to swap binary
	// should delete old binary
	// then should restart app with new binary
}

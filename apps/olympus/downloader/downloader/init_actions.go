package snapshot_init

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/zeus/pkg/utils/ephemery_reset"
)

func InitAction(ctx context.Context, w WorkloadInfo) {
	// the below uses a switch case to download if an ephemeralClientName is used
	ephemery_reset.ExtractAndDecEphemeralTestnetConfig(Workload.DataDir, w.ClientName)
	switch w.WorkloadType {
	case "validatorClient":
		err := w.DataDir.WipeDirIn()
		if err != nil {
			log.Ctx(ctx).Panic().Err(err).Msg("failed to wipe validator dir on startup")
			panic(err)
		}
	}
}

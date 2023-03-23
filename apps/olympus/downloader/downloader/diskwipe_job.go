package olympus_snapshot_init

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	pg_poseidon "github.com/zeus-fyi/olympus/datastores/postgres/apps/poseidon"
)

func CheckForDiskWipeJobBeacon(ctx context.Context, w WorkloadInfo) {
	dw := pg_poseidon.DiskWipeOrchestration{
		ClientName: w.ClientName,
		OrchestrationJob: artemis_orchestrations.OrchestrationJob{
			Orchestrations: artemis_autogen_bases.Orchestrations{},
			Scheduled: artemis_autogen_bases.OrchestrationsScheduledToCloudCtxNs{
				Status: pg_poseidon.Pending,
			},
			CloudCtxNs: w.CloudCtxNs,
		},
	}
	shouldDiskWipe, err := dw.CheckForPendingDiskWipeJob(ctx)
	if err != nil {
		log.Ctx(ctx).Panic().Err(err).Msgf("failed to wipe %s data dir on startup", w.ClientName)
		panic(err)
	}
	if !shouldDiskWipe {
		log.Ctx(ctx).Info().Interface("w", w).Msg("no disk wipe job found, skipping")
		return
	}
	log.Ctx(ctx).Info().Interface("w", w).Msg("disk wipe job found, starting disk wipe")
	err = w.DataDir.WipeDirIn()
	if err != nil {
		log.Ctx(ctx).Panic().Err(err).Msgf("failed to wipe %s data dir on startup", w.ClientName)
		panic(err)
	}
	err = dw.MarkDiskWipeComplete(ctx)
	if err != nil {
		log.Ctx(ctx).Panic().Err(err).Msgf("failed to mark disk wipe complete for %s", w.ClientName)
		panic(err)
	}
}

package poseidon_orchestrations

import (
	"context"

	"github.com/rs/zerolog/log"
	pg_poseidon "github.com/zeus-fyi/olympus/datastores/postgres/apps/poseidon"
)

func (d *PoseidonSyncActivities) ScheduleDiskWipe(ctx context.Context, params pg_poseidon.DiskWipeOrchestration) error {
	err := params.ScheduleDiskWipe(ctx)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("PoseidonSyncActivities: ScheduleDiskWipe")
		return err
	}
	return err
}

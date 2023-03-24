package poseidon_orchestrations

import (
	"context"

	"github.com/rs/zerolog/log"
	pg_poseidon "github.com/zeus-fyi/olympus/datastores/postgres/apps/poseidon"
)

func (d *PoseidonSyncActivities) ScheduleDiskUpload(ctx context.Context, params pg_poseidon.UploadDataDirOrchestration) error {
	err := params.ScheduleUpload(ctx)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("PoseidonSyncActivities: ScheduleDiskUpload")
		return err
	}
	return err
}

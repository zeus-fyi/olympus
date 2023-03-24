package pg_poseidon

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
)

type UploadDataDirOrchestration struct {
	ClientName string
	artemis_orchestrations.OrchestrationJob
}

const (
	UploadOrchestration = "UploadDataDirOrchestration"
)

func (d *UploadDataDirOrchestration) ScheduleUpload(ctx context.Context) error {
	d.OrchestrationName = GetUploadDataDirJobName(d.ClientName)
	return d.InsertOrchestrationsScheduledToCloudCtxNsUsingName(ctx)
}

func (d *UploadDataDirOrchestration) MarkUploadComplete(ctx context.Context) error {
	d.OrchestrationName = GetUploadDataDirJobName(d.ClientName)
	d.Scheduled.Status = Complete
	return d.UpdateOrchestrationsScheduledToCloudCtxNs(ctx)
}

func (d *UploadDataDirOrchestration) CheckForPendingUploadJob(ctx context.Context) (bool, error) {
	d.OrchestrationName = GetUploadDataDirJobName(d.ClientName)
	d.Scheduled.Status = Pending
	isJobPending, err := d.SelectOrchestrationsAtCloudCtxNsWithStatus(ctx)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("error checking for pending upload job")
		return false, err
	}
	return isJobPending, nil
}

func GetUploadDataDirJobName(clientName string) string {
	return fmt.Sprintf("%sDataDirUpload", clientName)
}

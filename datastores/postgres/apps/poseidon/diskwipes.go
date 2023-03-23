package pg_poseidon

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
)

type DiskWipeOrchestration struct {
	ClientName string
	artemis_orchestrations.OrchestrationJob
}

const (
	Pending  = "Pending"
	Complete = "Complete"
	Sn       = "DiskWipeOrchestration"
)

func (d *DiskWipeOrchestration) ScheduleDiskWipe(ctx context.Context) error {
	d.OrchestrationName = GetDiskWipeDataDirJobName(d.ClientName)
	return d.InsertOrchestrationsScheduledToCloudCtxNsUsingName(ctx)
}

func (d *DiskWipeOrchestration) MarkDiskWipeComplete(ctx context.Context) error {
	d.OrchestrationName = GetDiskWipeDataDirJobName(d.ClientName)
	d.Scheduled.Status = Complete
	return d.UpdateOrchestrationsScheduledToCloudCtxNs(ctx)
}

func (d *DiskWipeOrchestration) CheckForPendingDiskWipeJob(ctx context.Context) (bool, error) {
	d.OrchestrationName = GetDiskWipeDataDirJobName(d.ClientName)
	d.Scheduled.Status = Pending
	isJobPending, err := d.SelectOrchestrationsAtCloudCtxNsWithStatus(ctx)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("Error checking for pending disk wipe job")
		return false, err
	}
	return isJobPending, nil
}

func GetDiskWipeDataDirJobName(clientName string) string {
	return fmt.Sprintf("%sDataDirDiskWipe", clientName)
}

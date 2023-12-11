package ai_platform_service_orchestrations

import (
	"context"

	"github.com/rs/zerolog/log"
)

func (z *ZeusAiPlatformActivities) CancelRun(ctx context.Context, wfID string) error {
	err := ZeusAiPlatformWorker.ExecuteCancelWorkflow(ctx, wfID)
	if err != nil {
		log.Err(err).Msg("CancelRun: ExecuteCancelWorkflow failed")
		return err
	}
	return nil
}

package ai_platform_service_orchestrations

import (
	"context"
	"strings"

	"github.com/rs/zerolog/log"
)

func (z *ZeusAiPlatformActivities) CancelRun(ctx context.Context, wfID string) error {
	err := ZeusAiPlatformWorker.ExecuteCancelWorkflow(ctx, wfID)
	if err != nil && !strings.Contains(err.Error(), "already completed") && !strings.Contains(err.Error(), "not found") {
		log.Err(err).Msg("CancelRun: ExecuteCancelWorkflow failed")
		return err
	}
	return nil
}

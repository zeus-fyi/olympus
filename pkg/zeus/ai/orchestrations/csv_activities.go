package ai_platform_service_orchestrations

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
)

func (z *ZeusAiPlatformActivities) SaveCsvTaskOutput(ctx context.Context, wr *artemis_orchestrations.AIWorkflowAnalysisResult, cp *MbChildSubProcessParams) (int, error) {
	if cp == nil {
		return 0, fmt.Errorf("SaveTaskOutput: cp is nil")
	}
	wio, werr := gs3wfs(ctx, cp)
	if werr != nil {
		log.Err(werr).Msg("TokenOverflowReduction: failed to select workflow io")
		return 0, werr
	}
	// todo add csv results
	fmt.Println(wio)
	return wr.WorkflowResultID, nil
}

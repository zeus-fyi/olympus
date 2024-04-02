package ai_platform_service_orchestrations

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
)

func (z *ZeusAiPlatformActivities) SaveCsvTaskOutput(ctx context.Context, cp *MbChildSubProcessParams, wr *artemis_orchestrations.AIWorkflowAnalysisResult) (int, error) {
	if cp == nil || wr == nil {
		return 0, fmt.Errorf("SaveTaskOutput: cp or wr is nil")
	}
	//err := artemis_orchestrations.InsertAiWorkflowAnalysisResult(ctx, wr)
	//if err != nil {
	//	log.Err(err).Interface("wr", wr).Interface("wr", wr).Msg("SaveTaskOutput: failed")
	//	return 0, err
	//}
	wio, werr := gs3wfs(ctx, cp)
	if werr != nil {
		log.Err(werr).Msg("TokenOverflowReduction: failed to select workflow io")
		return 0, werr
	}
	// todo add csv results
	fmt.Println(wio)

	if wio.PromptReduction != nil && wio.PromptReduction.PromptReductionSearchResults != nil {

		for _, sgpt := range wio.PromptReduction.PromptReductionSearchResults.OutSearchGroups {
			fmt.Println(aws.StringValue(sgpt.RetrievalName))
		}
	}

	/*
		at this stage:
			1. should add results to final csv for cycle

		done
			1. get s3 of main csv input ref from entities
		todo
			2. use csv-merge entity from input
			3. merge results
			4. save wf output to correct stage: ie final processed output

			DirIn:  fmt.Sprintf("/%s/%s/cycle/%d", ogk, wfRunName, cp.Wsr.RunCycle),
	*/

	gens, err := GetGlobalEntitiesFromRef(ctx, cp.Ou, cp.WfExecParams.WorkflowOverrides.WorkflowEntityRefs)
	if err != nil {
		log.Err(err).Msg("GetGlobalEntitiesFromRef: failed to select workflow io")
		return 0, err
	}

	var newCsvEntities []artemis_entities.UserEntity
	for _, gv := range gens {
		if artemis_entities.SearchLabelsForMatch("csv:source", gv) {
			mvs, merr := findMatchingNicknamesCsvMerge(gv, cp.WfExecParams.WorkflowOverrides.WorkflowEntities, wio)
			if merr != nil {
				return 0, merr
			}
			newCsvEntities = append(newCsvEntities, *mvs)
		}
	}

	// save newCsvEntities
	// test export
	return wr.WorkflowResultID, nil
}

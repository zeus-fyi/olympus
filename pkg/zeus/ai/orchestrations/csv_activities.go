package ai_platform_service_orchestrations

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
)

// TODO: next step SaveCsvTaskOutput
/*
	main csv is kept in global if using entity filter lookup; make copy of this for mutations

	fmt.Println(wio)
	if wio.PromptReduction != nil && wio.PromptReduction.PromptReductionSearchResults != nil {
		for _, sgpt := range wio.PromptReduction.PromptReductionSearchResults.OutSearchGroups {
			fmt.Println(aws.StringValue(sgpt.RetrievalName))
		}
	}
*/

func (z *ZeusAiPlatformActivities) SaveCsvTaskOutput(ctx context.Context, cp *MbChildSubProcessParams, wr *artemis_orchestrations.AIWorkflowAnalysisResult) (int, error) {
	if cp == nil || wr == nil {
		return 0, fmt.Errorf("SaveTaskOutput: cp or wr is nil")
	}
	// gets globals where needed
	gens, err := GetGlobalEntitiesFromRef(ctx, cp.Ou, cp.WfExecParams.WorkflowOverrides.WorkflowEntityRefs)
	if err != nil {
		log.Err(err).Msg("GetGlobalEntitiesFromRef: failed to select workflow io")
		return 0, err
	}

	// gets cycle stage values
	wio, werr := gs3wfs(ctx, cp)
	if werr != nil {
		log.Err(werr).Msg("TokenOverflowReduction: failed to select workflow io")
		return 0, werr
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

	mergeCsvEntities, err := getGlobalCsvMergedEntities(gens, cp, wio)
	if err != nil {
		log.Err(err).Msg("GetGlobalEntitiesFromRef: failed to select workflow io")
		return 0, err
	}
	fmt.Println(mergeCsvEntities, "mergeCsvEntities")
	// now merge these: newCsvEntities

	// save newCsvEntities
	// test export
	return wr.WorkflowResultID, nil
}

func getGlobalCsvMergedEntities(gens []artemis_entities.UserEntity, cp *MbChildSubProcessParams, wio *WorkflowStageIO) ([]artemis_entities.UserEntity, error) {
	var newCsvEntities []artemis_entities.UserEntity
	for _, gv := range gens {
		// since gens == global; use global label; csvSrcGlobalLabel
		if artemis_entities.SearchLabelsForMatch(csvSrcGlobalLabel, gv) {
			// todo verify label matching; they should share a label
			mvs, merr := FindAndMergeMatchingNicknamesByLabel(gv, cp.WfExecParams.WorkflowOverrides.WorkflowEntities, wio, csvSrcGlobalMergeLabel)
			if merr != nil {
				return nil, merr
			}
			newCsvEntities = append(newCsvEntities, *mvs)
		}
	}
	return newCsvEntities, nil
}

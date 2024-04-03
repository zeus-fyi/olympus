package ai_platform_service_orchestrations

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
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

/*
	wf override
	wfRunName := cp.WfExecParams.WorkflowOverrides.WorkflowRunName
*/

func (z *ZeusAiPlatformActivities) SaveCsvTaskOutput(ctx context.Context, cp *MbChildSubProcessParams, wr *artemis_orchestrations.AIWorkflowAnalysisResult) (int, error) {
	if cp == nil || wr == nil {
		return 0, fmt.Errorf("SaveTaskOutput: cp or wr is nil")
	}
	wfRunName := cp.WfExecParams.WorkflowOverrides.WorkflowRunName
	if len(wfRunName) <= 0 {
		return 0, fmt.Errorf("no wf run name provided")
	}
	// gets globals where needed
	gens, gerr := GetGlobalEntitiesFromRef(ctx, cp.Ou, cp.WfExecParams.WorkflowOverrides.WorkflowEntityRefs)
	if gerr != nil {
		log.Err(gerr).Msg("SaveCsvTaskOutput: GetGlobalEntitiesFromRef: failed to select workflow io")
		return 0, gerr
	}
	switch cp.Tc.TaskType {
	case AggTask:
		//for _, wfe := range cp.WfExecParams.WorkflowOverrides.WorkflowEntities {
		//
		//}
		// get all analysis output csvs and merge
		var sgss []*hera_search.SearchResultGroup
		for _, tv := range cp.WfExecParams.WorkflowTasks {
			tmn := cp.Tc.TaskName
			cp.Tc.TaskName = tv.AnalysisTaskName
			// gets cycle stage values
			wio, werr := gs3wfs(ctx, cp)
			if werr != nil {
				log.Err(werr).Msg("SaveCsvTaskOutput: gs3wfs failed to select workflow io")
				return 0, werr
			}
			sgss = append(sgss, wio.GetOutSearchGroups()...)
			cp.Tc.TaskName = tmn
		}
	default:
		// gets cycle stage values
		wio, werr := gs3wfs(ctx, cp)
		if werr != nil {
			log.Err(werr).Msg("SaveCsvTaskOutput: gs3wfs failed to select workflow io")
			return 0, werr
		}
		/*
			cp.WfExecParams.WorkflowOverrides.WorkflowEntities
			- this will have col name & json data for emRow

			data coming from search groups inputs
		*/
		//fmt.Println(mergeCsvEntities, "mergeCsvEntities")
		mergeCsvEntities, err := getGlobalCsvMergedEntities(gens, cp, wio)
		if err != nil {
			log.Err(err).Msg("SaveCsvTaskOutput: GetGlobalEntitiesFromRef: failed to select workflow io")
			return 0, err
		}
		for i, nev := range mergeCsvEntities {
			log.Info().Interface("i", i).Interface("nn", nev.Nickname).Msg("mergeCsvEntities")
			//mergeCsvEntities[i].Nickname = wfRunName
			// for now just saved under wf, later use label, csv under platform csv-exports
			_, err = S3WfRunImports(ctx, cp.Ou, wfRunName, &nev)
			if err != nil {
				log.Err(err).Msg("S3WfRunImports: failed to save merged result")
				return 0, err
			}
		}
	}
	err := artemis_orchestrations.InsertAiWorkflowAnalysisResult(ctx, wr)
	if err != nil {
		log.Err(err).Interface("wr", wr).Interface("wr", wr).Msg("SaveTaskOutput: failed")
		return 0, err
	}
	return wr.WorkflowResultID, nil
}

func getGlobalCsvMergedEntities(gens []artemis_entities.UserEntity, cp *MbChildSubProcessParams, wio *WorkflowStageIO) ([]artemis_entities.UserEntity, error) {
	var newCsvEntities []artemis_entities.UserEntity
	for _, gv := range gens {
		// since gens == global; use global label; csvSrcGlobalLabel
		if artemis_entities.SearchLabelsForMatch(csvSrcGlobalLabel, gv) {
			mvs, merr := FindAndMergeMatchingNicknamesByLabelPrefix(gv, cp.WfExecParams.WorkflowOverrides.WorkflowEntities, wio, csvSrcGlobalMergeLabel)
			if merr != nil {
				log.Err(merr).Msg("getGlobalCsvMergedEntities")
				return nil, merr
			}
			newCsvEntities = append(newCsvEntities, *mvs)
		}
	}
	return newCsvEntities, nil
}

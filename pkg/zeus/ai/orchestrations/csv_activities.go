package ai_platform_service_orchestrations

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
)

func (z *ZeusAiPlatformActivities) SaveCsvTaskOutput(ctx context.Context, cp *MbChildSubProcessParams, wr *artemis_orchestrations.AIWorkflowAnalysisResult) (int, error) {
	if cp == nil || wr == nil {
		return 0, fmt.Errorf("SaveTaskOutput: cp or wr is nil")
	}
	wfRunName := cp.WfExecParams.WorkflowOverrides.WorkflowRunName
	if len(wfRunName) <= 0 {
		return 0, fmt.Errorf("no wf run name provided")
	}
	// gets globals where needed
	gens, err := GetGlobalEntitiesFromRef(ctx, cp.Ou, cp.WfExecParams.WorkflowOverrides.WorkflowEntityRefs)
	if err != nil {
		log.Err(err).Msg("SaveCsvTaskOutput: GetGlobalEntitiesFromRef: failed to select workflow io")
		return 0, err
	}
	var merg []artemis_entities.UserEntity
	switch cp.Tc.TaskType {
	case AggTask:
		// get all analysis output csvs and merge
		for _, tv := range cp.WfExecParams.WorkflowTasks {
			tmn := cp.Tc.TaskName
			cp.Tc.TaskName = tv.AnalysisTaskName
			// gets cycle stage values
			wio, werr := gs3wfs(ctx, cp)
			if werr != nil {
				log.Err(werr).Msg("SaveCsvTaskOutput: gs3wfs failed to select workflow io")
				return 0, werr
			}
			if wio == nil {
				continue
			}
			merg, err = getGlobalCsvMergedEntities(gens, cp, wio)
			if err != nil {
				log.Err(err).Msg("SaveCsvTaskOutput: GetGlobalEntitiesFromRef: failed to select workflow io")
				return 0, err
			}
			gens = merg
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
		merg, err = getGlobalCsvMergedEntities(gens, cp, wio)
		if err != nil {
			log.Err(err).Msg("SaveCsvTaskOutput: GetGlobalEntitiesFromRef: failed to select workflow io")
			return 0, err
		}
		for i, nev := range merg {
			log.Info().Interface("i", i).Interface("nn", nev.Nickname).Msg("mergeCsvEntities")
			//mergeCsvEntities[i].Nickname = wfRunName
			// for now just saved under wf, later use label, csv under platform csv-exports
			if len(nev.Nickname) <= 0 {
				nev.Nickname = wfRunName
			}
			_, err = S3WfRunImports(ctx, cp.Ou, wfRunName, &nev)
			if err != nil {
				log.Err(err).Msg("S3WfRunImports: failed to save merged result")
				return 0, err
			}
		}
	}
	err = artemis_orchestrations.InsertAiWorkflowAnalysisResult(ctx, wr)
	if err != nil {
		log.Err(err).Interface("wr", wr).Interface("wr", wr).Msg("SaveTaskOutput: failed")
		return 0, err
	}
	return wr.WorkflowResultID, nil
}

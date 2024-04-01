package ai_platform_service_orchestrations

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
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

	/*
		at this stage:
			1. should add results to final csv for cycle
		todo
			1. get s3 of main csv input ref from entities
			2. use csv-merge entity from input
			3. merge results
			4. save wf output to correct stage: ie final processed output

			DirIn:  fmt.Sprintf("/%s/%s/cycle/%d", ogk, wfRunName, cp.Wsr.RunCycle),
	*/

	var gens []artemis_entities.UserEntity
	for _, ev := range cp.WfExecParams.WorkflowOverrides.WorkflowEntityRefs {
		ue := &artemis_entities.UserEntity{
			Nickname: ev.Nickname,
			Platform: ev.Platform,
		}
		ue, err := GetS3GlobalOrg(ctx, cp.Ou, ue)
		if err != nil {
			log.Err(err).Msg("TokenOverflowReduction: failed to select workflow io")
			return 0, err
		}
		if len(ue.MdSlice) > 0 {
			gens = append(gens, *ue)
		}
	}

	return wr.WorkflowResultID, nil
}

func GetGlobalEntitiesFromRef(ctx context.Context, ou org_users.OrgUser, refs []artemis_entities.EntitiesFilter) ([]artemis_entities.UserEntity, error) {
	var gens []artemis_entities.UserEntity
	for _, ev := range refs {
		ue := &artemis_entities.UserEntity{
			Nickname: ev.Nickname,
			Platform: ev.Platform,
		}
		ue, err := GetS3GlobalOrg(ctx, ou, ue)
		if err != nil {
			log.Err(err).Msg("TokenOverflowReduction: failed to select workflow io")
			return nil, err
		}
		if len(ue.MdSlice) > 0 {
			gens = append(gens, *ue)
		}
	}
	return gens, nil
}

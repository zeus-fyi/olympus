package ai_platform_service_orchestrations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	utils_csv "github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/csv"
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
	switch cp.Tc.TaskType {
	default:
		// gets cycle stage values
		wio, werr := gs3wfs(ctx, cp)
		if werr != nil {
			log.Err(werr).Msg("SaveCsvTaskOutput: gs3wfs failed to select workflow io")
			return 0, werr
		}
		var payloadMaps []map[string]interface{}
		sgo := wio.GetOutSearchGroups()
		if sgo != nil && len(sgo) > 0 {
			for _, sgi := range sgo {
				if sgi.ApiResponseResults != nil && len(sgi.ApiResponseResults) > 0 {
					for _, sr := range sgi.ApiResponseResults {
						if sr.WebResponse.Body != nil {
							payloadMaps = append(payloadMaps, sr.WebResponse.Body)
						}
					}
				}
			}
		}
		for i, source := range gens {
			tmp := source.MdSlice
			for mind, mi := range cp.WfExecParams.WorkflowOverrides.WorkflowEntities {
				for _, minv := range mi.MdSlice {
					if minv.JsonData != nil && string(minv.JsonData) != "null" {
						var cme utils_csv.CsvMergeEntity
						jerr := json.Unmarshal(minv.JsonData, &cme)
						if jerr != nil {
							log.Err(jerr).Interface("minv.JsonData", minv.JsonData).Msg(" json.Unmarshal(minv.JsonData, &emRow)")
							continue
						}
						log.Info().Interface("cme.MergeColName", cme.MergeColName).Msg("cme.MergeColName")
						//fmt.Println(cv)
						_, ms, merr := utils_csv.MergeCsvEntity(source, payloadMaps, cme, mind)
						if merr != nil {
							log.Err(merr).Msg("SaveCsvTaskOutput: MergeCsvEntity")
							return 0, err
						}
						if ms == nil {
							log.Warn().Msg("GenerateCycleReports: merged nil")
							continue
						}
						for _, vs := range ms {
							source = artemis_entities.UserEntity{
								Platform: "csv-exports",
								MdSlice: []artemis_entities.UserEntityMetadata{
									{
										TextData: aws.String(vs),
									},
								},
							}
						}
						rna := cp.GetRunName()
						if i > 0 {
							rna = fmt.Sprintf("%s_%d", rn, i)
						}
						_, err = S3WfRunImports(ctx, cp.Ou, rna, &source)
						if err != nil {
							log.Err(err).Msg("S3WfRunImports: failed to save merged result")
							return 0, err
						}
					}
				}
			}
			fmt.Println(source)
			fmt.Println(tmp)
		}
	}
	err = artemis_orchestrations.InsertAiWorkflowAnalysisResult(ctx, wr)
	if err != nil {
		log.Err(err).Interface("wr", wr).Interface("wr", wr).Msg("SaveTaskOutput: failed")
		return 0, err
	}
	return wr.WorkflowResultID, nil
}

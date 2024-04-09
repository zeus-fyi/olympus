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

func (z *ZeusAiPlatformActivities) GenerateCycleReports(ctx context.Context, cp *MbChildSubProcessParams) error {
	var sourceTaskIds []int
	for _, wfi := range cp.WfExecParams.WorkflowTasks {
		if wfi.AnalysisTaskID > 0 && wfi.AggTaskID == nil {
			sourceTaskIds = append(sourceTaskIds, wfi.AnalysisTaskID)
		}
		if wfi.AggTaskID != nil && *wfi.AggTaskID > 0 {
			sourceTaskIds = append(sourceTaskIds, *wfi.AggTaskID)
		}
	}
	uin := &artemis_entities.UserEntity{
		Nickname: cp.GetRunName(),
		Platform: "csv-exports",
	}
	var gens []artemis_entities.UserEntity
	ue, err := S3WfRunExport(ctx, cp.Ou, cp.GetRunName(), uin)
	if err != nil {
		log.Warn().Err(err).Msg("GenerateCycleReports: S3WfRunExport")
		err = nil
	}
	if ue != nil && ue.MdSlice != nil && err == nil {
		log.Info().Interface("ue.m", ue.MdSlice).Msg("UserEntity")
		gens = []artemis_entities.UserEntity{*ue}
	} else {
		// gets globals where needed
		gens, err = GetGlobalEntitiesFromRef(ctx, cp.Ou, cp.WfExecParams.WorkflowOverrides.WorkflowEntityRefs)
		if err != nil {
			log.Err(err).Msg("SaveCsvTaskOutput: GetGlobalEntitiesFromRef: failed to select workflow io")
			return err
		}
	}
	results, err := artemis_orchestrations.SelectAiWorkflowAnalysisResults(ctx, cp.Window, []int{cp.Oj.OrchestrationID}, sourceTaskIds)
	if err != nil {
		log.Err(err).Msg("GenerateCycleReports: SelectAiWorkflowAnalysisResults failed")
		return err
	}
	var resp []InputDataAnalysisToAgg
	for _, r := range results {
		b, berr := gs3wfsCustomTaskName(ctx, cp, fmt.Sprintf("%d", r.WorkflowResultID))
		if berr != nil {
			log.Err(berr).Msg("GenerateCycleReports: failed")
			continue
		}
		if b == nil {
			continue
		}
		tmp := InputDataAnalysisToAgg{}
		jerr := json.Unmarshal(b.Bytes(), &tmp)
		if jerr != nil {
			log.Err(jerr).Msg("GenerateCycleReports: failed")
			continue
		}
		resp = append(resp, tmp)
	}
	var payloadMaps []map[string]interface{}
	m := make(map[string]bool)
	var jsr []artemis_orchestrations.JsonSchemaDefinition
	for _, v := range resp {
		if v.CsvResponse != nil {
			payloadMaps = append(payloadMaps, v.CsvResponse)
		} else if v.ChatCompletionQueryResponse != nil {
			if v.ChatCompletionQueryResponse.JsonResponseResults != nil {
				// create report
				for _, vi := range v.ChatCompletionQueryResponse.JsonResponseResults {
					ht, herr := artemis_entities.HashParams(0, []interface{}{vi})
					if herr != nil {
						log.Err(herr).Msg("GenerateCycleReports: failed")
						continue
					}
					if _, ok := m[ht]; ok {
						continue
					}
					jsr = append(jsr, vi)
					m[ht] = true
				}
			}
		} else if v.SearchResultGroup != nil {
			if v.SearchResultGroup.RegexSearchResults != nil {
				fmt.Println("GenerateCycleReports: RegexSearchResults", v.SearchResultGroup.RegexSearchResults)
			} else if v.SearchResultGroup.ApiResponseResults != nil {
				fmt.Println("ApiResponseResults", v.SearchResultGroup.ApiResponseResults)
			} else if v.SearchResultGroup.SearchResults != nil {
				fmt.Println("SearchResults", v.SearchResultGroup.SearchResults)
			}
		}
	}
	if len(jsr) > 0 {
		payloadMaps = append(payloadMaps, artemis_orchestrations.CreateMapInterfaceFromAssignedSchemaFields(jsr)...)
	}

	for i, source := range gens {
		//tmp := source.MdSlice
		for _, mi := range cp.WfExecParams.WorkflowOverrides.WorkflowEntities {
			for _, minv := range mi.MdSlice {
				if minv.JsonData != nil && string(minv.JsonData) != "null" {
					var cme utils_csv.CsvMergeEntity
					jerr := json.Unmarshal(minv.JsonData, &cme)
					if jerr != nil {
						log.Err(jerr).Interface("minv.JsonData", minv.JsonData).Msg(" json.Unmarshal(minv.JsonData, &emRow)")
						continue
					}
					// {"MergeColName":"Email","Rows":{"alex@zeus.fyi":[0,2],"leevar@gmail.com":[1,3]}}
					cnT := cme.MergeColName
					log.Info().Interface("cme.MergeColName", cme.MergeColName).Msg("cme.MergeColName")
					cv := convEntityToCsvCol(cnT, payloadMaps)
					//fmt.Println(cv)
					merged, merr := utils_csv.MergeCsvEntity(source, cv, cme)
					if merr != nil {
						log.Err(merr).Msg("GenerateCycleReports: MergeCsvEntity")
						return err
					}
					if merged == nil {
						log.Warn().Msg("GenerateCycleReports: merged nil")
						continue
					}
					mergedCsvStr, merr := utils_csv.PayloadToCsvString(merged)
					if merr != nil {
						log.Err(merr).Msg("GenerateCycleReports: MergeCsvEntity")
						return merr
					}
					source = artemis_entities.UserEntity{
						Platform: "csv-exports",
						MdSlice: []artemis_entities.UserEntityMetadata{
							{
								TextData: aws.String(mergedCsvStr),
							},
						},
					}
					rna := cp.GetRunName()
					if i > 0 {
						rna = fmt.Sprintf("%s_%d", rn, i)
					}
					_, err = S3WfRunImports(ctx, cp.Ou, rna, &source)
					if err != nil {
						log.Err(err).Msg("S3WfRunImports: failed to save merged result")
						return err
					}
				}
			}
		}
		//fmt.Println(source)
		//fmt.Println(tmp)
	}
	return nil
}

func convEntityToCsvCol(cn string, plms []map[string]interface{}) []map[string]interface{} {
	//m := make(map[string]interface{})
	for i, pl := range plms {
		if v, ok := pl["entity"]; ok {
			delete(pl, "entity")
			pl[cn] = v
			plms[i] = pl
		}
	}
	return mergeMaps(plms, cn)
}

func mergeMaps(plms []map[string]interface{}, uniqueKey string) []map[string]interface{} {
	// Use a map to track the combined entries
	combinedEntries := make(map[interface{}]map[string]interface{})
	var result []map[string]interface{}

	for _, plm := range plms {
		keyVal := plm[uniqueKey]
		if combined, ok := combinedEntries[keyVal]; ok {
			// If the entry already exists, merge it
			for k, v := range plm {
				if k != uniqueKey {
					combined[k] = v
				}
			}
		} else {
			// If it's a new entry, create it and add to the combinedEntries
			newEntry := make(map[string]interface{})
			for k, v := range plm {
				newEntry[k] = v
			}
			combinedEntries[keyVal] = newEntry
		}
	}

	// Convert the combined entries map to a slice
	for _, v := range combinedEntries {
		result = append(result, v)
	}

	return result
}

//merg, err = getGlobalCsvMergedEntities(gens, cp, wio)
//if err != nil {
//log.Err(err).Msg("SaveCsvTaskOutput: GetGlobalEntitiesFromRef: failed to select workflow io")
//return 0, err

/*
	ret-only:
		ws.PromptReduction.PromptReductionSearchResults.OutSearchGroups
	analysis/agg-json:
		need wk-result id for every chunk
*/

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
	db := AiAggregateAnalysisRetrievalTaskInputDebug{
		Cp: cp,
	}
	db.Save()
	var sourceTaskIds []int
	for _, wfi := range cp.WfExecParams.WorkflowTasks {
		if wfi.AnalysisTaskID > 0 && wfi.AggTaskID == nil {
			sourceTaskIds = append(sourceTaskIds, wfi.AnalysisTaskID)
		}
		if wfi.AggTaskID != nil && *wfi.AggTaskID > 0 {
			sourceTaskIds = append(sourceTaskIds, *wfi.AggTaskID)
		}
	}
	results, err := artemis_orchestrations.SelectAiWorkflowAnalysisResults(ctx, cp.Window, []int{cp.Oj.OrchestrationID}, sourceTaskIds)
	if err != nil {
		log.Err(err).Msg("AiAggregateAnalysisRetrievalTask: SelectAiWorkflowAnalysisResults failed")
		return err
	}
	var resp []InputDataAnalysisToAgg
	for _, r := range results {
		b, berr := gs3wfsCustomTaskName(ctx, cp, fmt.Sprintf("%d", r.WorkflowResultID))
		if berr != nil {
			log.Err(berr).Msg("AiAggregateAnalysisRetrievalTask: failed")
			continue
		}
		if b == nil {
			continue
		}
		tmp := InputDataAnalysisToAgg{}
		jerr := json.Unmarshal(b.Bytes(), &tmp)
		if jerr != nil {
			log.Err(jerr).Msg("AiAggregateAnalysisRetrievalTask: failed")
			continue
		}
		resp = append(resp, tmp)
	}
	m := make(map[string]bool)
	var jsr []artemis_orchestrations.JsonSchemaDefinition
	for _, v := range resp {
		if v.ChatCompletionQueryResponse != nil {
			if v.ChatCompletionQueryResponse.JsonResponseResults != nil {
				// create report
				for _, vi := range v.ChatCompletionQueryResponse.JsonResponseResults {
					ht, herr := artemis_entities.HashParams(0, []interface{}{vi})
					if herr != nil {
						log.Err(herr).Msg("AiAggregateAnalysisRetrievalTask: failed")
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
				fmt.Println("RegexSearchResults", v.SearchResultGroup.RegexSearchResults)
			} else if v.SearchResultGroup.ApiResponseResults != nil {
				fmt.Println("ApiResponseResults", v.SearchResultGroup.ApiResponseResults)
			} else if v.SearchResultGroup.SearchResults != nil {
				fmt.Println("SearchResults", v.SearchResultGroup.SearchResults)
			}
		}
	}
	var payloadMaps []map[string]interface{}
	if len(jsr) > 0 {
		payloadMaps = append(payloadMaps, artemis_orchestrations.CreateMapInterfaceFromAssignedSchemaFields(jsr)...)
	}
	// gets globals where needed
	gens, err := GetGlobalEntitiesFromRef(ctx, cp.Ou, cp.WfExecParams.WorkflowOverrides.WorkflowEntityRefs)
	if err != nil {
		log.Err(err).Msg("SaveCsvTaskOutput: GetGlobalEntitiesFromRef: failed to select workflow io")
		return err
	}
	for i, source := range gens {
		tmp := source.MdSlice
		var cme utils_csv.CsvMergeEntity
		for _, mi := range cp.WfExecParams.WorkflowOverrides.WorkflowEntities {
			for _, minv := range mi.MdSlice {
				if minv.JsonData != nil && string(minv.JsonData) != "null" {
					jerr := json.Unmarshal(minv.JsonData, &cme)
					if jerr != nil {
						log.Err(jerr).Interface("minv.JsonData", minv.JsonData).Msg(" json.Unmarshal(minv.JsonData, &emRow)")
						continue
					}
				}
				cv := convEntityToCsvCol(cme.MergeColName, payloadMaps)
				fmt.Println(cv)
				merged, merr := utils_csv.MergeCsvEntity(source, cv, cme)
				if merr != nil {
					log.Err(merr).Msg("GenerateCycleReports: MergeCsvEntity")
					return err
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
				rn := cp.GetRunName()
				if i > 0 {
					rn = fmt.Sprintf("%s_%d", rn, i)
				}
				_, err = S3WfRunImports(ctx, cp.Ou, rn, &source)
				if err != nil {
					log.Err(err).Msg("S3WfRunImports: failed to save merged result")
					return err
				}
			}
		}
		fmt.Println(source)
		fmt.Println(tmp)
	}
	return nil
}

func convEntityToCsvCol(cn string, plms []map[string]interface{}) []map[string]interface{} {
	for i, pl := range plms {
		if v, ok := pl["entity"]; ok {
			delete(pl, "entity")
			pl[cn] = v
			plms[i] = pl
		}
		if v, ok := pl["summary"]; ok {
			delete(pl, "summary")
			pl[fmt.Sprintf("%s_%s", cn, "AI_Response")] = v
			plms[i] = pl
		}
	}
	return plms
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

package ai_platform_service_orchestrations

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	utils_csv "github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/csv"
)

const (
	mergeRetTag = "csv:global:merge:ret:"
)

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

// FindAndMergeMatchingNicknamesByLabelPrefix finds using retrieval name on search group and gets web response body agg
func FindAndMergeMatchingNicknamesByLabelPrefix(source artemis_entities.UserEntity, entities []artemis_entities.UserEntity, wsi *WorkflowStageIO, label string) (*artemis_entities.UserEntity, error) {
	if wsi == nil {
		return nil, nil
	}
	if source.Nickname == "" {
		return nil, fmt.Errorf("source nn empty")
	}
	fnn := source.Nickname
	// assume known for now ^
	var mes []artemis_entities.UserEntity
	for _, ev := range entities {
		if ev.Nickname == fnn && artemis_entities.SearchLabelsForPrefixMatch(label, ev) {
			mes = append(mes, ev)
		}
	}
	log.Info().Interface("mes", mes).Msg("findMatchingNicknamesByLabel: SearchLabelsForMatch(iter)")
	return mergeCsvs(source, mes, wsi)
}

// deprecated
func mergeCsvs(source artemis_entities.UserEntity, mergeIn []artemis_entities.UserEntity, wsi *WorkflowStageIO) (*artemis_entities.UserEntity, error) {
	var results []hera_search.SearchResult
	// todo; multi?
	cme := utils_csv.CsvMergeEntity{}
	for _, mi := range mergeIn {
		for _, minv := range mi.MdSlice {
			if minv.JsonData != nil && string(minv.JsonData) != "null" {
				jerr := json.Unmarshal(minv.JsonData, &cme)
				if jerr != nil {
					log.Err(jerr).Interface("minv.JsonData", minv.JsonData).Msg(" json.Unmarshal(minv.JsonData, &emRow)")
					continue
				}
			}
			gl := mi.GetStrLabels()
			rm := mergeRets(gl)
			sgs := wsi.GetSearchGroupsOutByRetNameMatch(rm)
			for _, sg := range sgs {
				if sg.ApiResponseResults != nil {
					results = append(results, sg.ApiResponseResults...)
				} else if sg.RegexSearchResults != nil {
					results = append(results, sg.RegexSearchResults...)
				}
			}
		}
	}
	var appendCsvEntry []map[string]interface{}
	for _, v := range results {
		if v.WebResponse.Body != nil {
			log.Info().Interface(" v.WebResponse.Body", v.WebResponse.Body).Msg("appendCsvEntry: results")
			if len(v.WebResponse.Body) > 0 {
				appendCsvEntry = append(appendCsvEntry, v.WebResponse.Body)
			}
		}
	}
	merged, _, err := utils_csv.MergeCsvEntity(source, appendCsvEntry, cme)
	if err != nil {
		log.Err(err).Msg("mergeCore")
		return nil, err
	}
	mergedCsvStr, err := utils_csv.PayloadToCsvString(merged)
	if err != nil {
		log.Err(err).Msg("PayloadToCsvString")
		return nil, err
	}
	log.Info().Interface("mergedCsvStr", mergedCsvStr).Msg("mergeCsvs: PayloadToCsvString")
	csvMerge := &artemis_entities.UserEntity{
		Platform: "csv-exports",
		MdSlice: []artemis_entities.UserEntityMetadata{
			{
				TextData: aws.String(mergedCsvStr),
			},
		},
	}
	return csvMerge, nil
}

func mergeRets(lbs []string) map[string]bool {
	rets := make(map[string]bool)
	for _, lb := range lbs {
		if strings.HasPrefix(lb, mergeRetTag) {
			rets[strings.TrimPrefix(lb, mergeRetTag)] = true
		}
	}
	if len(rets) <= 0 {
		log.Warn().Interface("lbs", lbs).Msg("mergeRets: empty rets")
	}
	return rets
}

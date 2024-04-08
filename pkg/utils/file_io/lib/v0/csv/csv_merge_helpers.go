package utils_csv

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
)

type CsvMergeEntity struct {
	MergeColName string
	Rows         map[string][]int
}

func NewCsvMergeEntity() CsvMergeEntity {
	return CsvMergeEntity{
		Rows: make(map[string][]int),
	}
}

func NewCsvMergeEntityFromSrc(colName string, emRow map[string][]int) CsvMergeEntity {
	return CsvMergeEntity{
		MergeColName: colName,
		Rows:         emRow,
	}
}

func NewCsvMergeEntityFromSrcBin(colName string, emRow map[string][]int) ([]byte, error) {
	cme := CsvMergeEntity{
		MergeColName: colName,
		Rows:         emRow,
	}

	b, err := json.Marshal(cme)
	if err != nil {
		log.Err(err).Msg("failed to marshal emRow")
		return nil, err
	}
	return b, nil
}

func MergeCsvEntity(source artemis_entities.UserEntity, appendCsvEntry []map[string]interface{}, cme CsvMergeEntity) ([]map[string]string, error) {
	if len(appendCsvEntry) == 0 {
		for _, v := range source.MdSlice {
			if v.TextData != nil && *v.TextData != "" {
				csvMapMerge, err := ParseCsvStringToMap(*v.TextData)
				if err != nil {
					log.Warn().Err(err).Msg("ParseCsvStringToMap")
					err = nil
					return nil, nil
				} else {
					return csvMapMerge, nil
				}
			}
		}
	}
	var merged []map[string]string
	for _, v := range source.MdSlice {
		if v.TextData != nil && len(*v.TextData) > 0 {
			csvMap, err := ParseCsvStringToMap(*v.TextData)
			if err != nil {
				return nil, err
			}
			pscsv, perr := PayloadV2ToCsvString(appendCsvEntry)
			if perr != nil {
				log.Err(perr).Interface("appendCsvEntry", appendCsvEntry).Msg("PayloadV2ToCsvString: ")
				return nil, perr
			}
			csvMapMerge, err := ParseCsvStringToMap(pscsv)
			if err != nil {
				return nil, err
			}
			log.Info().Interface("csvMapMerge", csvMapMerge).Msg("ParseCsvStringToMap: csvMapMerge")
			merged, err = AppendCsvData(csvMap, csvMapMerge, cme.MergeColName, cme.Rows)
			if err != nil {
				log.Err(err).Interface("merged", merged).Msg("mergeRets: empty rets")
				return nil, err
			}
			log.Info().Interface("merged", merged).Msg("appendCsvData: merged")
		}
	}
	return merged, nil
}

func AppendCsvData(inputCsvData, csvData []map[string]string, colName string, emRow map[string][]int) ([]map[string]string, error) {
	// Iterate through csvData to find and merge matching rows
	for _, dataRow := range csvData {
		email := dataRow[colName]
		if indices, ok := emRow[email]; ok {
			// If a matching row is found, merge the data
			for _, index := range indices {
				for key, value := range dataRow {
					inputCsvData[index][key] = value
				}
			}
		}
	}
	return inputCsvData, nil
}
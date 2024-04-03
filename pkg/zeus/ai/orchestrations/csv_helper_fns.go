package ai_platform_service_orchestrations

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
)

const mergeRetTag = "csv:global:merge:ret:"

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

func appendCsvData(inputCsvData, csvData []map[string]string, colName string, emRow map[string][]int) ([]map[string]string, error) {
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

func PayloadV2ToCsvString(payload []map[string]interface{}) (string, error) {
	if len(payload) == 0 {
		return "", fmt.Errorf("empty or nil payload")
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Extract and sort headers from the first map to ensure consistent column ordering
	headers := make([]string, 0, len(payload[0]))
	for key := range payload[0] {
		headers = append(headers, key)
	}
	sort.Strings(headers)

	// Write CSV header
	if err := writer.Write(headers); err != nil {
		return "", fmt.Errorf("error writing header to CSV: %w", err)
	}

	// Write each map as a CSV row
	for _, record := range payload {
		row := make([]string, len(headers))
		for i, header := range headers {
			value, ok := record[header]
			if !ok {
				// If a key is missing in a record, an empty string will be used as its value
				row[i] = ""
				continue
			}
			// Convert the interface value to string; more sophisticated conversion might be needed based on actual types
			row[i] = fmt.Sprintf("%v", value)
		}
		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("error writing record to CSV: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("error flushing CSV writer: %w", err)
	}

	return buf.String(), nil
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

func mergeCsvs(source artemis_entities.UserEntity, mergeIn []artemis_entities.UserEntity, wsi *WorkflowStageIO) (*artemis_entities.UserEntity, error) {
	var results []hera_search.SearchResult
	var colName string
	var emRow map[string][]int
	// todo; multi?
	for _, mi := range mergeIn {
		for _, minv := range mi.MdSlice {
			if minv.TextData != nil && len(*minv.TextData) > 0 {
				colName = *minv.TextData
			}
			if minv.JsonData != nil && string(minv.JsonData) != "null" {
				jerr := json.Unmarshal(minv.JsonData, &emRow)
				if jerr != nil {
					log.Err(jerr).Interface("minv.JsonData", minv.JsonData).Msg(" json.Unmarshal(minv.JsonData, &emRow)")
				}
			}
			//
			gl := mi.GetStrLabels()
			sgs := wsi.GetSearchGroupsOutByRetNameMatch(mergeRets(gl))
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

	var merged []map[string]string
	for _, v := range source.MdSlice {
		if v.JsonData != nil && string(v.JsonData) != "null" {
			err := json.Unmarshal(v.JsonData, &emRow)
			if err != nil {
				log.Err(err).Interface("v.JsonData", v.JsonData)
				return nil, err
			}
		}
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
			merged, err = appendCsvData(csvMap, csvMapMerge, colName, emRow)
			if err != nil {
				log.Err(err).Interface("merged", merged).Msg("mergeRets: empty rets")
				return nil, err
			}
			log.Info().Interface("merged", merged).Msg("appendCsvData: merged")
		}
	}
	mergedCsvStr, err := PayloadToCsvString(merged)
	if err != nil {
		log.Err(err).Msg("PayloadToCsvString")
		return nil, err
	}
	log.Info().Interface("mergedCsvStr", mergedCsvStr).Msg("mergeCsvs: PayloadToCsvString")
	csvMerge := &artemis_entities.UserEntity{
		Nickname: wsi.WorkflowOverrides.WorkflowRunName,
		Platform: "csv-exports",
		MdSlice: []artemis_entities.UserEntityMetadata{
			{
				TextData: aws.String(mergedCsvStr),
			},
		},
	}
	return csvMerge, nil
}

func ParseCsvStringToMap(csvString string) ([]map[string]string, error) {
	// Create a new reader from the CSV string
	reader := csv.NewReader(strings.NewReader(csvString))

	// Read the first row to get column headers
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	var records []map[string]string

	// Read each record after the header
	for {
		record, rerr := reader.Read()
		if rerr != nil {
			break // Stop reading when we reach the end of the file or an error
		}

		// Create a map for each record
		rowMap := make(map[string]string)
		for i, value := range record {
			rowMap[headers[i]] = value
		}

		// Append the map to the slice of records
		records = append(records, rowMap)
	}

	// If the loop exits due to an error other than EOF, return the error
	if err != nil && err.Error() != "EOF" {
		return nil, err
	}

	return records, nil
}

func PayloadToCsvString(payload []map[string]string) (string, error) {
	if payload == nil {
		return "", fmt.Errorf("empty or nil payload")
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write CSV header if necessary
	// Assuming all maps in ContactsCsv have the same keys
	headers := make([]string, 0, len(payload[0]))
	for key := range payload[0] {
		headers = append(headers, key)
	}
	if err := writer.Write(headers); err != nil {
		return "", fmt.Errorf("error writing header to CSV: %w", err)
	}

	// Write data
	for _, record := range payload {
		row := make([]string, 0, len(record))
		for _, header := range headers {
			row = append(row, record[header])
		}
		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("error writing record to CSV: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("error flushing CSV writer: %w", err)
	}
	return buf.String(), nil
}

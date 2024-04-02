package ai_platform_service_orchestrations

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
)

const mergeRetTag = "merge:ret"

func mergeRets(lbs []string) map[string]bool {
	rets := make(map[string]bool)
	for _, lb := range lbs {
		if strings.HasPrefix(lb, mergeRetTag) {
			rets[strings.TrimPrefix(lb, mergeRetTag)] = true
		}
	}
	return rets
}

func appendCsvData(csvContacts, csvData []map[string]string, colName string, emRow map[string][]int) ([]map[string]string, error) {
	// Iterate through csvData to find and merge matching rows
	for _, dataRow := range csvData {
		email := dataRow[colName]
		if indices, ok := emRow[email]; ok {
			// If a matching row is found, merge the data
			for _, index := range indices {
				for key, value := range dataRow {
					csvContacts[index][key] = value
				}
			}
		}
	}
	return csvContacts, nil
}

func mergeCsvs(source artemis_entities.UserEntity, mergeIn []artemis_entities.UserEntity, wsi *WorkflowStageIO) (*artemis_entities.UserEntity, error) {
	for _, mi := range mergeIn {
		sgs := wsi.GetSearchGroupsOutByRetNameMatch(mergeRets(mi.GetStrLabels()))
		for _, sg := range sgs {
			if sg.ApiResponseResults != nil {

			} else if sg.RegexSearchResults != nil {

			}
		}
	}

	var tmpCsv []map[string]string
	var merged []map[string]string
	var colName string
	var emRow map[string][]int
	for _, v := range source.MdSlice {
		if v.TextData != nil && len(*v.TextData) > 0 {
			csvMap, err := ParseCsvStringToMap(*v.TextData)
			if err != nil {
				return nil, err
			}
			merged, err = appendCsvData(csvMap, tmpCsv, colName, emRow)
			if err != nil {
				return nil, err
			}
		}
	}
	/*
		todo: get body, transform to csv type
		then call appendCsvData
	*/

	mergedCsvStr, err := PayloadToCsvString(merged)
	if err != nil {
		return nil, err
	}
	// now get bodies; then merge via
	csvMerge := &artemis_entities.UserEntity{
		Nickname: "",
		Platform: "flows",
		MdSlice: []artemis_entities.UserEntityMetadata{
			{
				TextData: aws.String(mergedCsvStr),
			},
		},
	}
	return csvMerge, nil
}

func findMatchingNicknamesCsvMerge(source artemis_entities.UserEntity, entities []artemis_entities.UserEntity, wsi *WorkflowStageIO) (*artemis_entities.UserEntity, error) {
	if wsi == nil {
		return nil, nil
	}
	fnn := source.Nickname
	// assume known for now ^
	var mes []artemis_entities.UserEntity
	for _, ev := range entities {
		if ev.Nickname == fnn && artemis_entities.SearchLabelsForMatch("csv:merge", ev) {
			mes = append(mes, ev)
		}
	}
	return mergeCsvs(source, mes, wsi)
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

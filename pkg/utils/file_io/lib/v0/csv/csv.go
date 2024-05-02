package utils_csv

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"sort"
	"strings"
)

func ParseCsvStringOrderedHeaders(csvString string) ([]string, error) {
	// Create a new reader from the CSV string
	reader := csv.NewReader(strings.NewReader(csvString))

	// Read the first row to get column headers
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	// Return the headers, which are the first row of the CSV
	return headers, nil
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
	if payload == nil || len(payload) == 0 {
		return "", fmt.Errorf("empty or nil payload")
	}
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	// Extract and sort the headers
	headers := make([]string, 0, len(payload[0]))
	for key := range payload[0] {
		headers = append(headers, key)
	}
	sort.Strings(headers) // Sort the headers alphabetically
	// Write CSV header
	if err := writer.Write(headers); err != nil {
		return "", fmt.Errorf("error writing header to CSV: %w", err)
	}
	// Write data rows
	for _, record := range payload {
		row := make([]string, 0, len(headers))
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

func SortCSV(csvString string, orderedHeaders []string) (string, error) {
	reader := csv.NewReader(strings.NewReader(csvString))
	records, err := reader.ReadAll()
	if err != nil {
		return "", err
	}

	if len(records) < 1 {
		return "", errors.New("csv contains no data")
	}

	originalHeaders := records[0]
	headerMap := make(map[string]int) // Map to store header to index
	for i, header := range originalHeaders {
		headerMap[header] = i
	}

	// Prepare the new header order
	newHeaders := make([]string, 0, len(originalHeaders))
	knownHeaders := make(map[string]bool)
	for _, header := range orderedHeaders {
		if _, exists := headerMap[header]; exists {
			newHeaders = append(newHeaders, header)
			knownHeaders[header] = true
		}
	}

	// Append unknown headers in alphabetical order
	unknownHeaders := make([]string, 0)
	for header := range headerMap {
		if !knownHeaders[header] {
			unknownHeaders = append(unknownHeaders, header)
		}
	}
	sort.Strings(unknownHeaders)
	newHeaders = append(newHeaders, unknownHeaders...)

	// Create an index map based on new headers
	newIndexMap := make(map[string]int)
	for i, header := range newHeaders {
		newIndexMap[header] = i
	}

	// Rearrange each record according to the new header order
	sortedRecords := make([][]string, len(records))
	sortedRecords[0] = newHeaders // Set the new headers as the first record
	for i, record := range records {
		if i == 0 {
			continue
		}
		newRecord := make([]string, len(newHeaders))
		for j, header := range newHeaders {
			if oldIndex, exists := headerMap[header]; exists {
				newRecord[j] = record[oldIndex]
			}
		}
		sortedRecords[i] = newRecord
	}

	// Convert the records back to a CSV string
	var sb strings.Builder
	writer := csv.NewWriter(&sb)
	if err := writer.WriteAll(sortedRecords); err != nil {
		return "", err
	}
	writer.Flush()

	return sb.String(), nil
}

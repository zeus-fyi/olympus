package zeus_v1_ai

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

type FlowsActionsRequest struct {
	FlowsCsvPayload `json:",inline"`
	Stages          map[string]bool   `json:"stages"`
	CommandPrompts  map[string]string `json:"commandPrompts"`
}

type FlowsCsvPayload struct {
	ContactsCsvStr string              `json:"contentContactsCsvStr"`
	ContactsCsv    []map[string]string `json:"contentContactsCsv,omitempty"`
	PromptsCsvStr  string              `json:"promptsCsvStr"`
	PromptsCsv     []map[string]string `json:"promptsCsv,omitempty,omitempty"`
}

func FlowsActionsRequestHandler(c echo.Context) error {
	request := new(FlowsActionsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return PostFlowsActionsRequest(c, *request)
}

func PostFlowsActionsRequest(c echo.Context, fa FlowsActionsRequest) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	fmt.Println(ou, fa)
	return c.JSON(http.StatusOK, nil)
}

func parseCsvStringToMap(csvString string) ([]map[string]string, error) {
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

func payloadToCsvString(payload *FlowsCsvPayload) (string, error) {
	if payload == nil || len(payload.ContactsCsv) == 0 {
		return "", fmt.Errorf("empty or nil payload")
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write CSV header if necessary
	// Assuming all maps in ContactsCsv have the same keys
	headers := make([]string, 0, len(payload.ContactsCsv[0]))
	for key := range payload.ContactsCsv[0] {
		headers = append(headers, key)
	}
	if err := writer.Write(headers); err != nil {
		return "", fmt.Errorf("error writing header to CSV: %w", err)
	}

	// Write data
	for _, record := range payload.ContactsCsv {
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

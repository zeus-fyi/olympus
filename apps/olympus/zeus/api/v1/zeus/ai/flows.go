package zeus_v1_ai

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	ai_platform_service_orchestrations "github.com/zeus-fyi/olympus/pkg/zeus/ai/orchestrations"
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

type FlowsActionsGetRequest struct{}

func FlowsExportCsvRequestHandler(c echo.Context) error {
	request := new(FlowsActionsGetRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Err(err).Msg("invalid ID parameter")
		return c.JSON(http.StatusBadRequest, "invalid ID parameter")
	}
	return ExportRunCsvRequest(c, id)
}

func ExportRunCsvRequest(c echo.Context, id int) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	ojsRuns, err := artemis_orchestrations.SelectAiSystemOrchestrations(c.Request().Context(), ou, id)
	if err != nil {
		log.Err(err).Msg("failed to get runs")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if len(ojsRuns) < 1 {
		return c.JSON(http.StatusOK, nil)
	}
	ojr := ojsRuns[0]
	log.Info().Interface("oj.OrchestrationName", ojr.OrchestrationName).Msg("ExportRunCsvRequest")
	ue := artemis_entities.UserEntity{
		Nickname: ojr.OrchestrationName,
		Platform: "csv-exports",
	}
	_, err = ai_platform_service_orchestrations.S3WfRunExport(c.Request().Context(), ou, ojr.OrchestrationName, &ue)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	j := FlowsActionsRequest{}
	var tmp []interface{}
	for _, v := range ue.MdSlice {
		err = json.Unmarshal(v.JsonData, &j)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, nil)
		}
		tmp = append(tmp, v)
	}
	c.Response().Header().Set("Content-Type", "text/csv")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.csv"`, ue.Nickname))
	scsv, err := payloadToCsvString(&j.FlowsCsvPayload)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	_, err = c.Response().Write([]byte(scsv))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return nil
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

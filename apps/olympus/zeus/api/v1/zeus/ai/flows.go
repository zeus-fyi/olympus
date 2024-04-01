package zeus_v1_ai

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	ai_platform_service_orchestrations "github.com/zeus-fyi/olympus/pkg/zeus/ai/orchestrations"
)

type FlowsActionsRequest struct {
	FlowsCsvPayload `json:",inline"`
	Stages          map[string]bool   `json:"stages"`
	CommandPrompts  map[string]string `json:"commandPrompts"`
}

type FlowsCsvPayload struct {
	ContactsCsv []map[string]string `json:"contentContactsCsv"`
	PromptsCsv  []map[string]string `json:"promptsCsv,omitempty"`
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

func FlowsExportCsvRequestHandler(c echo.Context) error {
	request := new(FlowsActionsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return ExportRunCsvRequest(c, *request)
}

// TODO: get and use param value, replace ue
// add wf results from s3

func ExportRunCsvRequest(c echo.Context, fa FlowsActionsRequest) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	var err error
	ue := artemis_entities.UserEntity{
		Nickname: "78f6d017bf9dd3d8d974024aef62ecefc7d281d8e0c857638022319d101cf99278f6d017bf9dd3d8d974024aef62ecefc7d281d8e0c857638022319d101cf992",
		Platform: "flows",
	}
	_, err = ai_platform_service_orchestrations.GetS3GlobalOrg(c.Request().Context(), ou, &ue)
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

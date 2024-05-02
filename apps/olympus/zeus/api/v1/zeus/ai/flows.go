package zeus_v1_ai

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

type FlowsActionsRequest struct {
	FlowsCsvPayload `json:",inline"`
	PreviewCount    int               `json:"previewCount"`
	Stages          map[string]bool   `json:"stages"`
	CommandPrompts  map[string]string `json:"commandPrompts"`
}

type FlowsCsvPayload struct {
	ContactsCsvFilename string              `json:"contactsCsvFilename"`
	ContactsCsvStr      string              `json:"contentContactsCsvStr"`
	ContactsCsv         []map[string]string `json:"contentContactsCsv,omitempty"`
	PromptsCsvStr       string              `json:"promptsCsvStr"`
	PromptsCsv          []map[string]string `json:"promptsCsv,omitempty,omitempty"`
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

package zeus_v1_ai

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	kronos_helix "github.com/zeus-fyi/olympus/pkg/kronos/helix"
)

type GetWorkflowsRequest struct {
}

func GetWorkflowsRequestHandler(c echo.Context) error {
	request := new(GetWorkflowsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.GetWorkflows(c)
}

func (w *GetWorkflowsRequest) GetWorkflows(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	ojs, err := artemis_orchestrations.SelectAiSystemOrchestrationsWithInstructionsByGroupType(c.Request().Context(), ou.OrgID, "ai", "workflows")
	if err != nil {
		log.Err(err).Msg("failed to get workflows")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, ojs)
}

type PostWorkflowsRequest struct {
	WorkflowName string                      `json:"workflowName"`
	StepSize     int                         `json:"stepSize"`
	StepSizeUnit string                      `json:"stepSizeUnit"`
	Models       []WorkflowModelInstructions `json:"models"`
}

type WorkflowModelInstructions struct {
	InstructionType       string `json:"instructionType"`
	Prompt                string `json:"prompt"`
	Model                 string `json:"model"`
	MaxTokens             int    `json:"maxTokens,omitempty"`
	TokenOverflowStrategy string `json:"tokenOverflowStrategy,omitempty"`
	CycleCount            int    `json:"cycleCount"`
}

func PostWorkflowsRequestHandler(c echo.Context) error {
	request := new(PostWorkflowsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateOrUpdateWorkflow(c)
}

func (w *PostWorkflowsRequest) CreateOrUpdateWorkflow(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if w.WorkflowName == "" || w.StepSize == 0 || w.StepSizeUnit == "" || len(w.Models) == 0 {
		return c.JSON(http.StatusBadRequest, nil)
	}
	inst := kronos_helix.Instructions{
		GroupName: "ai",
		Type:      "workflows",
		Alerts:    kronos_helix.AlertInstructions{},
		AiInstruction: kronos_helix.AiInstructions{
			WorkflowName: w.WorkflowName,
			StepSize:     w.StepSize,
			StepSizeUnit: w.StepSizeUnit,
			Map:          make(map[string]kronos_helix.AiInstruction),
		},
	}
	for _, m := range w.Models {
		if len(m.InstructionType) == 0 {
			return c.JSON(http.StatusBadRequest, nil)
		}
		inst.AiInstruction.Map[m.InstructionType] = kronos_helix.AiInstruction{
			InstructionType:       m.InstructionType,
			Prompt:                m.Prompt,
			Model:                 m.Model,
			MaxTokens:             m.MaxTokens,
			TokenOverflowStrategy: m.TokenOverflowStrategy,
			CycleCount:            m.CycleCount,
		}
	}
	b, err := json.Marshal(inst)
	if err != nil {
		log.Err(err).Msg("failed to marshal instructions")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(ou.OrgID, w.WorkflowName, "ai", "workflows")
	ojs, err := artemis_orchestrations.InsertOrchestration(c.Request().Context(), oj, b)
	if err != nil {
		log.Err(err).Msg("failed to insert workflow")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, ojs)
}

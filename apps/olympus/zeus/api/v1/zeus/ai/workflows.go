package zeus_v1_ai

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
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
	//var tmp []kronos_helix.AiModelInstruction
	//for i, _ := range ojs {
	//	var ins kronos_helix.Instructions
	//	err = json.Unmarshal([]byte(ojs[i].Instructions), &ins)
	//	if err != nil {
	//		log.Err(err).Msg("failed to unmarshal instructions")
	//		return c.JSON(http.StatusInternalServerError, nil)
	//	}
	//	if len(ins.AiInstruction.Map) > 0 {
	//		for _, v := range ins.AiInstruction.Map {
	//			tmp = append(tmp, v...)
	//		}
	//	}
	//}
	tasks, err := artemis_orchestrations.SelectTasks(c.Request().Context(), ou.OrgID)
	if err != nil {
		log.Err(err).Msg("failed to get tasks")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, AiWorkflowWrapper{
		Workflows: ojs,
		Tasks:     tasks,
	})
}

type AiWorkflowWrapper struct {
	Workflows []artemis_autogen_bases.Orchestrations `json:"workflows"`
	Tasks     []artemis_orchestrations.AITaskLibrary `json:"tasks"`
}

type PostWorkflowsRequest struct {
	WorkflowName         string               `json:"workflowName"`
	StepSize             int                  `json:"stepSize"`
	StepSizeUnit         string               `json:"stepSizeUnit"`
	Models               TaskMap              `json:"models"`
	AggregateSubTasksMap AggregateSubTasksMap `json:"aggregateSubTasksMap"`
}

type AggregateSubTasksMap map[int]map[int]bool
type TaskMap map[int]TaskModelInstructions

// TaskModelInstructions represents the equivalent of the TypeScript interface TaskModelInstructions
type TaskModelInstructions struct {
	TaskID                int    `json:"taskID"`
	Model                 string `json:"model"`
	TaskType              string `json:"taskType"`
	TaskGroup             string `json:"taskGroup"`
	TaskName              string `json:"taskName"`
	MaxTokens             int    `json:"maxTokens"`
	TokenOverflowStrategy string `json:"tokenOverflowStrategy"`
	Prompt                string `json:"prompt"`
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
	if w.WorkflowName == "" || len(w.Models) == 0 {
		return c.JSON(http.StatusBadRequest, nil)
	}
	inst := kronos_helix.Instructions{
		GroupName: "ai",
		Type:      "workflows",
		AiInstruction: kronos_helix.AiInstructions{
			WorkflowName: w.WorkflowName,
			Map:          make(map[string][]kronos_helix.AiModelInstruction),
		},
	}

	totalCycles := 0
	for _, m := range w.Models {
		if len(m.TaskType) == 0 {
			return c.JSON(http.StatusBadRequest, nil)
		}
		totalCycles += m.CycleCount
		if m.CycleCount == 0 {
			continue
		}
		switch m.TaskType {
		case "analysis", "aggregation":
			inst.AiInstruction.Map[m.TaskType] = append(inst.AiInstruction.Map[m.TaskType],
				kronos_helix.AiModelInstruction{
					Prompt:                m.Prompt,
					Model:                 m.Model,
					MaxTokens:             m.MaxTokens,
					TokenOverflowStrategy: m.TokenOverflowStrategy,
					CycleCount:            m.CycleCount,
				})
		default:
			return c.JSON(http.StatusBadRequest, nil)
		}
	}
	if totalCycles == 0 {
		return c.JSON(http.StatusBadRequest, nil)
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

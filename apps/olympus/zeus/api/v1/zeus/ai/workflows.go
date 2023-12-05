package zeus_v1_ai

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
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
	//ojs, err := artemis_orchestrations.SelectAiSystemOrchestrationsWithInstructionsByGroupType(c.Request().Context(), ou.OrgID, "ai", "workflows")
	//if err != nil {
	//	log.Err(err).Msg("failed to get workflows")
	//	return c.JSON(http.StatusInternalServerError, nil)
	//}

	ojs, err := artemis_orchestrations.SelectWorkflowTemplates(c.Request().Context(), ou)
	if err != nil {
		log.Err(err).Msg("failed to get workflows")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	ret, err := artemis_orchestrations.SelectRetrievals(c.Request().Context(), ou)
	if err != nil {
		log.Err(err).Msg("failed to get retrievals")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	tasks, err := artemis_orchestrations.SelectTasks(c.Request().Context(), ou.OrgID)
	if err != nil {
		log.Err(err).Msg("failed to get tasks")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, AiWorkflowWrapper{
		Workflows:  ojs,
		Tasks:      tasks,
		Retrievals: ret,
	})
}

type AiWorkflowWrapper struct {
	Workflows  []artemis_orchestrations.WorkflowTemplate `json:"workflows"`
	Tasks      []artemis_orchestrations.AITaskLibrary    `json:"tasks"`
	Retrievals []artemis_orchestrations.RetrievalItem    `json:"retrievals"`
}

type PostWorkflowsRequest struct {
	WorkflowName          string                `json:"workflowName"`
	WorkflowGroupName     string                `json:"workflowGroupName"`
	StepSize              int                   `json:"stepSize"`
	StepSizeUnit          string                `json:"stepSizeUnit"`
	Models                TaskMap               `json:"models"`
	AggregateSubTasksMap  AggregateSubTasksMap  `json:"aggregateSubTasksMap"`
	AnalysisRetrievalsMap AnalysisRetrievalsMap `json:"analysisRetrievalsMap"`
}

type AnalysisRetrievalsMap map[int]map[int]bool
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
	wt := artemis_orchestrations.WorkflowTemplate{
		WorkflowName:              w.WorkflowName,
		WorkflowGroup:             w.WorkflowGroupName,
		FundamentalPeriod:         w.StepSize,
		FundamentalPeriodTimeUnit: w.StepSizeUnit,
	}
	wft := artemis_orchestrations.WorkflowTasks{
		AggTasks:          []artemis_orchestrations.AggTask{},
		AnalysisOnlyTasks: []artemis_orchestrations.AITaskLibrary{},
	}

	ms := make(map[int]int)
	for _, m := range w.Models {
		if m.CycleCount < 1 {
			m.CycleCount = 1
		}
		switch m.TaskType {
		case "aggregation":
			agt := artemis_orchestrations.AggTask{
				AggId:      m.TaskID,
				CycleCount: m.CycleCount,
				Tasks:      []artemis_orchestrations.AITaskLibrary{},
			}
			for k, v := range w.AggregateSubTasksMap {
				if k == m.TaskID {
					for at, isTrue := range v {
						if isTrue {
							ms[at] = m.TaskID
							agt.Tasks = append(agt.Tasks, artemis_orchestrations.AITaskLibrary{
								TaskID:     at,
								CycleCount: m.CycleCount,
							})
						}
					}
				}
			}
			wft.AggTasks = append(wft.AggTasks, agt)
		case "analysis":
			at := artemis_orchestrations.AITaskLibrary{
				TaskID:                m.TaskID,
				OrgID:                 ou.OrgID,
				UserID:                ou.UserID,
				MaxTokensPerTask:      m.MaxTokens,
				TaskType:              m.TaskType,
				TaskName:              m.TaskName,
				TaskGroup:             m.TaskGroup,
				TokenOverflowStrategy: m.TokenOverflowStrategy,
				Model:                 m.Model,
				Prompt:                m.Prompt,
				CycleCount:            m.CycleCount,
				RetrievalDependencies: []artemis_orchestrations.RetrievalItem{},
			}
			for k, v := range w.AnalysisRetrievalsMap {
				for rt, isTrue := range v {
					if isTrue && rt == m.TaskID {
						at.RetrievalDependencies = append(at.RetrievalDependencies, artemis_orchestrations.RetrievalItem{
							RetrievalID: k,
						})
					}
				}
			}
			wft.AnalysisOnlyTasks = append(wft.AnalysisOnlyTasks, at)
		default:
			return c.JSON(http.StatusBadRequest, nil)
		}
	}

	aggSubTaskMap := make(map[int]artemis_orchestrations.AITaskLibrary)
	var analysisOnlyTasks []artemis_orchestrations.AITaskLibrary

	for _, v := range wft.AnalysisOnlyTasks {
		aggId, tok := ms[v.TaskID]
		if tok {
			aggSubTaskMap[aggId] = v
		} else {
			analysisOnlyTasks = append(analysisOnlyTasks, v)
		}
	}
	wft.AnalysisOnlyTasks = analysisOnlyTasks
	for _, v := range wft.AggTasks {
		sbt, tok := aggSubTaskMap[v.AggId]
		if tok {
			for ti, aa := range v.Tasks {
				if aa.TaskID == sbt.TaskID {
					v.Tasks[ti] = sbt
				}
			}
		}
	}
	err := artemis_orchestrations.InsertWorkflowWithComponents(c.Request().Context(), ou, &wt, wft)
	if err != nil {
		log.Err(err).Msg("failed to insert workflow")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, wt)
}

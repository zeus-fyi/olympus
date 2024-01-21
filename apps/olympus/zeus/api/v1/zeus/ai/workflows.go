package zeus_v1_ai

import (
	"net/http"
	"sort"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_openai_dbmodels "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
	hestia_stripe "github.com/zeus-fyi/olympus/pkg/hestia/stripe"
	ai_platform_service_orchestrations "github.com/zeus-fyi/olympus/pkg/zeus/ai/orchestrations"
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
	tasks, err := artemis_orchestrations.SelectTasks(c.Request().Context(), ou)
	if err != nil {
		log.Err(err).Msg("failed to get tasks")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	ojsRuns, err := artemis_orchestrations.SelectAiSystemOrchestrations(c.Request().Context(), ou)
	if err != nil {
		log.Err(err).Msg("failed to get runs")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	si, err := hera_openai_dbmodels.GetSearchIndexersByOrg(c.Request().Context(), ou)
	if err != nil {
		log.Err(err).Msg("failed to get search indexers")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	actions, err := artemis_orchestrations.SelectTriggerActionsByOrgAndOptParams(c.Request().Context(), ou, 0)
	if err != nil {
		log.Err(err).Msg("failed to get actions")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	schemas, err := artemis_orchestrations.SelectJsonSchemaByOrg(c.Request().Context(), ou)
	if err != nil {
		log.Err(err).Msg("failed to get actions")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	evals, err := artemis_orchestrations.SelectEvalFnsByOrgIDAndID(c.Request().Context(), ou, 0)
	if err != nil {
		log.Err(err).Msg("failed to get actions")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	var assistants []artemis_orchestrations.AiAssistant
	sv, err := ai_platform_service_orchestrations.GetMockingBirdSecrets(c.Request().Context(), ou)
	if err == nil && sv != nil && sv.ApiKey != "" {
		oc := hera_openai.InitOrgHeraOpenAI(sv.ApiKey)
		al, lerr := oc.ListAssistants(c.Request().Context(), nil, nil, nil, nil)
		if lerr != nil {
			log.Err(lerr).Msg("failed to get assistants")
			return c.JSON(http.StatusInternalServerError, nil)
		}
		for _, a := range al.Assistants {
			assistants = append(assistants, artemis_orchestrations.AiAssistant{
				Assistant: a,
			})
		}
	}
	//
	//assistants, err = artemis_orchestrations.SelectAssistants(c.Request().Context(), ou)
	//if err != nil {
	//	log.Err(err).Msg("failed to get assistants")
	//	return c.JSON(http.StatusInternalServerError, nil)
	//}
	sortWorkflowsByTemplateID(ojs.WorkflowTemplateSlice)
	return c.JSON(http.StatusOK, AiWorkflowWrapper{
		Workflows:      ojs.WorkflowTemplateSlice,
		Tasks:          tasks,
		Retrievals:     ret,
		Runs:           ojsRuns,
		SearchIndexers: si,
		TriggerActions: actions,
		Evals:          evals,
		Assistants:     assistants,
		Schemas:        schemas.Slice,
	})
}

func sortWorkflowsByTemplateID(workflows []artemis_orchestrations.WorkflowTemplateValue) {
	sort.Slice(workflows, func(i, j int) bool {
		return workflows[i].WorkflowTemplateID > workflows[j].WorkflowTemplateID
	})
}

type AiWorkflowWrapper struct {
	Workflows      []artemis_orchestrations.WorkflowTemplateValue  `json:"workflows"`
	Runs           []artemis_orchestrations.OrchestrationsAnalysis `json:"runs"`
	Tasks          []artemis_orchestrations.AITaskLibrary          `json:"tasks"`
	Retrievals     []artemis_orchestrations.RetrievalItem          `json:"retrievals"`
	SearchIndexers []hera_openai_dbmodels.SearchIndexerParams      `json:"searchIndexers"`
	Evals          []artemis_orchestrations.EvalFn                 `json:"evalFns"`
	TriggerActions []artemis_orchestrations.TriggerAction          `json:"triggerActions"`
	Assistants     []artemis_orchestrations.AiAssistant            `json:"assistants"`
	Schemas        []artemis_orchestrations.JsonSchemaDefinition   `json:"schemas"`
}

type PostWorkflowsRequest struct {
	WorkflowName          string                                `json:"workflowName"`
	WorkflowGroupName     string                                `json:"workflowGroupName"`
	StepSize              int                                   `json:"stepSize"`
	StepSizeUnit          string                                `json:"stepSizeUnit"`
	Models                TaskMap                               `json:"models"`
	EvalsMap              map[int]artemis_orchestrations.EvalFn `json:"evalsMap,omitempty"`
	EvalTasksMap          TaskEvalsMap                          `json:"evalTasksMap,omitempty"`
	AggregateSubTasksMap  AggregateSubTasksMap                  `json:"aggregateSubTasksMap,omitempty"`
	AnalysisRetrievalsMap AnalysisRetrievalsMap                 `json:"analysisRetrievalsMap"`
}

type TaskEvalsMap map[int]map[int]bool
type AnalysisRetrievalsMap map[int]map[int]bool
type AggregateSubTasksMap map[int]map[int]bool
type TaskMap map[int]TaskModelInstructions

// TaskModelInstructions represents the equivalent of the TypeScript interface TaskModelInstructions
type TaskModelInstructions struct {
	TaskID                int                             `json:"taskID"`
	Model                 string                          `json:"model"`
	TaskType              string                          `json:"taskType"`
	TaskGroup             string                          `json:"taskGroup"`
	TaskName              string                          `json:"taskName"`
	MaxTokens             int                             `json:"maxTokens"`
	TokenOverflowStrategy string                          `json:"tokenOverflowStrategy"`
	Prompt                string                          `json:"prompt"`
	CycleCount            int                             `json:"cycleCount"`
	EvalFns               []artemis_orchestrations.EvalFn `json:"evalFns,omitempty"`
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
	isBillingSetup, err := hestia_stripe.DoesUserHaveBillingMethod(c.Request().Context(), ou.UserID)
	if err != nil {
		log.Error().Err(err).Msg("failed to check if user has billing method")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if !isBillingSetup {
		return c.JSON(http.StatusPreconditionFailed, nil)
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
			if w.EvalTasksMap != nil {
				if evm, tok := w.EvalTasksMap[m.TaskID]; tok {
					for k, v := range evm {
						if v {
							agt.EvalFns = append(agt.EvalFns, w.EvalsMap[k])
						}
					}
				}
			}
			for k, v := range w.AggregateSubTasksMap {
				if k == m.TaskID {
					for at, isTrue := range v {
						if isTrue {
							ms[at] = m.TaskID

							ta := artemis_orchestrations.AITaskLibrary{
								TaskID:     at,
								CycleCount: m.CycleCount,
							}
							if w.EvalTasksMap != nil {
								if evm, tok := w.EvalTasksMap[at]; tok {
									for ke, ve := range evm {
										if ve {
											agt.EvalFns = append(agt.EvalFns, w.EvalsMap[ke])
										}
									}
								}
							}
							agt.Tasks = append(agt.Tasks, ta)
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

			if w.EvalTasksMap != nil {
				if evm, tok := w.EvalTasksMap[m.TaskID]; tok {
					for ke, ve := range evm {
						if ve {
							at.EvalFns = append(at.EvalFns, w.EvalsMap[ke])
						}
					}
				}
			}
			for k, v := range w.AnalysisRetrievalsMap {
				for rt, isTrue := range v {
					if isTrue && rt == m.TaskID && k > 0 {
						at.RetrievalDependencies = append(at.RetrievalDependencies, artemis_orchestrations.RetrievalItem{
							RetrievalID: &k,
						})
					}
				}
			}
			wft.AnalysisOnlyTasks = append(wft.AnalysisOnlyTasks, at)
		default:
			return c.JSON(http.StatusBadRequest, nil)
		}
	}
	err = artemis_orchestrations.InsertWorkflowWithComponents(c.Request().Context(), ou, &wt, wft)
	if err != nil {
		log.Err(err).Msg("failed to insert workflow")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, wt)
}

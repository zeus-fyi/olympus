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
	if ojs.WorkflowTemplateSlice == nil {
		ojs.WorkflowTemplateSlice = []artemis_orchestrations.WorkflowTemplateValue{}
	}
	ret, err := artemis_orchestrations.SelectRetrievals(c.Request().Context(), ou, 0)
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
		log.Err(err).Interface("ou", ou).Msg("GetWorkflowsRequest: failed to get search indexers")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	tp := artemis_orchestrations.TriggersWorkflowQueryParams{Ou: ou}
	actions, err := artemis_orchestrations.SelectTriggerActionsByOrgAndOptParams(c.Request().Context(), tp)
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
		} else {
			for _, a := range al.Assistants {
				assistants = append(assistants, artemis_orchestrations.AiAssistant{
					Assistant: a,
				})
			}
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
	Schemas        []*artemis_orchestrations.JsonSchemaDefinition  `json:"schemas"`
}

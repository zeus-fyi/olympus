package zeus_v1_ai

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_stripe "github.com/zeus-fyi/olympus/pkg/hestia/stripe"
	ai_platform_service_orchestrations "github.com/zeus-fyi/olympus/pkg/zeus/ai/orchestrations"
)

type WorkflowsActionsRequest struct {
	Action        string `json:"action"`
	UnixStartTime int    `json:"unixStartTime,omitempty"`
	Duration      int    `json:"duration,omitempty"`
	DurationUnit  string `json:"durationUnit,omitempty"`

	IsStrictTimeWindow           bool   `json:"isStrictTimeWindow,omitempty"`
	CustomBasePeriod             bool   `json:"customBasePeriod,omitempty"`
	CustomBasePeriodStepSize     int    `json:"customBasePeriodStepSize,omitempty"`
	CustomBasePeriodStepSizeUnit string `json:"customBasePeriodStepSizeUnit,omitempty"`

	WfSchemaFieldOverrides artemis_orchestrations.WorkflowSchemaOverrides       `json:"wfSchemaFieldOverrides,omitempty"`
	WfRetrievalOverrides   map[string]artemis_orchestrations.RetrievalOverrides `json:"wfRetrievalPayloadOverrides,omitempty"`
	WfTaskOverrides        map[string]artemis_orchestrations.TaskOverrides      `json:"wfTaskOverrides,omitempty"`

	TaskOverrides             artemis_orchestrations.TaskOverrides                 `json:"taskOverrides,omitempty"`
	SchemaFieldOverrides      artemis_orchestrations.SchemaOverrides               `json:"schemaFieldOverrides,omitempty"`
	WorkflowEntityRefs        []artemis_entities.EntitiesFilter                    `json:"workflowEntitiesRef,omitempty"`
	WorkflowEntities          []artemis_entities.UserEntity                        `json:"workflowEntities,omitempty"`
	WorkflowEntitiesOverrides artemis_orchestrations.WorkflowUserEntitiesOverrides `json:"workflowEntitiesOverrides,omitempty"`
	Workflows                 []artemis_orchestrations.WorkflowTemplate            `json:"workflows,omitempty"`
}

func WorkflowsActionsRequestHandler(c echo.Context) error {
	request := new(WorkflowsActionsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.Process(c)
}

func (w *WorkflowsActionsRequest) Process(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Info().Interface("ou", ou)
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

	var rid int
	if w.CustomBasePeriod && w.CustomBasePeriodStepSize > 0 && w.CustomBasePeriodStepSizeUnit != "" {
		for i, _ := range w.Workflows {
			w.Workflows[i].FundamentalPeriod = w.CustomBasePeriodStepSize
			w.Workflows[i].FundamentalPeriodTimeUnit = w.CustomBasePeriodStepSizeUnit
		}
	}
	switch w.Action {
	case "start":
		window, isCycleStepped := w.GetTimeSeriesIterInst()
		err = w.ConvertWfStrIDs()
		if err != nil {
			log.Err(err).Interface("ou", ou).Interface("[]WorkflowTemplate", w.Workflows).Msg("WorkflowsActionsRequestHandler: ConvertWfStrIDs failed")
			return c.JSON(http.StatusInternalServerError, nil)
		}
		resp, rerr := artemis_orchestrations.GetAiOrchestrationParams(c.Request().Context(), ou, &window, w.Workflows)
		if rerr != nil {
			log.Err(rerr).Interface("ou", ou).Interface("[]WorkflowTemplate", w.Workflows).Msg("WorkflowsActionsRequestHandler: GetAiOrchestrationParams failed")
			return c.JSON(http.StatusInternalServerError, nil)
		}
		for ri, _ := range resp {
			resp[ri].WorkflowExecTimekeepingParams.IsCycleStepped = isCycleStepped
			if isCycleStepped {
				resp[ri].WorkflowExecTimekeepingParams.RunCycles = w.Duration
			}
			//for ti, task := range resp[ri].WorkflowTasks {
			//	tov, tok := w.TaskOverrides[task.AnalysisTaskName]
			//	if tok {
			//		resp[ri].WorkflowTasks[ti].AnalysisPrompt = tov.ReplacePrompt
			//	}
			//	if task.AggTaskID != nil && *task.AggTaskID > 0 {
			//		tov, tok = w.TaskOverrides[strconv.Itoa(*task.AggTaskID)]
			//		if tok {
			//			resp[ri].WorkflowTasks[ti].AggPrompt = aws.String(tov.ReplacePrompt)
			//		}
			//	}
			//}
			if w.SchemaFieldOverrides != nil {
				resp[ri].WorkflowOverrides.SchemaFieldOverrides = w.SchemaFieldOverrides
			}
			resp[ri].WorkflowExecTimekeepingParams.IsStrictTimeWindow = w.IsStrictTimeWindow
			rid, err = ai_platform_service_orchestrations.ZeusAiPlatformWorker.ExecuteRunAiWorkflowProcess(c.Request().Context(), ou, resp[ri])
			if err != nil {
				log.Err(err).Interface("ou", ou).Interface("WorkflowExecParams", resp).Msg("WorkflowsActionsRequestHandler: ExecuteRunAiWorkflowProcess failed")
				return c.JSON(http.StatusInternalServerError, nil)
			}
		}
	case "stop":
		// do y
	}
	return c.JSON(http.StatusOK, fmt.Sprintf("%d", rid))
}

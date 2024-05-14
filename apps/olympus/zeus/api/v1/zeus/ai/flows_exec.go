package zeus_v1_ai

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	ai_platform_service_orchestrations "github.com/zeus-fyi/olympus/pkg/zeus/ai/orchestrations"
)

type ExecFlowsActionsRequest struct {
	WorkflowsActionsRequest `json:",inline"`
	FlowsActionsRequest     `json:",inline"`
}

func FlowsExecActionsRequestHandler(c echo.Context) error {
	request := new(ExecFlowsActionsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.ProcessFlow(c)
}

func (w *ExecFlowsActionsRequest) ProcessFlow(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Info().Interface("ou", ou)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	_, err := w.SetupFlow(c.Request().Context(), ou)
	if err != nil {
		log.Err(err).Interface("w", w).Msg("SaveImport failed")
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if len(w.Workflows) > 0 || len(w.CustomWorkflows) > 0 {
		w.Action = "start"
	}
	w.Duration = 1
	w.DurationUnit = "cycles"
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
		tmpOu := ou
		tmpOu.OrgID = 1685378241971196000
		resp, rerr := artemis_orchestrations.GetAiOrchestrationParams(c.Request().Context(), tmpOu, &window, w.Workflows)
		if rerr != nil {
			log.Err(rerr).Interface("ou", ou).Interface("[]WorkflowTemplate", w.Workflows).Msg("WorkflowsActionsRequestHandler: GetAiOrchestrationParams failed")
			return c.JSON(http.StatusInternalServerError, nil)
		}
		resp2, rerr := artemis_orchestrations.GetAiOrchestrationParams(c.Request().Context(), ou, &window, w.CustomWorkflows)
		if rerr != nil {
			log.Err(rerr).Interface("ou", ou).Interface("[]WorkflowTemplate", w.Workflows).Msg("WorkflowsActionsRequestHandler: GetAiOrchestrationParams failed")
			return c.JSON(http.StatusInternalServerError, nil)
		}
		resp = append(resp, resp2...)
		var wfExecParams artemis_orchestrations.WorkflowExecParams
		for ri, _ := range resp {
			wfName := resp[ri].WorkflowTemplate.WorkflowName
			if ve, wok := w.WorkflowEntitiesOverrides[wfName]; wok {
				resp[ri].WorkflowOverrides.WorkflowEntities = ve
			} else {
				resp[ri].WorkflowOverrides.WorkflowEntities = w.WorkflowEntities
			}
			resp[ri].WorkflowExecTimekeepingParams.IsCycleStepped = isCycleStepped
			if isCycleStepped {
				resp[ri].WorkflowExecTimekeepingParams.RunCycles = w.Duration
			}
			//for ti, task := range resp[ri].WorkflowTasks {
			//tov, tok := w.TaskOverrides[task.AnalysisTaskName]
			//if tok {
			//	resp[ri].WorkflowTasks[ti].AnalysisPrompt = tov.ReplacePrompt
			//}
			//if task.AggTaskID != nil && *task.AggTaskID > 0 {
			//	tov, tok = w.TaskOverrides[strconv.Itoa(*task.AggTaskID)]
			//	if tok {
			//		resp[ri].WorkflowTasks[ti].AggPrompt = aws.String(tov.ReplacePrompt)
			//	}
			//}
			//}
			if w.WfRetrievalOverrides != nil {
				if wv, wok := w.WfRetrievalOverrides[wfName]; wok {
					resp[ri].WorkflowOverrides.RetrievalOverrides = wv
				}
			}
			if w.WfSchemaFieldOverrides != nil {
				if wv, wok := w.WfSchemaFieldOverrides[wfName]; wok {
					resp[ri].WorkflowOverrides.SchemaFieldOverrides = wv
				}
			}
			if w.WfTaskOverrides != nil {
				if wv, wok := w.WfTaskOverrides[wfName]; wok {
					resp[ri].WorkflowOverrides.TaskPromptOverrides = wv
				}
			}
			resp[ri].WorkflowExecTimekeepingParams.IsStrictTimeWindow = w.IsStrictTimeWindow
			resp[ri].WorkflowOverrides.IsUsingFlows = true
			resp[ri].WorkflowOverrides.WorkflowEntityRefs = w.WorkflowEntityRefs
			wfExecParams.WorkflowExecTimekeepingParams = resp[ri].WorkflowExecTimekeepingParams
			wfExecParams.WorkflowTasks = append(wfExecParams.WorkflowTasks, resp[ri].WorkflowTasks...)
			wfExecParams.MergeWorkflowTaskRelationships(resp[ri].WorkflowTaskRelationships)
			wfExecParams.MergeCycleCountTaskRelative(resp[ri].CycleCountTaskRelative)
			wfExecParams.MergeWorkflowOverrides(resp[ri].WorkflowOverrides)
		}

		wfExecParams.WorkflowOverrides.IsUsingFlows = true
		wfExecParams.WorkflowTemplate.WorkflowName = "csv-analysis"
		wfExecParams.WorkflowTemplate.WorkflowGroup = w.ContactsCsvFilename
		rid, err = ai_platform_service_orchestrations.ZeusAiPlatformWorker.ExecuteRunAiWorkflowProcess(c.Request().Context(), ou, wfExecParams)
		if err != nil {
			log.Err(err).Interface("ou", ou).Interface("WorkflowExecParams", resp).Msg("WorkflowsActionsRequestHandler: ExecuteRunAiWorkflowProcess failed")
			return c.JSON(http.StatusInternalServerError, nil)
		}

	case "stop":
		// do y
	}
	return c.JSON(http.StatusOK, fmt.Sprintf("%d", rid))
}

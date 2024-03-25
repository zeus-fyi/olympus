package zeus_v1_ai

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
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

const (
	googWf = "google-query-regex-index-wf"
	liWf   = "linkedin-rapid-api-profiles-wf"
)

func (w *ExecFlowsActionsRequest) GoogleSearchSetup() error {
	if v, ok := w.Stages["googleSearch"]; !ok || !v {
		return nil
	}
	b, err := json.Marshal(w.ContactsCsv)
	if err != nil {
		log.Err(err).Msg("failed to marshal gs")
		return err
	}
	if w.TaskOverrides == nil {
		w.TaskOverrides = make(map[string]TaskOverride)
	}
	w.TaskOverrides["zeusfyi-verbatim"] = TaskOverride{ReplacePrompt: string(b)}
	if v, ok := w.CommandPrompts["googleSearch"]; ok && v != "" {
		if w.SchemaFieldOverrides == nil {
			w.SchemaFieldOverrides = make(map[string]map[string]string)
			w.SchemaFieldOverrides["google-results-agg"] = map[string]string{
				"summary": v,
			}
		}
	}
	w.Workflows = append(w.Workflows, artemis_orchestrations.WorkflowTemplate{
		WorkflowName: googWf,
	})
	return nil
}

// Can you tell me what this person does in their current role; and the company they work at now?

func (w *ExecFlowsActionsRequest) LinkedInScraperSetup() error {
	if v, ok := w.Stages["linkedIn"]; !ok || !v {
		return nil
	}
	b, err := json.Marshal(w.ContactsCsv)
	if err != nil {
		log.Err(err).Msg("failed to marshal li")
		return err
	}
	if w.TaskOverrides == nil {
		w.TaskOverrides = make(map[string]TaskOverride)
	}
	w.TaskOverrides["linkedin-profiles-rapid-api-qps"] = TaskOverride{ReplacePrompt: string(b)}
	if v, ok := w.CommandPrompts["linkedIn"]; ok && v != "" {
		if w.SchemaFieldOverrides == nil {
			w.SchemaFieldOverrides = make(map[string]map[string]string)
			w.SchemaFieldOverrides["results-agg"] = map[string]string{
				"summary": v,
			}
		}
	}
	w.Workflows = append(w.Workflows, artemis_orchestrations.WorkflowTemplate{
		WorkflowName: liWf,
	})
	return nil
}

func (w *ExecFlowsActionsRequest) LinkedInBizScraperSetup() error {
	if v, ok := w.Stages["linkedInBiz"]; !ok || !v {
		return nil
	}
	b, err := json.Marshal(w.ContactsCsv)
	if err != nil {
		log.Err(err).Msg("failed to marshal li")
		return err
	}
	if w.TaskOverrides == nil {
		w.TaskOverrides = make(map[string]TaskOverride)
	}
	w.TaskOverrides["linkedin-biz-profiles-rapid-api-qps"] = TaskOverride{ReplacePrompt: string(b)}
	if v, ok := w.CommandPrompts["linkedIn"]; ok && v != "" {
		if w.SchemaFieldOverrides == nil {
			w.SchemaFieldOverrides = make(map[string]map[string]string)
			w.SchemaFieldOverrides["results-agg"] = map[string]string{
				"summary": v,
			}
		}
	}
	w.Workflows = append(w.Workflows, artemis_orchestrations.WorkflowTemplate{
		WorkflowName: liWf,
	})
	return nil
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
	err := w.GoogleSearchSetup()
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	err = w.LinkedInScraperSetup()
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	if len(w.Workflows) > 0 {
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
		if w.UnixStartTime == 0 {
			w.UnixStartTime = int(time.Now().Unix())
		}
		isCycleStepped := false
		endTime := 0
		switch w.DurationUnit {
		case "minutes", "minute":
			minutes := time.Minute
			if w.Duration < 0 {
				w.UnixStartTime += int(minutes.Seconds()) * w.Duration
				w.Duration = -1 * w.Duration
			}
			endTime = w.UnixStartTime + int(minutes.Seconds())*w.Duration
		case "hours", "hour", "hrs", "hr":
			hours := time.Hour
			if w.Duration < 0 {
				w.UnixStartTime += int(hours.Seconds()) * w.Duration
				w.Duration = -1 * w.Duration
			}
			endTime = w.UnixStartTime + int(hours.Seconds())*w.Duration
		case "days", "day":
			days := 24 * time.Hour
			if w.Duration < 0 {
				w.UnixStartTime += int(days.Seconds()) * w.Duration
				w.Duration = -1 * w.Duration
			}
			endTime = w.UnixStartTime + int(days.Seconds())*w.Duration
		case "weeks", "week":
			weeks := 7 * 24 * time.Hour
			if w.Duration < 0 {
				w.UnixStartTime += int(weeks.Seconds()) * w.Duration
				w.Duration = -1 * w.Duration
			}
			endTime = w.UnixStartTime + int(weeks.Seconds())*w.Duration
		case "cycles":
			isCycleStepped = true
		}
		window := artemis_orchestrations.Window{
			UnixStartTime: w.UnixStartTime,
			UnixEndTime:   endTime,
		}
		for wfi, _ := range w.Workflows {
			if w.Workflows[wfi].WorkflowTemplateStrID != "" {
				wid, werr := strconv.Atoi(w.Workflows[wfi].WorkflowTemplateStrID)
				if werr != nil {
					log.Err(werr).Msg("failed to parse int")
					return c.JSON(http.StatusBadRequest, nil)
				}
				w.Workflows[wfi].WorkflowTemplateID = wid
			}
		}
		for _, wfn := range w.Workflows {
			switch wfn.WorkflowName {
			case googWf:
			case liWf:
			default:
				return c.JSON(http.StatusBadRequest, nil)
			}
		}
		tmpOu := ou
		tmpOu.OrgID = 1685378241971196000
		resp, rerr := artemis_orchestrations.GetAiOrchestrationParams(c.Request().Context(), tmpOu, &window, w.Workflows)
		if rerr != nil {
			log.Err(rerr).Interface("ou", ou).Interface("[]WorkflowTemplate", w.Workflows).Msg("WorkflowsActionsRequestHandler: GetAiOrchestrationParams failed")
			return c.JSON(http.StatusInternalServerError, nil)
		}
		for ri, _ := range resp {
			resp[ri].WorkflowExecTimekeepingParams.IsCycleStepped = isCycleStepped
			if isCycleStepped {
				resp[ri].WorkflowExecTimekeepingParams.RunCycles = w.Duration
			}
			for ti, task := range resp[ri].WorkflowTasks {
				tov, tok := w.TaskOverrides[task.AnalysisTaskName]
				if tok {
					resp[ri].WorkflowTasks[ti].AnalysisPrompt = tov.ReplacePrompt
				}
				if task.AggTaskID != nil && *task.AggTaskID > 0 {
					tov, tok = w.TaskOverrides[strconv.Itoa(*task.AggTaskID)]
					if tok {
						resp[ri].WorkflowTasks[ti].AggPrompt = aws.String(tov.ReplacePrompt)
					}
				}
			}
			if w.SchemaFieldOverrides != nil {
				resp[ri].WorkflowOverrides.SchemaFieldOverrides = w.SchemaFieldOverrides
			}
			resp[ri].WorkflowExecTimekeepingParams.IsStrictTimeWindow = w.IsStrictTimeWindow
			resp[ri].WorkflowOverrides.IsUsingFlows = true
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

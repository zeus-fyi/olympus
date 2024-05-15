package zeus_v1_ai

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func GetRunReportsRequestHandler(c echo.Context) error {
	request := new(RunsActionsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.GetRuns(c)
}

func GetUIRunReportsRequestHandler(c echo.Context) error {
	request := new(RunsActionsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.GetRunsUI(c)
}

func (w *RunsActionsRequest) GetRunsUI(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Info().Interface("ou", ou)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	//isBillingSetup, err := hestia_stripe.DoesUserHaveBillingMethod(c.Request().Context(), ou.UserID)
	//if err != nil {
	//	log.Error().Err(err).Msg("failed to check if user has billing method")
	//	return c.JSON(http.StatusInternalServerError, nil)
	//}
	//if !isBillingSetup {
	//	return c.JSON(http.StatusPreconditionFailed, nil)
	//}
	ojsRuns, err := artemis_orchestrations.SelectAiSystemOrchestrationsUI(context.Background(), ou, 0)
	if err != nil {
		log.Err(err).Msg("failed to get runs")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, ojsRuns)
}

func (w *RunsActionsRequest) GetRuns(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Info().Interface("ou", ou)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	//isBillingSetup, err := hestia_stripe.DoesUserHaveBillingMethod(c.Request().Context(), ou.UserID)
	//if err != nil {
	//	log.Error().Err(err).Msg("failed to check if user has billing method")
	//	return c.JSON(http.StatusInternalServerError, nil)
	//}
	//if !isBillingSetup {
	//	return c.JSON(http.StatusPreconditionFailed, nil)
	//}
	ojsRuns, err := artemis_orchestrations.SelectAiSystemOrchestrations(context.Background(), ou, 0)
	if err != nil {
		log.Err(err).Msg("failed to get runs")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, ojsRuns)
}

type GetRunsActionsRequest struct {
}

func GetRunActionsRequestHandler(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Err(err).Msg("invalid ID parameter")
		return c.JSON(http.StatusBadRequest, "invalid ID parameter")
	}
	request := new(GetRunsActionsRequest)
	if err = c.Bind(request); err != nil {
		return err
	}
	return request.GetRun(c, id) // Pass the ID to the method
}

func (w *GetRunsActionsRequest) GetRun(c echo.Context, id int) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Info().Interface("ou", ou)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	//isBillingSetup, err := hestia_stripe.DoesUserHaveBillingMethod(c.Request().Context(), ou.UserID)
	//if err != nil {
	//	log.Error().Err(err).Msg("failed to check if user has billing method")
	//	return c.JSON(http.StatusInternalServerError, nil)
	//}
	//if !isBillingSetup {
	//	return c.JSON(http.StatusPreconditionFailed, nil)
	//}
	ojsRuns, err := artemis_orchestrations.SelectAiSystemOrchestrations(context.Background(), ou, id)
	if err != nil {
		log.Err(err).Msg("failed to get runs")
		return c.JSON(http.StatusInternalServerError, nil)
	}

	for oi, ojv := range ojsRuns {
		pcUsed := 0
		cmpUsed := 0
		tokensUsed := 0
		lv := len(ojv.AggregatedData)
		if lv >= 10 {
			var tmp []artemis_orchestrations.AggregatedData
			for iv, v := range ojv.AggregatedData {
				if iv < 10 {
					tmp = append(tmp, v)
				}
				tokensUsed += v.TotalTokens
				pcUsed += v.PromptTokens
				cmpUsed += v.CompletionTokens
			}
			agd := artemis_orchestrations.AggregatedData{
				AIWorkflowAnalysisResult: artemis_orchestrations.AIWorkflowAnalysisResult{},
				TaskName:                 "task completions",
				TaskType:                 fmt.Sprintf("success %d", lv),
				PromptTokens:             pcUsed,
				CompletionTokens:         cmpUsed,
				TotalTokens:              tokensUsed,
			}
			tmp = append(tmp, agd)
			ojsRuns[oi].AggregatedData = tmp
		}
		rvl := len(ojv.AggregatedRetrievalResults)
		if rvl >= 10 {
			errCount := 0
			success := 0
			unknown := 0
			var tmp []artemis_orchestrations.AIWorkflowRetrievalResult
			var agr artemis_orchestrations.AIWorkflowRetrievalResult
			for iv, v := range ojv.AggregatedRetrievalResults {
				if iv == 0 {
					agr = v
				}
				if iv < 10 {
					tmp = append(tmp, v)
				}
				if strings.Contains(v.Status, "complete") {
					success += 1
				} else if strings.Contains(v.Status, "error") {
					errCount += 1
				} else {
					unknown += 1
				}
			}
			agr.RunningCycleNumber = 0
			agr.IterationCount = 0
			agr.ChunkOffset = 0
			agr.RetrievalName = "summary retrievals"
			agr.Status = fmt.Sprintf("success %d, errors %d, unknown %d, total %d", success, errCount, unknown, rvl)
			tmp = append(tmp, agr)
			ojsRuns[oi].AggregatedRetrievalResults = tmp
		}
	}
	return c.JSON(http.StatusOK, ojsRuns)
}

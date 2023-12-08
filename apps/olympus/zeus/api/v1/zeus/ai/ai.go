package zeus_v1_ai

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	ai_platform_service_orchestrations "github.com/zeus-fyi/olympus/pkg/zeus/ai/orchestrations"
)

const (
	internalOrgID = 7138983863666903883
)

func AiV1Routes(e *echo.Group) *echo.Group {
	e.POST("/search", AiSearchRequestHandler)
	e.POST("/search/analyze", AiSearchAnalyzeRequestHandler)
	e.GET("/workflows/ai", GetWorkflowsRequestHandler)
	e.POST("/workflows/ai", PostWorkflowsRequestHandler)
	e.POST("/workflows/ai/start", WorkflowsActionsRequestHandler)
	e.POST("/tasks/ai", CreateOrUpdateTaskRequestHandler)
	e.POST("/retrievals/ai", CreateOrUpdateRetrievalRequestHandler)

	// destructive
	e.DELETE("/workflows/ai", WorkflowsDeletionRequestHandler)
	return e
}

type AiSearchRequest struct {
	hera_search.AiSearchParams `json:"searchParams"`
}

func AiSearchRequestHandler(c echo.Context) error {
	request := new(AiSearchRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.Search(c)
}

func (r *AiSearchRequest) Search(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if ou.OrgID != internalOrgID {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	ts := time.Now()
	switch r.TimeRange {
	case "1 hour":
		r.SearchInterval[0] = ts.Add(-1 * time.Hour)
		r.SearchInterval[1] = ts
	case "24 hours":
		r.SearchInterval[0] = ts.AddDate(0, 0, -1)
		r.SearchInterval[1] = ts
	case "7 days":
		r.SearchInterval[0] = ts.AddDate(0, 0, -7)
		r.SearchInterval[1] = ts
	case "30 days":
		r.SearchInterval[0] = ts.AddDate(0, 0, -30)
		r.SearchInterval[1] = ts
	case "all":
		r.SearchInterval[0] = ts.AddDate(-4, 0, 0)
		r.SearchInterval[1] = ts
	case "window":
		log.Info().Interface("searchInterval", r.SearchInterval).Msg("window")
	}
	res, err := hera_search.PerformPlatformSearches(c.Request().Context(), ou, r.AiSearchParams)
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg("error performing platform searches")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	hera_search.SortSearchResults(res)
	if len(r.Retrieval.RetrievalPrompt) > 0 {
		aiResp, arr := ai_platform_service_orchestrations.AiAggregateTask(c.Request().Context(), ou, res, r.AiSearchParams)
		if arr != nil {
			log.Err(arr).Msg("error aggregating tasks")
			return c.JSON(http.StatusInternalServerError, nil)
		}
		fullResp := fmt.Sprintf("%s\n%s\n", aiResp.Choices[0].Message.Content, hera_search.FormatSearchResultsV2(res))
		return c.JSON(http.StatusOK, fullResp)
	}
	return c.JSON(http.StatusOK, hera_search.FormatSearchResultsV2(res))
}

func AiSearchAnalyzeRequestHandler(c echo.Context) error {
	request := new(AiSearchRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.SearchAnalyze(c)
}

func (r *AiSearchRequest) SearchAnalyze(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if ou.OrgID != internalOrgID {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	res, err := hera_search.PerformPlatformSearches(c.Request().Context(), ou, r.AiSearchParams)
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg("error performing platform searches")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	hera_search.SortSearchResults(res)
	resp, err := ai_platform_service_orchestrations.AiAggregateTask(c.Request().Context(), ou, res, r.AiSearchParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, resp.Choices[0].Message.Content)
}

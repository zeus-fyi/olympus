package zeus_v1_ai

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_stripe "github.com/zeus-fyi/olympus/pkg/hestia/stripe"
	ai_platform_service_orchestrations "github.com/zeus-fyi/olympus/pkg/zeus/ai/orchestrations"
)

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

	ts := time.Now()
	switch r.TimeRange {
	case "1 hour":
		r.Window.Start = ts.Add(-1 * time.Hour)
		r.Window.End = ts
	case "24 hours":
		r.Window.Start = ts.AddDate(0, 0, -1)
		r.Window.End = ts
	case "7 days":
		r.Window.Start = ts.AddDate(0, 0, -7)
		r.Window.End = ts
	case "30 days":
		r.Window.Start = ts.AddDate(0, 0, -30)
		r.Window.End = ts
	case "all":
		r.Window.Start = ts.AddDate(-4, 0, 0)
		r.Window.End = ts
	case "window":
		log.Info().Interface("searchInterval", r.Window).Msg("window")
	}
	res, err := hera_search.PerformPlatformSearches(c.Request().Context(), ou, r.AiSearchParams)
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg("error performing platform searches")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	hera_search.SortSearchResults(res)
	if len(r.Retrieval.RetrievalPrompt) > 0 {
		isBillingSetup, berr := hestia_stripe.DoesUserHaveBillingMethod(c.Request().Context(), ou.UserID)
		if berr != nil {
			log.Error().Err(berr).Msg("failed to check if user has billing method")
			return c.JSON(http.StatusInternalServerError, nil)
		}
		if !isBillingSetup {
			return c.JSON(http.StatusPreconditionFailed, nil)
		}
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
	isBillingSetup, berr := hestia_stripe.DoesUserHaveBillingMethod(c.Request().Context(), ou.UserID)
	if berr != nil {
		log.Error().Err(berr).Msg("failed to check if user has billing method")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if !isBillingSetup {
		return c.JSON(http.StatusPreconditionFailed, nil)
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

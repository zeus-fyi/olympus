package zeus_v1_ai

import (
	"net/http"

	"github.com/labstack/echo/v4"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	ai_platform_service_orchestrations "github.com/zeus-fyi/olympus/pkg/zeus/ai/ai/orchestrations"
)

const (
	internalOrgID = 7138983863666903883
)

func AiV1Routes(e *echo.Group) *echo.Group {
	e.POST("/search", AiSearchRequestHandler)
	e.POST("/search/analyze", AiSearchAnalyzeRequestHandler)
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
	res, err := hera_search.SearchTelegram(c.Request().Context(), ou, r.AiSearchParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, hera_search.FormatSearchResults(res))
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
	res, err := hera_search.SearchTelegram(c.Request().Context(), ou, r.AiSearchParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp, err := ai_platform_service_orchestrations.AiTelegramTask(c.Request().Context(), ou, res, r.AiSearchParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, resp.Choices[0].Message.Content)
}

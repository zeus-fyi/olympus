package zeus_v1_ai

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

const (
	internalOrgID = 7138983863666903883
)

func AiV1Routes(e *echo.Group) *echo.Group {
	e.POST("/search", AiSearchRequestHandler)
	e.POST("/search/analyze", AiSearchRequestHandler)
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

	return c.JSON(http.StatusOK, formatSearchResults(res))
}
func formatSearchResults(results []hera_search.SearchResult) string {
	var builder strings.Builder

	for _, result := range results {
		line := fmt.Sprintf("%d | %s | %s | %s | %s \n",
			result.UnixTimestamp,
			escapeString(result.Source),
			escapeString(result.Group),
			escapeString(result.Metadata.Username),
			escapeString(result.Value))
		builder.WriteString(line)
	}

	return builder.String()
}

func escapeString(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "\\u003c", "<"), "\\u003e", ">")
}
func AiSearchAnalyzeRequestHandler(c echo.Context) error {
	request := new(AiSearchRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.Search(c)
}

func (r *AiSearchRequest) SearchAnalyze(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	if ou.OrgID != internalOrgID {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, nil)
}

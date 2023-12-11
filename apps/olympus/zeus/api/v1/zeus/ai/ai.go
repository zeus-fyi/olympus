package zeus_v1_ai

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_stripe "github.com/zeus-fyi/olympus/pkg/hestia/stripe"
)

func AiV1Routes(e *echo.Group) *echo.Group {
	e.POST("/search", AiSearchRequestHandler)
	e.POST("/search/analyze", AiSearchAnalyzeRequestHandler)
	e.POST("/search/indexer", AiSearchIndexerRequestHandler)

	e.GET("/workflows/ai", GetWorkflowsRequestHandler)
	e.POST("/workflows/ai", PostWorkflowsRequestHandler)
	e.POST("/workflows/ai/actions", WorkflowsActionsRequestHandler)
	e.POST("/runs/ai/actions", RunsActionsRequestHandler)

	e.POST("/tasks/ai", CreateOrUpdateTaskRequestHandler)
	e.POST("/retrievals/ai", CreateOrUpdateRetrievalRequestHandler)

	// destructive
	e.DELETE("/workflows/ai", WorkflowsDeletionRequestHandler)
	return e
}

type AiSearchIndexerRequest struct {
	hera_search.SearchIndexerParams `json:"searchIndexer"`
	PlatformSecretReference         `json:"platformSecretReference"`
}

type PlatformSecretReference struct {
	SecretGroupName string `json:"secretGroupName"`
	SecretKeyName   string `json:"secretKeyName,omitempty"`
}

func AiSearchIndexerRequestHandler(c echo.Context) error {
	request := new(AiSearchIndexerRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateOrUpdateSearchIndex(c)
}

func (r *AiSearchIndexerRequest) CreateOrUpdateSearchIndex(c echo.Context) error {
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
	if len(r.SecretKeyName) <= 0 {
		r.SecretKeyName = r.Platform
	}
	r.MaxResults = 100
	switch r.Platform {
	case "reddit":
		resp, err := hera_search.InsertRedditSearchQuery(c.Request().Context(), ou, r.SearchGroupName, r.Query, r.MaxResults)
		if err != nil {
			log.Err(err).Msg("error inserting reddit search query")
			return c.JSON(http.StatusInternalServerError, nil)
		}
		return c.JSON(http.StatusOK, resp)
	case "twitter":
		resp, err := hera_search.InsertTwitterSearchQuery(c.Request().Context(), ou, r.SearchGroupName, r.Query, r.MaxResults)
		if err != nil {
			log.Err(err).Msg("error inserting twitter search query")
			return c.JSON(http.StatusInternalServerError, nil)
		}
		return c.JSON(http.StatusOK, resp)
	case "discord":
		resp, err := hera_search.InsertDiscordSearchQuery(c.Request().Context(), ou, r.SearchGroupName, r.Query, r.MaxResults)
		if err != nil {
			log.Err(err).Msg("error inserting discord search query")
			return c.JSON(http.StatusInternalServerError, nil)
		}
		return c.JSON(http.StatusOK, resp)
	case "telegram":
	}
	return c.JSON(http.StatusOK, nil)
}

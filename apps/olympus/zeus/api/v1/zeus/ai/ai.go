package zeus_v1_ai

import (
	"github.com/labstack/echo/v4"
)

func AiV1Routes(e *echo.Group) *echo.Group {
	e.POST("/search", AiSearchRequestHandler)
	e.POST("/search/analyze", AiSearchAnalyzeRequestHandler)
	e.POST("/search/indexer", AiSearchIndexerRequestHandler)
	e.POST("/search/indexer/actions", SearchIndexerActionsRequestHandler)

	e.GET("/workflows/ai", GetWorkflowsRequestHandler)
	e.POST("/workflows/ai", PostWorkflowsRequestHandler)
	e.POST("/workflows/ai/actions", WorkflowsActionsRequestHandler)
	e.POST("/runs/ai/actions", RunsActionsRequestHandler)

	e.POST("/assistants/ai", CreateOrUpdateAssistantRequestHandler)
	e.POST("/actions/ai", AiActionsHandler)
	e.POST("/tasks/ai", CreateOrUpdateTaskRequestHandler)
	e.POST("/evals/ai", CreateOrUpdateEvalsRequestHandler)

	e.POST("/retrievals/ai", CreateOrUpdateRetrievalRequestHandler)
	// destructive
	e.DELETE("/workflows/ai", WorkflowsDeletionRequestHandler)
	return e
}

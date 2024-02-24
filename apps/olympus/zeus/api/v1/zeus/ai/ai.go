package zeus_v1_ai

import (
	"github.com/labstack/echo/v4"
)

func AiV1Routes(e *echo.Group) *echo.Group {
	// search
	e.POST("/search", AiSearchRequestHandler)
	e.POST("/search/analyze", AiSearchAnalyzeRequestHandler)
	e.POST("/search/indexer", AiSearchIndexerRequestHandler)
	e.POST("/search/indexer/actions", SearchIndexerActionsRequestHandler)

	// schemas
	e.GET("/schemas/ai", AiSchemaHHandler)
	e.POST("/schemas/ai", AiSchemasHandler)

	e.POST("/assistants/ai", CreateOrUpdateAssistantRequestHandler)

	e.GET("/entities/ai", SelectEntitiesRequestHandler)
	e.POST("/entities/ai", CreateOrUpdateEntitiesRequestHandler)

	// tasks
	e.GET("/tasks/ai", GetTasksRequestHandler)
	e.POST("/tasks/ai", CreateOrUpdateTaskRequestHandler)

	e.GET("/evals/ai", GetEvalsRequestHandler)
	e.POST("/evals/ai", CreateOrUpdateEvalsRequestHandler)

	e.POST("/retrievals/ai", CreateOrUpdateRetrievalRequestHandler)

	e.GET("/actions/ai", AiActionsReaderHandler)
	e.POST("/actions/ai", AiActionsHandler)
	e.PUT("/actions/ai", AiActionsApprovalHandler)

	// workflows
	e.GET("/workflows/ai", GetWorkflowsRequestHandler)
	e.POST("/workflows/ai", PostWorkflowsRequestHandler)
	e.POST("/workflows/ai/actions", WorkflowsActionsRequestHandler)
	e.POST("/runs/ai/actions", RunsActionsRequestHandler)

	// runs
	e.GET("/runs/ai", GetRunReportsRequestHandler)

	// destructive
	e.DELETE("/workflows/ai", WorkflowsDeletionRequestHandler)
	return e
}

package zeus_v1_ai

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func AiV1Routes(e *echo.Group) *echo.Group {
	// search
	e.POST("/search", AiSearchRequestHandler)
	e.POST("/search/analyze", AiSearchAnalyzeRequestHandler)
	e.POST("/search/indexer", AiSearchIndexerRequestHandler)
	e.POST("/search/indexer/actions", SearchIndexerActionsRequestHandler)

	// schemas
	e.GET("/schemas/ai", AiSchemaHandler)
	e.POST("/schemas/ai", AiSchemasHandler)

	e.POST("/assistants/ai", CreateOrUpdateAssistantRequestHandler)

	e.POST("/search/entities/ai", SelectEntitiesRequestHandler)
	e.POST("/entity/ai", CreateOrUpdateEntityRequestHandler)
	e.POST("/entities/ai", CreateOrUpdateEntitiesRequestHandler)

	e.GET("/task/ai/:id", GetTaskRequestHandler)
	e.GET("/tasks/ai", GetTasksRequestHandler)
	e.POST("/tasks/ai", CreateOrUpdateTaskRequestHandler)

	e.GET("/eval/ai/:id", GetEvalRequestHandler)
	e.GET("/evals/ai", GetEvalsRequestHandler)
	e.POST("/evals/ai", CreateOrUpdateEvalsRequestHandler)

	e.GET("/retrievals/ai", GetRetrievalsRequestHandler)
	e.GET("/retrieval/ai/:id", GetRetrievalRequestHandler)

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
	e.GET("/run/ai/:id", GetRunActionsRequestHandler)
	e.GET("/runs/ai/ui", GetUIRunReportsRequestHandler)
	e.GET("/runs/ai", GetRunReportsRequestHandler)

	e.POST("/flows", FlowsActionsRequestHandler)
	// for a 10M for 10 MB limit
	e.POST("/flows/exec", FlowsExecActionsRequestHandler, middleware.BodyLimit("100M"))

	// destructive
	e.DELETE("/workflows/ai", WorkflowsDeletionRequestHandler)
	return e
}

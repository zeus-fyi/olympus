package workload_state

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// UpdateWorkloadStateHandler TODO must verify this is auth is scoped to user only
func UpdateWorkloadStateHandler(c echo.Context) error {
	ctx := context.Background()
	request := new(InternalWorkloadStatusUpdate)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	err := request.InsertStatus(ctx)
	if err != nil {
		log.Err(err).Msg("UpdateWorkloadStateHandler")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, nil)
}

package v1_iris

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Proxy(c echo.Context) error {
	// TODO set queue
	return c.String(http.StatusOK, "Healthy")
}

package v1_tyche

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Echo) *echo.Echo {
	e.GET("/health", Health)
	e.POST("/mev/tx", TxProcessingRequestHandler)
	return e
}

func Health(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy")
}

/*

needs an endpoint to kickoff active trade analysis

peers -> get total peer count

decoder
  -> get tx to address
  -> get starting token/end token
  -> get trade method

*/

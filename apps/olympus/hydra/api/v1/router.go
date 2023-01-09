package v1_hypnos

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func Routes(e *echo.Echo) *echo.Echo {
	e.GET("/health", Health)
	e.POST("/api/v1/eth2/sign/:identifier", Eth2SignRequest)
	return e
}

func Eth2SignRequest(c echo.Context) error {
	identifier := c.Param("identifier")
	log.Info().Str("identifier", identifier)

	resp := Eth2SignResponse{Signature: ""}

	// TODO, Send here, using temporal, wait for reply, then send back
	return c.JSON(http.StatusOK, resp)
}

type Eth2SignResponse struct {
	Signature string `json:"signature"`
}

func Health(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy")
}

package hydra_eth2_web3signer

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

const (
	Eth2SignRoute = "/api/v1/eth2/sign/:identifier"
)

type Eth2SignResponse struct {
	Signature string `json:"signature"`
}

func Eth2SignRequest(c echo.Context) error {
	identifier := c.Param("identifier")
	log.Info().Str("identifier", identifier)

	var jsonBody any
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		return c.JSON(http.StatusBadRequest, jsonBody)
	}

	// TODO, Send here, using temporal, wait for reply, then send back
	resp := Eth2SignResponse{Signature: ""}
	return c.JSON(http.StatusOK, resp)
}

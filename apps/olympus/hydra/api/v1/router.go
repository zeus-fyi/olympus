package v1_hypnos

import (
	"net/http"

	"github.com/labstack/echo/v4"
	hydra_eth2_web3signer "github.com/zeus-fyi/olympus/hydra/api/v1/web3signer"
)

func Routes(e *echo.Echo) *echo.Echo {
	e.GET("/health", Health)
	e.POST(hydra_eth2_web3signer.Eth2SignRoute, hydra_eth2_web3signer.HydraEth2SignRequestHandler)
	return e
}

func Health(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy")
}

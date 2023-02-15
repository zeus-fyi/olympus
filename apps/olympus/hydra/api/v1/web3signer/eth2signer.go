package hydra_eth2_web3signer

import (
	"context"
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

type Web3SignerRequest struct {
	Body echo.Map
}

func HydraEth2SignRequestHandler(c echo.Context) error {
	request := new(Web3SignerRequest)
	request.Body = echo.Map{}
	if err := c.Bind(&request.Body); err != nil {
		log.Err(err)
		return err
	}
	return request.Eth2SignRequest(c)
}

func (w *Web3SignerRequest) Eth2SignRequest(c echo.Context) error {
	log.Info().Msg("Eth2SignRequest")
	pubkey := c.Param("identifier")
	ctx := context.Background()
	sr, err := Watermarking(ctx, pubkey, w)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("pubkey", pubkey).Interface("body", w.Body).Msg("Eth2SignRequest")
		return c.JSON(http.StatusInternalServerError, err)
	}
	resp, err := WaitForSignature(ctx, sr)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("pubkey", pubkey).Interface("body", w.Body).Msg("Eth2SignRequest")
		return c.JSON(http.StatusInternalServerError, err)
	}
	if resp.Signature == "" {
		log.Ctx(ctx).Err(err).Interface("pubkey", pubkey).Interface("body", w.Body).Msg("Eth2SignRequest: Signature Field Was Empty")
		return c.JSON(http.StatusRequestTimeout, resp)
	}
	return c.JSON(http.StatusOK, resp)
}

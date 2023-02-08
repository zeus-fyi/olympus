package hydra_eth2_web3signer

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	aegis_inmemdbs "github.com/zeus-fyi/zeus/pkg/aegis/inmemdbs"
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

	err := Watermarking(pubkey, w)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	// TODO, Send here, using temporal? or just send to a queue?
	// maybe just have it wait on a channel?
	// lookup function/endpoint & send & wait for reply, then send back

	// Unique identifier for this request (e.g. UUID), type, pubkey

	// Aggregate requests, then send in batch
	resp := Eth2SignResponse{Signature: ""}
	return c.JSON(http.StatusNotFound, resp)
}

/*
// TODO queue up requests, then send in batch?
var tsCache = cache.New(1*time.Second, 2*time.Second)
*/

// TODO, this is a mock response, for testing delete/replace
func MockResponse(c echo.Context) error {
	respMsgMap := make(map[string]aegis_inmemdbs.EthereumBLSKeySignatureResponse)
	signedEventResponse := aegis_inmemdbs.EthereumBLSKeySignatureResponses{
		Map: respMsgMap,
	}
	return c.JSON(http.StatusNotFound, signedEventResponse)
}

package hydra_eth2_web3signer

import (
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

const (
	Eth2SignRoute = "/api/v1/eth2/sign/:identifier"

	ATTESTATION                           = "ATTESTATION"
	AGGREGATION_SLOT                      = "AGGREGATION_SLOT"
	AGGREGATE_AND_PROOF                   = "AGGREGATE_AND_PROOF"
	BLOCK                                 = "BLOCK"
	BLOCK_V2                              = "BLOCK_V2"
	RANDAO_REVEAL                         = "RANDAO_REVEAL"
	SYNC_COMMITTEE_MESSAGE                = "SYNC_COMMITTEE_MESSAGE"
	SYNC_COMMITTEE_SELECTION_PROOF        = "SYNC_COMMITTEE_SELECTION_PROOF"
	SYNC_COMMITTEE_CONTRIBUTION_AND_PROOF = "SYNC_COMMITTEE_CONTRIBUTION_AND_PROOF"
	VALIDATOR_REGISTRATION                = "VALIDATOR_REGISTRATION"
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

	var sr SignRequest

	signType := w.Body["type"]
	switch signType {
	case ATTESTATION:
		// TODO watermark
		log.Info().Interface("body", w.Body).Msg("ATTESTATION")
	case AGGREGATION_SLOT:
	case AGGREGATE_AND_PROOF:
	case BLOCK:
		// TODO watermark
		log.Info().Interface("body", w.Body).Msg("BLOCK")
	case BLOCK_V2:
		log.Info().Interface("body", w.Body).Msg("BLOCK_V2")
		// TODO watermark
	case RANDAO_REVEAL:
		log.Info().Interface("body", w.Body).Msg("RANDAO_REVEAL")
	case SYNC_COMMITTEE_MESSAGE:
	case SYNC_COMMITTEE_SELECTION_PROOF:
	case SYNC_COMMITTEE_CONTRIBUTION_AND_PROOF:
	case VALIDATOR_REGISTRATION:
		// TODO watermark?
	default:

	}
	signingRoot := w.Body["signingRoot"]

	sr.UUID = uuid.New()
	sr.Type = signType.(string)
	sr.Pubkey = pubkey
	sr.SigningRoot = signingRoot.(string)

	// TODO, Send here, using temporal? or just send to a queue?
	// maybe just have it wait on a channel?

	// lookup function/endpoint & send & wait for reply, then send back

	// ideally can aggregate requests, and send in batch

	// Unique identifier for this request (e.g. UUID), type, pubkey
	resp := Eth2SignResponse{Signature: ""}
	return c.JSON(http.StatusNotFound, resp)
}

type SignRequest struct {
	UUID uuid.UUID `json:"uuid"`

	Pubkey      string `json:"pubkey"`
	Type        string `json:"type"`
	SigningRoot string `json:"signingRoot"`
}

// TODO queue up requests, then send in batch?
var tsCache = cache.New(1*time.Second, 2*time.Second)

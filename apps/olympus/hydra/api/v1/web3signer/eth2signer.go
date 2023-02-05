package hydra_eth2_web3signer

import (
	"encoding/json"
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

func Eth2SignRequest(c echo.Context) error {
	pubkey := c.Param("identifier")
	log.Info().Str("identifier", pubkey)

	var jsonBody map[string]any
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		return c.JSON(http.StatusBadRequest, jsonBody)
	}

	var sr SignRequest

	signType := jsonBody["type"]
	switch signType {
	case ATTESTATION:
		// TODO watermark
	case AGGREGATION_SLOT:
	case AGGREGATE_AND_PROOF:
	case BLOCK:
		// TODO watermark
	case BLOCK_V2:
		// TODO watermark
	case RANDAO_REVEAL:
	case SYNC_COMMITTEE_MESSAGE:
	case SYNC_COMMITTEE_SELECTION_PROOF:
	case SYNC_COMMITTEE_CONTRIBUTION_AND_PROOF:
	case VALIDATOR_REGISTRATION:
		// TODO watermark?
	default:

	}
	signingRoot := jsonBody["signingRoot"]

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
	return c.JSON(http.StatusOK, resp)
}

type SignRequest struct {
	UUID uuid.UUID `json:"uuid"`

	Pubkey      string `json:"pubkey"`
	Type        string `json:"type"`
	SigningRoot string `json:"signingRoot"`
}

// TODO queue up requests, then send in batch?
var tsCache = cache.New(1*time.Second, 2*time.Second)

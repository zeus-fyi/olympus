package hydra_eth2_web3signer

import (
	"encoding/json"
	"github.com/patrickmn/go-cache"
	aegis_inmemdbs "github.com/zeus-fyi/zeus/pkg/aegis/inmemdbs"
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

	signType := jsonBody["type"]
	switch signType {
	case ATTESTATION:
	case AGGREGATION_SLOT:
	case AGGREGATE_AND_PROOF:
	case BLOCK:
	case BLOCK_V2:
	case RANDAO_REVEAL:
	case SYNC_COMMITTEE_MESSAGE:
	case SYNC_COMMITTEE_SELECTION_PROOF:
	case SYNC_COMMITTEE_CONTRIBUTION_AND_PROOF:
	case VALIDATOR_REGISTRATION:

	}
	// TODO, Send here, using temporal, wait for reply, then send back
	// Unique identifier for this request (e.g. UUID), type, pubkey
	resp := Eth2SignResponse{Signature: ""}
	return c.JSON(http.StatusOK, resp)
}

type SignRequest struct {
	aegis_inmemdbs.EthereumBLSKeySignatureRequests
}

// TODO queue up requests, then send in batch
var tsCache = cache.New(1*time.Second, 2*time.Second)

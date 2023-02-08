package hydra_eth2_web3signer

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	bls_serverless_signing "github.com/zeus-fyi/zeus/pkg/aegis/aws/serverless_signing"
	aegis_inmemdbs "github.com/zeus-fyi/zeus/pkg/aegis/inmemdbs"
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
		log.Info().Interface("body", w.Body).Msg("AGGREGATION_SLOT")
	case AGGREGATE_AND_PROOF:
		log.Info().Interface("body", w.Body).Msg("AGGREGATE_AND_PROOF")
	case BLOCK:
		// TODO watermark
		log.Info().Interface("body", w.Body).Msg("BLOCK")
	case BLOCK_V2:
		log.Info().Interface("body", w.Body).Msg("BLOCK_V2")
		// TODO watermark
	case RANDAO_REVEAL:
		log.Info().Interface("body", w.Body).Msg("RANDAO_REVEAL")
	case SYNC_COMMITTEE_MESSAGE:
		log.Info().Interface("body", w.Body).Msg("SYNC_COMMITTEE_MESSAGE")
	case SYNC_COMMITTEE_SELECTION_PROOF:
		log.Info().Interface("body", w.Body).Msg("SYNC_COMMITTEE_SELECTION_PROOF")
	case SYNC_COMMITTEE_CONTRIBUTION_AND_PROOF:
		log.Info().Interface("body", w.Body).Msg("SYNC_COMMITTEE_CONTRIBUTION_AND_PROOF")
	case VALIDATOR_REGISTRATION:
		log.Info().Interface("body", w.Body).Msg("VALIDATOR_REGISTRATION")
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
	sigReqs := bls_serverless_signing.SignatureRequests{
		SecretName:        "",
		SignatureRequests: aegis_inmemdbs.EthereumBLSKeySignatureRequests{Map: make(map[string]aegis_inmemdbs.EthereumBLSKeySignatureRequest)},
	}
	sigReqs.SignatureRequests.Map[sr.Pubkey] = aegis_inmemdbs.EthereumBLSKeySignatureRequest{Message: sr.SigningRoot}
	// Unique identifier for this request (e.g. UUID), type, pubkey

	// Aggregate requests, then send in batch
	resp := Eth2SignResponse{Signature: ""}
	return c.JSON(http.StatusNotFound, resp)
}

type SignRequest struct {
	UUID uuid.UUID `json:"uuid"`

	Pubkey      string `json:"pubkey"`
	Type        string `json:"type"`
	SigningRoot string `json:"signingRoot"`
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

package eth_validator_signature_requests

import (
	"context"
	"time"

	aegis_inmemdbs "github.com/zeus-fyi/zeus/pkg/aegis/inmemdbs"
)

const (
	waitForTxRxTimeout    = 15 * time.Minute
	submitSignedTxTimeout = 5 * time.Minute
)

type ArtemisEthereumValidatorSignatureRequestActivities struct {
}

func NewArtemisEthereumValidatorSignatureRequestActivities() ArtemisEthereumValidatorSignatureRequestActivities {
	return ArtemisEthereumValidatorSignatureRequestActivities{}
}

type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

type Eth2SignResponse struct {
	Signature string `json:"signature"`
}

func (d *ArtemisEthereumValidatorSignatureRequestActivities) GetActivities() ActivitiesSlice {
	return []interface{}{d.RequestValidatorSignatures}
}

func (d *ArtemisEthereumValidatorSignatureRequestActivities) RequestValidatorSignatures(ctx context.Context, sigRequests aegis_inmemdbs.EthereumBLSKeySignatureRequests) (aegis_inmemdbs.EthereumBLSKeySignatureResponses, error) {
	// TODO serverless request here
	// TODO, group pubkeys by serverless function then send requests

	m := make(map[string]aegis_inmemdbs.EthereumBLSKeySignatureRequests)
	// TODO: group by service url
	for pubkey, signReq := range sigRequests.Map {
		svcURL := MockGetServiceURL(pubkey)
		if _, ok := m[svcURL]; !ok {
			m[svcURL] = aegis_inmemdbs.EthereumBLSKeySignatureRequests{}
		}
		m[svcURL].Map[pubkey] = signReq
	}
	return aegis_inmemdbs.EthereumBLSKeySignatureResponses{}, nil
}

func MockGetServiceURL(pubkey string) string {
	// TODO lookup in cache
	return "http://localhost:8080"
}

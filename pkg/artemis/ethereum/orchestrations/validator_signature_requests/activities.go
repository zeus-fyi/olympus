package eth_validator_signature_requests

import (
	"context"
	aegis_inmemdbs "github.com/zeus-fyi/zeus/pkg/aegis/inmemdbs"
	"time"
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
	return []interface{}{d.RequestValidatorSignature}
}

func (d *ArtemisEthereumValidatorSignatureRequestActivities) RequestValidatorSignature(ctx context.Context, payload any) (aegis_inmemdbs.EthereumBLSKeySignatureResponses, error) {
	// TODO serverless request here
	return aegis_inmemdbs.EthereumBLSKeySignatureResponses{}, nil
}

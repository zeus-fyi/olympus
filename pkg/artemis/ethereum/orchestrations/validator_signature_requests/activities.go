package eth_validator_signature_requests

import (
	"context"
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

func (d *ArtemisEthereumValidatorSignatureRequestActivities) RequestValidatorSignature(ctx context.Context, payload any) (Eth2SignResponse, error) {

	// TODO serverless request here
	resp := Eth2SignResponse{}
	return resp, nil
}

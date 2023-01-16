package eth_validators_service_requests

import (
	"context"
	"time"
)

const (
	waitForTxRxTimeout    = 15 * time.Minute
	submitSignedTxTimeout = 5 * time.Minute
)

type ArtemisEthereumValidatorsServiceRequestActivities struct {
}

func NewArtemisEthereumValidatorSignatureRequestActivities() ArtemisEthereumValidatorsServiceRequestActivities {
	return ArtemisEthereumValidatorsServiceRequestActivities{}
}

type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (d *ArtemisEthereumValidatorsServiceRequestActivities) GetActivities() ActivitiesSlice {
	return []interface{}{d.SendValidatorsToCloudCtxNs}
}

func (d *ArtemisEthereumValidatorsServiceRequestActivities) SendValidatorsToCloudCtxNs(ctx context.Context, payload any) error {
	return nil
}

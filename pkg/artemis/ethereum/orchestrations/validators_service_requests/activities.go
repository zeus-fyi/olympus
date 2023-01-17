package eth_validators_service_requests

import (
	"context"
	"time"

	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
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
	return []interface{}{d.AssignValidatorsToCloudCtxNs}
}

type ArtemisEthereumValidatorsServiceRequestPayload struct {
	CloudCtxNsID, ProtocolNetworkID int
}

func (d *ArtemisEthereumValidatorsServiceRequestActivities) AssignValidatorsToCloudCtxNs(ctx context.Context, params ArtemisEthereumValidatorsServiceRequestPayload) error {
	err := artemis_validator_service_groups_models.SelectInsertUnplacedValidatorsIntoCloudCtxNs(ctx, params.ProtocolNetworkID, params.CloudCtxNsID)
	if err != nil {
		return err
	}

	return nil
}

func (d *ArtemisEthereumValidatorsServiceRequestActivities) SendValidatorsToCloudCtxNs(ctx context.Context, params ArtemisEthereumValidatorsServiceRequestPayload) error {
	// query all validators that should be in this cluster, then patch the validator yaml file

	return nil
}

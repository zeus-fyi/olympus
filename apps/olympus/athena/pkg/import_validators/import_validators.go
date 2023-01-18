package athena_import_validators

import (
	"context"

	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
)

// TODO finish
func SelectAssignedValidatorsToCluster(ctx context.Context, vsg artemis_validator_service_groups_models.ValidatorServiceCloudCtxNsProtocol) error {
	_, err := artemis_validator_service_groups_models.SelectValidatorsAssignedToCloudCtxNs(ctx, vsg)
	// []ethereum_web3signer_actions.LighthouseWeb3SignerRequest
	return err
}

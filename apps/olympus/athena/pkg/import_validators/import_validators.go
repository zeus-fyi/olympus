package athena_import_validators

import (
	"context"

	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	ethereum_web3signer_actions "github.com/zeus-fyi/zeus/cookbooks/ethereum/web3signers/actions"
)

// TODO finish
func SelectAndAssignLighthouseValidatorsToCluster(ctx context.Context, vsg artemis_validator_service_groups_models.ValidatorServiceCloudCtxNsProtocol) ([]ethereum_web3signer_actions.LighthouseWeb3SignerRequest, error) {
	lh, err := artemis_validator_service_groups_models.SelectValidatorsAssignedToCloudCtxNs(ctx, vsg)
	// []ethereum_web3signer_actions.LighthouseWeb3SignerRequest
	return lh, err
}

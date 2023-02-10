package artemis_validator_signature_service_routing

import (
	"context"
	"github.com/rs/zerolog/log"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	bls_serverless_signing "github.com/zeus-fyi/zeus/pkg/aegis/aws/serverless_signing"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

type ServiceRoutes struct {
	To []ServiceRoute
}

type ServiceRoute struct {
	AuthToken string // TODO

	ServiceURL string
	bls_serverless_signing.SignatureRequests
}

func GetServiceURLs(ctx context.Context, cctx zeus_common_types.CloudCtxNs) (artemis_validator_service_groups_models.ValidatorsSignatureServiceRoutes, error) {
	vsi := artemis_validator_service_groups_models.ValidatorServiceCloudCtxNsProtocol{}
	vsRoutes, err := artemis_validator_service_groups_models.SelectValidatorsServiceRoutesAssignedToCloudCtxNs(ctx, vsi, cctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("GetServiceURL")
		return vsRoutes, err
	}
	// TODO get auth info + chain messages to sign

	return vsRoutes, err
}

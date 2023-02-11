package artemis_validator_signature_service_routing

import (
	"context"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	bls_serverless_signing "github.com/zeus-fyi/zeus/pkg/aegis/aws/serverless_signing"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

// TODO get auth info + chain messages to sign

// map pubkey -> service url -> batchSignReq to this url

/*
 v.GroupName, ou.OrgID, v.ProtocolNetworkID)
*/

type ServiceRoutes struct {
	To []ServiceRoute
}

type ServiceRoute struct {
	AuthToken string // TODO

	ServiceURL string
	bls_serverless_signing.SignatureRequests
}

var ServiceAuthRouteCache = cache.New(12*time.Hour, 24*time.Hour)

func InitAsyncServiceAuthRoutePolling(ctx context.Context, cctx zeus_common_types.CloudCtxNs) {
	for {
		sr, err := GetServiceURLs(ctx, cctx)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("InitAsyncServiceAuthRoutePolling")
		}
		for pubkey, svc := range sr.Map {
			// TODO group by service url & pubkey
			ServiceAuthRouteCache.Set(pubkey, svc, cache.DefaultExpiration)
		}
		time.Sleep(60 * time.Second)
	}
}

func GetServiceURLs(ctx context.Context, cctx zeus_common_types.CloudCtxNs) (artemis_validator_service_groups_models.ValidatorsSignatureServiceRoutes, error) {
	vsi := artemis_validator_service_groups_models.ValidatorServiceCloudCtxNsProtocol{}
	vsRoutes, err := artemis_validator_service_groups_models.SelectValidatorsServiceRoutesAssignedToCloudCtxNs(ctx, vsi, cctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("GetServiceURL")
		return vsRoutes, err
	}
	return vsRoutes, err
}

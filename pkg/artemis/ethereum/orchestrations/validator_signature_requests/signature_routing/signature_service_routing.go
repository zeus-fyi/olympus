package artemis_validator_signature_service_routing

import (
	"context"
	"fmt"
	artemis_hydra_orchestrations_aws_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/aws_auth"
	aws_secrets "github.com/zeus-fyi/zeus/pkg/aegis/aws"
	aegis_inmemdbs "github.com/zeus-fyi/zeus/pkg/aegis/inmemdbs"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	"time"

	"github.com/patrickmn/go-cache"
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

var ServiceAuthRouteCache = cache.New(12*time.Hour, 24*time.Hour)

func InitAsyncServiceAuthRoutePolling(ctx context.Context, cctx zeus_common_types.CloudCtxNs) {
	for {
		sr, err := GetServiceURLs(ctx, cctx)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("InitAsyncServiceAuthRoutePolling")
		}
		for pubkey, svc := range sr.Map {
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

func formatSecret(groupName string, orgID, protocolNetworkID int) string {
	return fmt.Sprintf("%s-%d-%s", groupName, orgID, hestia_req_types.ProtocolNetworkIDToString(protocolNetworkID))
}

func GetServiceRoutesAuths(ctx context.Context, vsRoute artemis_validator_service_groups_models.ValidatorsSignatureServiceRoutes) error {
	si := aws_secrets.SecretInfo{
		Region: "us-west-1",
		Name:   "",
	}
	for _, v := range vsRoute.Map {
		si.Name = formatSecret(v.GroupName, v.OrgID, v.ProtocolNetworkID)
		s, err := artemis_hydra_orchestrations_aws_auth.GetServiceRoutesAuths(ctx, si)
		if err != nil {
			return err
		}
		fmt.Sprint(s)
	}

	return nil
}

func GroupSigRequestsByServiceURL(ctx context.Context, sr aegis_inmemdbs.EthereumBLSKeySignatureRequests) map[string]aegis_inmemdbs.EthereumBLSKeySignatureRequests {
	m := make(map[string]aegis_inmemdbs.EthereumBLSKeySignatureRequests)
	for pubkey, signReq := range sr.Map {
		svcURL := MockGetServiceURL(pubkey)
		if _, ok := m[svcURL]; !ok {
			m[svcURL] = aegis_inmemdbs.EthereumBLSKeySignatureRequests{Map: make(map[string]aegis_inmemdbs.EthereumBLSKeySignatureRequest)}
		}
		m[svcURL].Map[pubkey] = signReq
	}
	return m
}

func MockGetServiceURL(pubkey string) string {
	// TODO lookup in cache
	return "http://localhost:8080"
}

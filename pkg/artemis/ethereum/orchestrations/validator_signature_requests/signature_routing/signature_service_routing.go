package artemis_validator_signature_service_routing

import (
	"context"
	"errors"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
	bls_serverless_signing "github.com/zeus-fyi/zeus/pkg/aegis/aws/serverless_signing"
	aegis_inmemdbs "github.com/zeus-fyi/zeus/pkg/aegis/inmemdbs"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

var RouteMapInMemFS memfs.MemFS

const (
	serviceKeyGroupsDir = "./service_key_groups"
)

func InitRouteMapInMemFS(ctx context.Context) error {
	RouteMapInMemFS = memfs.NewMemFs()
	return nil
}

func InitAsyncServiceAuthRoutePolling(ctx context.Context, vsi artemis_validator_service_groups_models.ValidatorServiceCloudCtxNsProtocol, cctx zeus_common_types.CloudCtxNs) {
	log.Ctx(ctx).Info().Interface("cctx", cctx).Msg("InitAsyncServiceAuthRoutePolling")
	for {
		err := GetServiceAuthAndURLs(ctx, vsi, cctx)
		if err != nil {
			log.Ctx(ctx).Err(err).Msg("GetServiceAuthAndURLs")
		}
		time.Sleep(60 * time.Second)
	}
}

func GetServiceMetadata(ctx context.Context, vsi artemis_validator_service_groups_models.ValidatorServiceCloudCtxNsProtocol, cctx zeus_common_types.CloudCtxNs) (artemis_validator_service_groups_models.ValidatorsSignatureServiceRoutes, error) {
	log.Ctx(ctx).Info().Interface("cctx", cctx).Msg("GetServiceMetadata")
	vsRoutes, err := artemis_validator_service_groups_models.SelectValidatorsServiceRoutesAssignedToCloudCtxNs(ctx, vsi, cctx)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("GetServiceURL")
		return vsRoutes, err
	}
	return vsRoutes, err
}

func GetAllServiceMetadata(ctx context.Context) (artemis_validator_service_groups_models.ValidatorsSignatureServiceRoutes, error) {
	vsRoutes, err := artemis_validator_service_groups_models.SelectValidatorsServiceRoutes(ctx)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("GetServiceURL")
		return vsRoutes, err
	}
	return vsRoutes, err
}

func InitAsyncServiceAuthRoutePollingHeartbeatAll(ctx context.Context) {
	i := 0
	for {
		svcGroups, err := GetAllServiceAuthAndURLs(ctx)
		if err != nil {
			log.Ctx(ctx).Err(err).Msg("GetAllServiceAuthAndURLs")
		}
		SendHeartbeat(ctx, svcGroups)
		time.Sleep(30 * time.Second)
		if i > 30 {
			svcGroups, err = GetAllServiceAuthAndURLs(ctx)
			if err != nil {
				log.Ctx(ctx).Err(err).Msg("GetAllServiceAuthAndURLs")
			}
			i = 0
		}
	}
}

var ConcurrentHeartbeatSize = 3

func SendHeartbeat(ctx context.Context, svcGroups artemis_validator_service_groups_models.ValidatorsSignatureServiceRoutes) {
	log.Ctx(ctx).Info().Msg("SendHeartbeat")
	for groupName, _ := range svcGroups.GroupToServiceMap {
		log.Ctx(ctx).Info().Interface("groupName", groupName).Msg("SendHeartbeat")
		for i := 0; i < ConcurrentHeartbeatSize; i++ {
			go func(groupName string) {
				log.Ctx(ctx).Info().Interface("groupName", groupName).Msg("sending heartbeat message")
				auth, err := GetGroupAuthFromInMemFS(ctx, groupName)
				if err != nil {
					log.Ctx(ctx).Err(err).Msg("Failed to get group auth")
					return
				}
				signReqs := aegis_inmemdbs.EthereumBLSKeySignatureRequests{Map: make(map[string]aegis_inmemdbs.EthereumBLSKeySignatureRequest)}
				hexMessage, err := aegis_inmemdbs.RandomHex(10)
				if err != nil {
					log.Ctx(ctx).Err(err).Msg("failed to create random hex message")
					return
				}
				signReqs.Map["0x0000000"] = aegis_inmemdbs.EthereumBLSKeySignatureRequest{Message: hexMessage}
				sr := bls_serverless_signing.SignatureRequests{
					SecretName:        auth.SecretName,
					SignatureRequests: aegis_inmemdbs.EthereumBLSKeySignatureRequests{Map: signReqs.Map},
				}
				sigResponses := aegis_inmemdbs.EthereumBLSKeySignatureResponses{Map: make(map[string]aegis_inmemdbs.EthereumBLSKeySignatureResponse)}
				cfg := aegis_aws_auth.AuthAWS{
					Region:    "us-west-1",
					AccessKey: auth.AccessKey,
					SecretKey: auth.SecretKey,
				}
				r := Resty{}
				r.Client = resty.New()
				r.SetBaseURL(auth.AuthLamdbaAWS.ServiceURL)
				r.SetTimeout(5 * time.Second)
				r.SetRetryCount(2)
				r.SetRetryWaitTime(20 * time.Millisecond)
				reqAuth, err := cfg.CreateV4AuthPOSTReq(ctx, "lambda", auth.AuthLamdbaAWS.ServiceURL, sr)
				if err != nil {
					log.Ctx(ctx).Error().Err(err).Msg("failed to get service routes auths for lambda iam auth")
					return
				}
				log.Info().Interface("groupName", groupName).Msg("sending heartbeat")
				resp, err := r.R().
					SetHeaderMultiValues(reqAuth.Header).
					SetResult(&sigResponses).
					SetBody(sr).Post("/")
				if err != nil {
					log.Ctx(ctx).Err(err).Interface("groupName", groupName).Msg("failed to get response")
					return
				}
				if resp.StatusCode() != 200 {
					err = errors.New("non-200 status code")
					log.Ctx(ctx).Err(err).Interface("groupName", groupName).Msg("failed to get 200 status code")
					return
				} else {
					log.Ctx(ctx).Info().Interface("groupName", groupName).Msg("heartbeat OK")
				}
			}(groupName)
		}
	}
}

func GetAllServiceAuthAndURLs(ctx context.Context) (artemis_validator_service_groups_models.ValidatorsSignatureServiceRoutes, error) {
	log.Ctx(ctx).Info().Msg("GetServiceAuthAndURLs")
	sr, err := GetAllServiceMetadata(ctx)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("GetServiceAuthAndURLs: GetServiceMetadata")
		return sr, err
	}
	err = FetchAndSetServiceGroupsAuths(ctx, sr)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("GetServiceRoutesAuths: FetchAndSetServiceGroupsAuths")
		return sr, err
	}
	err = SetPubkeyToGroupInMemFS(ctx, sr)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("GetServiceAuthAndURLs: SetPubkeyToGroupInMemFS")
		return sr, err
	}
	return sr, err
}

func GetServiceAuthAndURLs(ctx context.Context, vsi artemis_validator_service_groups_models.ValidatorServiceCloudCtxNsProtocol, cctx zeus_common_types.CloudCtxNs) error {
	log.Ctx(ctx).Info().Interface("cctx", cctx).Msg("GetServiceAuthAndURLs")
	sr, err := GetServiceMetadata(ctx, vsi, cctx)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("GetServiceAuthAndURLs: GetServiceMetadata")
		return err
	}
	err = FetchAndSetServiceGroupsAuths(ctx, sr)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("GetServiceRoutesAuths: FetchAndSetServiceGroupsAuths")
		return err
	}
	err = SetPubkeyToGroupInMemFS(ctx, sr)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("GetServiceAuthAndURLs: SetPubkeyToGroupInMemFS")
		return err
	}
	return err
}

func SetPubkeyToGroupInMemFS(ctx context.Context, a artemis_validator_service_groups_models.ValidatorsSignatureServiceRoutes) error {
	for pubkey, gn := range a.PubkeyToGroupName {
		err := SetPubkeyToGroupService(ctx, pubkey, gn)
		if err != nil {
			log.Ctx(ctx).Err(err).Msg("SetPubkeyToGroupInMemFS")
			return err
		}
	}
	return nil
}

type Resty struct {
	*resty.Client
}

func SetPubkeyToGroupService(ctx context.Context, pubkey, groupName string) error {
	svcAuthPath := filepaths.Path{
		DirIn: serviceKeyGroupsDir,
		FnIn:  pubkey,
	}
	err := RouteMapInMemFS.MakeFileIn(&svcAuthPath, []byte(groupName))
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("SetPubkeyToGroupService")
		return err
	}
	return nil
}

func GroupSigRequestsByGroupName(ctx context.Context, sr aegis_inmemdbs.EthereumBLSKeySignatureRequests) map[string]aegis_inmemdbs.EthereumBLSKeySignatureRequests {
	m := make(map[string]aegis_inmemdbs.EthereumBLSKeySignatureRequests)
	for pubkey, signReq := range sr.Map {
		groupName, err := GetServiceGroupFromInMemFS(ctx, pubkey)
		if err != nil {
			log.Ctx(ctx).Err(err).Msg("GroupSigRequestsByServiceURL")
			continue
		}
		if _, ok := m[groupName]; !ok {
			m[groupName] = aegis_inmemdbs.EthereumBLSKeySignatureRequests{Map: make(map[string]aegis_inmemdbs.EthereumBLSKeySignatureRequest)}
		}
		m[groupName].Map[pubkey] = signReq
	}
	return m
}

func GetServiceGroupFromInMemFS(ctx context.Context, pubkey string) (string, error) {
	svcPath := filepaths.Path{
		DirIn: serviceKeyGroupsDir,
		FnIn:  pubkey,
	}
	groupName, err := RouteMapInMemFS.ReadFile(svcPath.FileInPath())
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("GetServiceURLFromInMemFS")
		return "", err
	}
	return string(groupName), nil
}

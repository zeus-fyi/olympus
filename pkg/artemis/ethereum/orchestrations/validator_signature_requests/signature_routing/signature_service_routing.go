package artemis_validator_signature_service_routing

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
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

func InitAsyncServiceAuthRoutePolling(ctx context.Context, cctx zeus_common_types.CloudCtxNs) {
	for {
		err := GetServiceAuthAndURLs(ctx, cctx)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("GetServiceAuthAndURLs")
		}
		time.Sleep(60 * time.Second)
	}
}

func GetServiceMetadata(ctx context.Context, cctx zeus_common_types.CloudCtxNs) (artemis_validator_service_groups_models.ValidatorsSignatureServiceRoutes, error) {
	vsi := artemis_validator_service_groups_models.ValidatorServiceCloudCtxNsProtocol{}
	vsRoutes, err := artemis_validator_service_groups_models.SelectValidatorsServiceRoutesAssignedToCloudCtxNs(ctx, vsi, cctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("GetServiceURL")
		return vsRoutes, err
	}
	return vsRoutes, err
}

func GetServiceAuthAndURLs(ctx context.Context, cctx zeus_common_types.CloudCtxNs) error {
	sr, err := GetServiceMetadata(ctx, cctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("GetServiceAuthAndURLs: GetServiceMetadata")
		return err
	}
	err = FetchAndSetServiceGroupsAuths(ctx, sr)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("GetServiceRoutesAuths: FetchAndSetServiceGroupsAuths")
		return err
	}
	err = SetPubkeyToGroupInMemFS(ctx, sr)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("GetServiceAuthAndURLs: SetPubkeyToGroupInMemFS")
		return err
	}
	return err
}

func SetPubkeyToGroupInMemFS(ctx context.Context, a artemis_validator_service_groups_models.ValidatorsSignatureServiceRoutes) error {
	for gn, pubkey := range a.PubkeyToGroupName {
		err := SetPubkeyToGroupService(ctx, pubkey, gn)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("SetPubkeyToGroupInMemFS")
			return err
		}
	}
	return nil
}

func SetPubkeyToGroupService(ctx context.Context, pubkey, groupName string) error {
	svcAuthPath := filepaths.Path{
		DirIn: serviceKeyGroupsDir,
		FnIn:  pubkey,
	}
	err := RouteMapInMemFS.MakeFileIn(&svcAuthPath, []byte(groupName))
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("SetPubkeyToGroupService")
		return err
	}
	return nil
}

func GroupSigRequestsByGroupName(ctx context.Context, sr aegis_inmemdbs.EthereumBLSKeySignatureRequests) map[string]aegis_inmemdbs.EthereumBLSKeySignatureRequests {
	m := make(map[string]aegis_inmemdbs.EthereumBLSKeySignatureRequests)
	for pubkey, signReq := range sr.Map {
		groupName, err := GetServiceGroupFromInMemFS(ctx, pubkey)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("GroupSigRequestsByServiceURL")
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
		log.Ctx(ctx).Error().Err(err).Msg("GetServiceURLFromInMemFS")
		return "", err
	}
	return string(groupName), nil
}

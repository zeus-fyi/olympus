package artemis_validator_signature_service_routing

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"

	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	artemis_hydra_orchestrations_aws_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/aws_auth"
	aws_secrets "github.com/zeus-fyi/zeus/pkg/aegis/aws"
)

const (
	serviceGroupsAuthsDir = "./service_group_auths"
)

func FetchAndSetServiceGroupsAuths(ctx context.Context, vsRoute artemis_validator_service_groups_models.ValidatorsSignatureServiceRoutes) error {
	si := aws_secrets.SecretInfo{
		Region: "us-west-1",
		Name:   "",
	}
	for groupName, v := range vsRoute.GroupToServiceMap {
		svcAuthPath := filepaths.Path{
			DirIn: serviceGroupsAuthsDir,
			FnIn:  groupName,
		}
		_, err := RouteMapInMemFS.Stat(svcAuthPath.FileInPath())
		if err == nil {
			continue
		}
		si.Name = FormatSecretNameAWS(v.GroupName, v.OrgID, v.ProtocolNetworkID)
		s, err := artemis_hydra_orchestrations_aws_auth.GetServiceRoutesAuths(ctx, si)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("GetServiceRoutesAuths")
			return err
		}
		err = SetGroupAuthInMemFS(ctx, s.GroupName, s.ServiceAuth)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("GetServiceRoutesAuths")
			return err
		}
	}
	return nil
}

func SetGroupAuthInMemFS(ctx context.Context, groupName string, serviceAuth hestia_req_types.ServiceAuthConfig) error {
	svcAuthPath := filepaths.Path{
		DirIn: serviceGroupsAuthsDir,
		FnIn:  groupName,
	}
	b, err := json.Marshal(serviceAuth)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("GetServiceRoutesAuths")
		return err
	}
	err = RouteMapInMemFS.MakeFileIn(&svcAuthPath, b)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("GetServiceRoutesAuths")
		return err
	}
	return nil
}

func GetGroupAuthFromInMemFS(ctx context.Context, groupName string) (hestia_req_types.ServiceAuthConfig, error) {
	svcAuthPath := filepaths.Path{
		DirIn: serviceGroupsAuthsDir,
		FnIn:  groupName,
	}
	authCfg := hestia_req_types.ServiceAuthConfig{}
	b, err := RouteMapInMemFS.ReadFile(svcAuthPath.FileInPath())
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("GetServiceRoutesAuths")
		return authCfg, err
	}
	err = json.Unmarshal(b, &authCfg)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("GetServiceRoutesAuths")
		return authCfg, err
	}
	return authCfg, nil
}

func FormatSecretNameAWS(groupName string, orgID, protocolNetworkID int) string {
	return fmt.Sprintf("%s-%d-%s", groupName, orgID, hestia_req_types.ProtocolNetworkIDToString(protocolNetworkID))
}

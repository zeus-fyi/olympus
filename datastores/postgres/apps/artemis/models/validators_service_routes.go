package artemis_validator_service_groups_models

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

// ValidatorsSignatureServiceRoutes uses the validator pubkey as the map key
type ValidatorsSignatureServiceRoutes struct {
	PubkeyToGroupName  map[string]string
	GroupToPubKeySlice map[string][]string
	GroupToServiceMap  map[string]ValidatorsSignatureServiceRoute
}

type ValidatorsSignatureServiceRoute struct {
	GroupName         string `json:"groupName"`
	ServiceURL        string `json:"serviceURL"`
	OrgID             int    `json:"orgID"`
	ProtocolNetworkID int    `json:"protocolNetworkID"`
}

// SelectValidatorsServiceRoutesAssignedToCloudCtxNs is used by hydra
func SelectValidatorsServiceRoutesAssignedToCloudCtxNs(ctx context.Context, validatorServiceInfo ValidatorServiceCloudCtxNsProtocol, cloudCtxNs zeus_common_types.CloudCtxNs) (ValidatorsSignatureServiceRoutes, error) {
	q := sql_query_templates.QueryParams{}
	serviceRoutes := ValidatorsSignatureServiceRoutes{}
	m := make(map[string]string)
	gtkm := make(map[string][]string)
	gts := make(map[string]ValidatorsSignatureServiceRoute)

	q.RawQuery = `	
				  SELECT vsg.pubkey, vsg.group_name, vsg.service_url, vsg.org_id, vsg.protocol_network_id
				  FROM validators_service_org_groups_cloud_ctx_ns vctx
				  INNER JOIN topologies_org_cloud_ctx_ns topctx ON topctx.cloud_ctx_ns_id = vctx.cloud_ctx_ns_id
				  INNER JOIN validators_service_org_groups vsg ON vsg.pubkey = vctx.pubkey
				  WHERE vsg.protocol_network_id=$1 AND vsg.enabled=true AND topctx.cloud_provider=$2 AND topctx.context=$3 AND topctx.region=$4 AND topctx.namespace=$5
				  `
	log.Debug().Interface("SelectValidatorsServiceRoutesAssignedToCloudCtxNs", q.LogHeader(ModelName))
	rows, err := apps.Pg.Query(ctx, q.RawQuery, validatorServiceInfo.ProtocolNetworkID, cloudCtxNs.CloudProvider, cloudCtxNs.Context, cloudCtxNs.Region, cloudCtxNs.Namespace)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(ModelName)); returnErr != nil {
		return serviceRoutes, err
	}
	defer rows.Close()
	for rows.Next() {
		var pubkey string
		vsr := ValidatorsSignatureServiceRoute{}
		rowErr := rows.Scan(
			&pubkey, &vsr.GroupName, &vsr.ServiceURL, &vsr.OrgID, &vsr.ProtocolNetworkID,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(ModelName))
			return serviceRoutes, rowErr
		}
		// pubkey to group name
		groupName := FormatSecretNameAWS(vsr.GroupName, vsr.OrgID, vsr.ProtocolNetworkID)
		m[pubkey] = groupName

		// group to pubkey slice
		tmp := gtkm[groupName]
		gtkm[groupName] = append(tmp, pubkey)

		// group to service map
		if _, ok := gts[groupName]; !ok {
			gts[groupName] = vsr
		}
	}
	serviceRoutes.PubkeyToGroupName = m
	serviceRoutes.GroupToPubKeySlice = gtkm
	serviceRoutes.GroupToServiceMap = gts
	return serviceRoutes, misc.ReturnIfErr(err, q.LogHeader(ModelName))
}

// SelectValidatorsServiceRoutes is used to heartbeat
func SelectValidatorsServiceRoutes(ctx context.Context) (ValidatorsSignatureServiceRoutes, error) {
	log.Ctx(ctx).Info().Msg("SelectValidatorsServiceRoutes")
	q := sql_query_templates.QueryParams{}
	serviceRoutes := ValidatorsSignatureServiceRoutes{}
	m := make(map[string]string)
	gtkm := make(map[string][]string)
	gts := make(map[string]ValidatorsSignatureServiceRoute)

	q.RawQuery = `	
				  SELECT vsg.pubkey, vsg.group_name, vsg.service_url, vsg.org_id, vsg.protocol_network_id
				  FROM validators_service_org_groups_cloud_ctx_ns vctx
				  INNER JOIN topologies_org_cloud_ctx_ns topctx ON topctx.cloud_ctx_ns_id = vctx.cloud_ctx_ns_id
				  INNER JOIN validators_service_org_groups vsg ON vsg.pubkey = vctx.pubkey
				  WHERE vsg.enabled=true
				  `
	log.Debug().Interface("SelectValidatorsServiceRoutes", q.LogHeader(ModelName))
	rows, err := apps.Pg.Query(ctx, q.RawQuery)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(ModelName)); returnErr != nil {
		return serviceRoutes, err
	}
	defer rows.Close()
	for rows.Next() {
		var pubkey string
		vsr := ValidatorsSignatureServiceRoute{}
		rowErr := rows.Scan(
			&pubkey, &vsr.GroupName, &vsr.ServiceURL, &vsr.OrgID, &vsr.ProtocolNetworkID,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(ModelName))
			return serviceRoutes, rowErr
		}
		// pubkey to group name
		groupName := FormatSecretNameAWS(vsr.GroupName, vsr.OrgID, vsr.ProtocolNetworkID)
		m[pubkey] = groupName

		// group to pubkey slice
		tmp := gtkm[groupName]
		gtkm[groupName] = append(tmp, pubkey)

		// group to service map
		if _, ok := gts[groupName]; !ok {
			gts[groupName] = vsr
		}
	}
	serviceRoutes.PubkeyToGroupName = m
	serviceRoutes.GroupToPubKeySlice = gtkm
	serviceRoutes.GroupToServiceMap = gts
	return serviceRoutes, misc.ReturnIfErr(err, q.LogHeader(ModelName))
}

func FormatSecretNameAWS(groupName string, orgID, protocolNetworkID int) string {
	return fmt.Sprintf("%s-%d-%s", groupName, orgID, hestia_req_types.ProtocolNetworkIDToString(protocolNetworkID))
}

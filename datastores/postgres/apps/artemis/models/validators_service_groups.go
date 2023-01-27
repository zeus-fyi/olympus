package artemis_validator_service_groups_models

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	ethereum_web3signer_actions "github.com/zeus-fyi/zeus/cookbooks/ethereum/web3signers/actions"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

type ValidatorServiceCloudCtxNsProtocol struct {
	ProtocolNetworkID int `json:"protocolNetworkID"`
	OrgID             int `json:"orgID"`
}

const ModelName = "ArtemisValidatorsServices"

func SelectUnplacedValidators(ctx context.Context, validatorServiceInfo ValidatorServiceCloudCtxNsProtocol) (OrgValidatorServices, error) {
	vos := OrgValidatorServices{}
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
				  SELECT vsg.pubkey, vsg.service_url, vsg.org_id
				  FROM validators_service_org_groups vsg
				  WHERE vsg.enabled=true AND vsg.protocol_network_id=$1 AND vsg.org_id=$2 AND NOT EXISTS (SELECT pubkey FROM validators_service_org_groups_cloud_ctx_ns) 
				  `
	log.Debug().Interface("SelectUnplacedValidators", q.LogHeader(ModelName))
	rows, err := apps.Pg.Query(ctx, q.RawQuery, validatorServiceInfo.ProtocolNetworkID, validatorServiceInfo.OrgID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(ModelName)); returnErr != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		vr := OrgValidatorService{}
		rowErr := rows.Scan(
			&vr.Pubkey, &vr.ServiceURL, &vr.OrgID,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(ModelName))
			return nil, rowErr
		}
		vos = append(vos, vr)
	}
	return vos, misc.ReturnIfErr(err, q.LogHeader(ModelName))
}

type OrgValidatorServices []OrgValidatorService

type OrgValidatorService struct {
	GroupName         string `json:"groupName"`
	Pubkey            string `json:"pubkey"`
	ProtocolNetworkID int    `json:"protocolNetworkID"`
	ServiceURL        string `json:"serviceURL"`
	OrgID             int    `json:"orgID"`
	Enabled           bool   `json:"enabled"`
}

func InsertVerifiedValidatorsToService(ctx context.Context, validatorServiceInfo OrgValidatorService, pubkeys hestia_req_types.ValidatorServiceOrgGroupSlice) error {
	var rows [][]interface{}
	for i, keyPair := range pubkeys {
		rows[i] = []interface{}{
			validatorServiceInfo.GroupName,
			fmt.Sprintf("%d", validatorServiceInfo.OrgID),
			keyPair.Pubkey,
			fmt.Sprintf("%d", validatorServiceInfo.ProtocolNetworkID),
			keyPair.FeeRecipient,
			validatorServiceInfo.Enabled,
			validatorServiceInfo.ServiceURL,
		}
	}
	columns := []string{"group_name", "org_id", "pubkey", "protocol_network_id", "fee_recipient", "enabled", "service_url"}
	// Use the `pgx.CopyFrom` method to insert the data into the table
	_, err := apps.Pg.Pgpool.CopyFrom(context.Background(), pgx.Identifier{"validators_service_org_groups"}, columns, pgx.CopyFromRows(rows))
	if err != nil {
		log.Ctx(ctx).Err(err)
		return err
	}
	return err
}

// TODO needs to also use capacity and client number assignments

func SelectInsertUnplacedValidatorsIntoCloudCtxNs(ctx context.Context, validatorServiceInfo ValidatorServiceCloudCtxNsProtocol, cloudCtxNs zeus_common_types.CloudCtxNs) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `WITH cte_unplaced_validators AS (		
					  SELECT pubkey, fee_recipient, 
					  FROM validators_service_org_groups
					  WHERE NOT EXISTS (SELECT pubkey FROM validators_service_org_groups_cloud_ctx_ns) AND enabled=true AND protocol_network_id=$1
				  ) INSERT INTO validators_service_org_groups_cloud_ctx_ns(pubkey, cloud_ctx_ns_id)
					SELECT pubkey, (SELECT cloud_ctx_ns_id FROM topologies_org_cloud_ctx_ns WHERE cloud_provider=$2 AND context=$3 AND region=$4 AND namespace=$5) FROM cte_unplaced_validators
				  `
	log.Debug().Interface("SelectInsertUnplacedValidators", q.LogHeader(ModelName))
	r, err := apps.Pg.Exec(ctx, q.RawQuery, validatorServiceInfo.ProtocolNetworkID, cloudCtxNs.CloudProvider, cloudCtxNs.Context, cloudCtxNs.Region, cloudCtxNs.Namespace)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(ModelName)); returnErr != nil {
		return err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("SelectUnplacedValidators: %s, Rows Affected: %d", q.LogHeader(ModelName), rowsAffected)
	return misc.ReturnIfErr(err, q.LogHeader(ModelName))
}

const HydraAddress = "http://zeus-hydra:9000"

func SelectValidatorsAssignedToCloudCtxNs(ctx context.Context, validatorServiceInfo ValidatorServiceCloudCtxNsProtocol, cloudCtxNs zeus_common_types.CloudCtxNs) ([]ethereum_web3signer_actions.LighthouseWeb3SignerRequest, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `	
				  SELECT vsg.pubkey, vsg.fee_recipient
				  FROM validators_service_org_groups_cloud_ctx_ns vctx
				  INNER JOIN topologies_org_cloud_ctx_ns topctx ON topctx.cloud_ctx_ns_id = vctx.cloud_ctx_ns_id
				  INNER JOIN validators_service_org_groups vsg ON vsg.pubkey = vctx.pubkey
				  WHERE vsg.protocol_network_id=$1 AND vsg.enabled=true AND topctx.cloud_provider=$2 AND topctx.context=$3 AND topctx.region=$4 AND topctx.namespace=$5
				  `
	log.Debug().Interface("SelectValidatorsAssignedToCloudCtxNs", q.LogHeader(ModelName))
	rows, err := apps.Pg.Query(ctx, q.RawQuery, validatorServiceInfo.ProtocolNetworkID, cloudCtxNs.CloudProvider, cloudCtxNs.Context, cloudCtxNs.Region, cloudCtxNs.Namespace)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(ModelName)); returnErr != nil {
		return nil, err
	}
	lhRemoteRequests := []ethereum_web3signer_actions.LighthouseWeb3SignerRequest{}
	defer rows.Close()
	for rows.Next() {
		w3rs := ethereum_web3signer_actions.LighthouseWeb3SignerRequest{
			Enable: true,
			Url:    HydraAddress,
		}
		rowErr := rows.Scan(
			&w3rs.VotingPublicKey, &w3rs.SuggestedFeeRecipient,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(ModelName))
			return nil, rowErr
		}
		lhRemoteRequests = append(lhRemoteRequests, w3rs)
	}
	return lhRemoteRequests, misc.ReturnIfErr(err, q.LogHeader(ModelName))
}

package artemis_validator_service_groups_models

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	ethereum_web3signer_actions "github.com/zeus-fyi/zeus/cookbooks/ethereum/web3signers/actions"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

type ValidatorServiceCloudCtxNsProtocol struct {
	ProtocolNetworkID, ValidatorClientNumber int
	zeus_common_types.CloudCtxNs
}

const ModelName = "ArtemisSelectInsertUnplacedValidators"

// TODO note, this selects ALL validators that aren't placed, revise once capacity observation is setup

func SelectInsertUnplacedValidatorsIntoCloudCtxNs(ctx context.Context, validatorServiceInfo ValidatorServiceCloudCtxNsProtocol) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `WITH cte_unplaced_validators AS (		
					  SELECT pubkey, fee_recipient
					  FROM validators_service_org_groups
					  WHERE NOT EXISTS (SELECT pubkey FROM validators_service_org_groups_cloud_ctx_ns) AND enabled=true AND protocol_network_id=$1
				  ) INSERT INTO validators_service_org_groups_cloud_ctx_ns(pubkey, cloud_ctx_ns_id)
					SELECT pubkey, (SELECT cloud_ctx_ns_id FROM topologies_org_cloud_ctx_ns WHERE cloud_provider=$2 AND context=$3 AND region=$4 AND namespace=$5) FROM cte_unplaced_validators
				  `
	log.Debug().Interface("SelectInsertUnplacedValidators", q.LogHeader(ModelName))
	r, err := apps.Pg.Exec(ctx, q.RawQuery, validatorServiceInfo.ProtocolNetworkID, validatorServiceInfo.CloudCtxNs.CloudProvider, validatorServiceInfo.CloudCtxNs.Context, validatorServiceInfo.CloudCtxNs.Region, validatorServiceInfo.Namespace)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(ModelName)); returnErr != nil {
		return err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("SelectUnplacedValidators: %s, Rows Affected: %d", q.LogHeader(ModelName), rowsAffected)
	return misc.ReturnIfErr(err, q.LogHeader(ModelName))
}

const HydraAddress = "http://zeus-hydra:9000"

func SelectValidatorsAssignedToCloudCtxNs(ctx context.Context, validatorServiceInfo ValidatorServiceCloudCtxNsProtocol) ([]ethereum_web3signer_actions.LighthouseWeb3SignerRequest, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `	
				  SELECT vsg.pubkey, vsg.fee_recipient
				  FROM validators_service_org_groups_cloud_ctx_ns vctx
				  INNER JOIN topologies_org_cloud_ctx_ns topctx ON topctx.cloud_ctx_ns_id = vctx.cloud_ctx_ns_id
				  INNER JOIN validators_service_org_groups vsg ON vsg.pubkey = vctx.pubkey
				  WHERE vsg.protocol_network_id=$1 AND vsg.enabled=true AND topctx.cloud_provider=$2 AND topctx.context=$3 AND topctx.region=$4 AND topctx.namespace=$5
				  `
	log.Debug().Interface("SelectValidatorsAssignedToCloudCtxNs", q.LogHeader(ModelName))
	rows, err := apps.Pg.Query(ctx, q.RawQuery, validatorServiceInfo.ProtocolNetworkID, validatorServiceInfo.CloudCtxNs.CloudProvider, validatorServiceInfo.CloudCtxNs.Context, validatorServiceInfo.CloudCtxNs.Region, validatorServiceInfo.Namespace)
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

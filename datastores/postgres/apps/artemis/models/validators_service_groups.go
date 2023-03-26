package artemis_validator_service_groups_models

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	ethereum_web3signer_actions "github.com/zeus-fyi/zeus/cookbooks/ethereum/web3signers/actions"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

type ValidatorServiceCloudCtxNsProtocol struct {
	ProtocolNetworkID     int `json:"protocolNetworkID"`
	OrgID                 int `json:"orgID"`
	ValidatorClientNumber int `json:"validatorClientNumber"`
}

const ModelName = "ArtemisValidatorsServices"

func SelectValidatorsServiceInfo(ctx context.Context, orgID int) (hestia_autogen_bases.ValidatorServiceOrgGroupSlice, error) {
	vsr := hestia_autogen_bases.ValidatorServiceOrgGroupSlice{}
	q := sql_query_templates.QueryParams{}
	params := []interface{}{orgID}
	q.RawQuery = `
				  SELECT vsg.group_name, vsg.pubkey, vsg.fee_recipient, vsg.protocol_network_id, vsg.enabled, vsg.mev_enabled
				  FROM validators_service_org_groups vsg
				  WHERE vsg.org_id=$1`
	rows, err := apps.Pg.Query(ctx, q.RawQuery, params...)
	log.Debug().Interface("SelectValidatorsServiceInfo", q.LogHeader(ModelName))
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(ModelName)); returnErr != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		vr := hestia_autogen_bases.ValidatorServiceOrgGroup{}
		rowErr := rows.Scan(
			&vr.GroupName, &vr.Pubkey, &vr.FeeRecipient, &vr.ProtocolNetworkID, &vr.Enabled, &vr.MevEnabled,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(ModelName))
			return nil, rowErr
		}
		vsr = append(vsr, vr)
	}
	return vsr, misc.ReturnIfErr(err, q.LogHeader(ModelName))
}

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
	MevEnabled        bool   `json:"mevEnabled"`
}

func FilterKeysThatExistAlready(ctx context.Context, pubkeys hestia_req_types.ValidatorServiceOrgGroupSlice) (hestia_req_types.ValidatorServiceOrgGroupSlice, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `	
				  WITH cte_check_keys AS (
   		 			  SELECT column1 AS pubkey
    				  FROM UNNEST($1::text[]) AS column1
				  ) SELECT pubkey FROM cte_check_keys
					WHERE NOT EXISTS (
						  SELECT pubkey 
						  FROM validators_service_org_groups 
						  WHERE pubkey = cte_check_keys.pubkey
					)
				  `
	log.Debug().Interface("FilterKeysThatExistAlready", q.LogHeader(ModelName))
	var pkSlice []interface{}
	for _, keyPair := range pubkeys {
		pkSlice = append(pkSlice, keyPair.Pubkey)
	}
	pubkeys = hestia_req_types.ValidatorServiceOrgGroupSlice{}
	rows, err := apps.Pg.Query(ctx, q.RawQuery, pq.Array(pkSlice))
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(ModelName)); returnErr != nil {
		return pubkeys, err
	}
	defer rows.Close()
	for rows.Next() {
		var dbPubKey string
		err = rows.Scan(&dbPubKey)
		if returnErr := misc.ReturnIfErr(err, q.LogHeader(ModelName)); returnErr != nil {
			return pubkeys, err
		}
		pubkey := hestia_req_types.ValidatorServiceOrgGroup{Pubkey: dbPubKey}
		pubkeys = append(pubkeys, pubkey)
	}
	return pubkeys, nil
}
func InsertVerifiedValidatorsToService(ctx context.Context, validatorServiceInfo OrgValidatorService, pubkeys hestia_req_types.ValidatorServiceOrgGroupSlice) error {
	checkedKeys, err := FilterKeysThatExistAlready(ctx, pubkeys)
	if err != nil {
		log.Ctx(ctx).Err(err)
		return err
	}

	m := make(map[string]string)
	for _, key := range checkedKeys {
		m[key.Pubkey] = key.Pubkey
	}

	newKeys := hestia_req_types.ValidatorServiceOrgGroupSlice{}
	for _, pk := range pubkeys {
		if _, ok := m[pk.Pubkey]; !ok {
			continue
		}
		newKeys = append(newKeys, pk)
	}

	var rows [][]interface{}
	for _, keyPair := range newKeys {
		rows = append(rows, []interface{}{
			validatorServiceInfo.GroupName,
			validatorServiceInfo.OrgID,
			keyPair.Pubkey,
			validatorServiceInfo.ProtocolNetworkID,
			keyPair.FeeRecipient,
			validatorServiceInfo.Enabled,
			validatorServiceInfo.ServiceURL,
			validatorServiceInfo.MevEnabled,
		})
	}
	columns := []string{"group_name", "org_id", "pubkey", "protocol_network_id", "fee_recipient", "enabled", "service_url", "mev_enabled"}
	// Use the `pgx.CopyFrom` method to insert the data into the table
	_, err = apps.Pg.Pgpool.CopyFrom(ctx, pgx.Identifier{"validators_service_org_groups"}, columns, pgx.CopyFromRows(rows))
	if err != nil {
		log.Ctx(ctx).Err(err)
		return err
	}
	return err
}

// SelectInsertUnplacedValidatorsIntoCloudCtxNs TODO needs to also use capacity and client number assignments
func SelectInsertUnplacedValidatorsIntoCloudCtxNs(ctx context.Context, validatorServiceInfo ValidatorServiceCloudCtxNsProtocol, cloudCtxNs zeus_common_types.CloudCtxNs) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `WITH cte_unplaced_validators AS (		
					  SELECT pubkey, fee_recipient 
					  FROM validators_service_org_groups vg
					  WHERE NOT EXISTS (
					  	SELECT 1
					  	FROM validators_service_org_groups_cloud_ctx_ns cg
					  	WHERE vg.pubkey = cg.pubkey AND enabled=true
					  ) AND protocol_network_id=$1
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

// SelectValidatorsAssignedToCloudCtxNs is used by athena
func SelectValidatorsAssignedToCloudCtxNs(ctx context.Context, validatorServiceInfo ValidatorServiceCloudCtxNsProtocol, cloudCtxNs zeus_common_types.CloudCtxNs) ([]ethereum_web3signer_actions.LighthouseWeb3SignerRequest, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `	
				  SELECT vsg.enabled, vsg.pubkey, vsg.fee_recipient, vsg.mev_enabled
				  FROM validators_service_org_groups_cloud_ctx_ns vctx
				  INNER JOIN topologies_org_cloud_ctx_ns topctx ON topctx.cloud_ctx_ns_id = vctx.cloud_ctx_ns_id
				  INNER JOIN validators_service_org_groups vsg ON vsg.pubkey = vctx.pubkey
				  WHERE vsg.protocol_network_id=$1 AND vsg.enabled=true AND topctx.cloud_provider=$2 AND topctx.context=$3 AND topctx.region=$4 AND topctx.namespace=$5 AND vctx.validator_client_number=$6
				  `
	log.Debug().Interface("SelectValidatorsAssignedToCloudCtxNs", q.LogHeader(ModelName))
	rows, err := apps.Pg.Query(ctx, q.RawQuery, validatorServiceInfo.ProtocolNetworkID, cloudCtxNs.CloudProvider, cloudCtxNs.Context, cloudCtxNs.Region, cloudCtxNs.Namespace, validatorServiceInfo.ValidatorClientNumber)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(ModelName)); returnErr != nil {
		return nil, err
	}
	lhRemoteRequests := []ethereum_web3signer_actions.LighthouseWeb3SignerRequest{}
	defer rows.Close()
	for rows.Next() {
		w3rs := ethereum_web3signer_actions.LighthouseWeb3SignerRequest{
			Type: ethereum_web3signer_actions.Web3SignerType,
			Url:  HydraAddress,
		}
		rowErr := rows.Scan(
			&w3rs.Enabled, &w3rs.VotingPublicKey, &w3rs.SuggestedFeeRecipient, &w3rs.BuilderProposals,
		)
		// No mev for ephemery
		if validatorServiceInfo.ProtocolNetworkID == hestia_req_types.EthereumEphemeryProtocolNetworkID {
			w3rs.BuilderProposals = false
		}
		w3rs.Graffiti = "zeusFyiServerlessValidators"
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(ModelName))
			return nil, rowErr
		}
		lhRemoteRequests = append(lhRemoteRequests, w3rs)
	}
	return lhRemoteRequests, misc.ReturnIfErr(err, q.LogHeader(ModelName))
}

package hestia_quicknode_models

import (
	"context"

	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type QuickNodeService struct {
	hestia_autogen_bases.ProvisionedQuickNodeServices
	ProvisionedQuicknodeServicesContractAddresses []hestia_autogen_bases.ProvisionedQuicknodeServicesContractAddresses
	ProvisionedQuicknodeServicesReferers          []hestia_autogen_bases.ProvisionedQuicknodeServicesReferers
}

/*
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
*/

/*
 WITH cte_check_keys AS (
  SELECT column1 AS pubkey
  FROM UNNEST($1::text[]) AS column1
*/

func InsertProvisionedQuickNodeService(ctx context.Context, ps QuickNodeService) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `WITH cte_insert_service AS (
					  INSERT INTO provisioned_quicknode_services(quicknode_id, endpoint_id, http_url, network, plan, active, org_id, wss_url, chain)
					  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
					  ON CONFLICT (quicknode_id) 
					  DO UPDATE SET 
					  updated_at = EXCLUDED.updated_at,
					  endpoint_id = EXCLUDED.endpoint_id,
					  http_url = EXCLUDED.http_url,
					  network = EXCLUDED.network,
					  plan = EXCLUDED.plan,
					  active = EXCLUDED.active,
					  org_id = EXCLUDED.org_id,
					  wss_url = EXCLUDED.wss_url,
					  chain = EXCLUDED.chain
					  RETURNING quicknode_id, endpoint_id
				  ), cte_unnest_ca AS (
					  SELECT column1 AS contract_address
 					  FROM UNNEST($10::text[]) AS column1
				  ), cte_insert_contract_addresses AS (
					  INSERT INTO provisioned_quicknode_services_contract_addresses(quicknode_id, endpoint_id, contract_address)
					  SELECT cte_insert_service.quicknode_id, cte_insert_service.endpoint_id, cte_unnest_ca.contract_address
					  FROM cte_insert_service, cte_unnest_ca
				  ), cte_unnest_ref AS (
					  SELECT column1 AS referer
 					  FROM UNNEST($11::text[]) AS column1
				  ), cte_insert_referers AS (
					  INSERT INTO provisioned_quicknode_services_referers(quicknode_id, endpoint_id, referer)
					  SELECT cte_insert_service.quicknode_id, cte_insert_service.endpoint_id, cte_unnest_ref.referer
					  FROM cte_insert_service, cte_unnest_ref
				  ) SELECT quicknode_id FROM cte_insert_service;`

	cas := make([]string, len(ps.ProvisionedQuicknodeServicesContractAddresses))
	for _, ca := range ps.ProvisionedQuicknodeServicesContractAddresses {
		cas = append(cas, ca.ContractAddress)
	}
	refs := make([]string, len(ps.ProvisionedQuicknodeServicesReferers))
	for _, ref := range ps.ProvisionedQuicknodeServicesReferers {
		refs = append(refs, ref.Referer)
	}
	result, err := apps.Pg.Exec(ctx, q.RawQuery, ps.QuickNodeID, ps.EndpointID, ps.HttpURL, ps.Network, ps.Plan, ps.Active, ps.OrgID, ps.WssURL, ps.Chain,
		pq.Array(cas), pq.Array(refs))
	if err != nil {
		// Log the error here using ZeroLog
		log.Error().Err(err).Msg("failed to execute query")
	} else {
		log.Info().Msg("query executed successfully.")
		// You can inspect result here
		rowsAffected := result.RowsAffected()
		log.Info().Int("rows_affected", int(rowsAffected)).Msg("number of rows affected")
	}
	return misc.ReturnIfErr(err, q.LogHeader("InsertProvisionedQuickNodeService"))
}

func UpdateProvisionedQuickNodeService(ctx context.Context, ps hestia_autogen_bases.ProvisionedQuickNodeServices) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `UPDATE provisioned_quicknode_services
    			  SET endpoint_id = $1, http_url = $2, network = $3, plan = $4, wss_url = $5, chain = $6
				  WHERE org_id = $1 AND quicknode_id = $2
			      RETURNING quicknode_id;
				  `
	qnID := ""
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, ps.OrgID, ps.EndpointID, ps.HttpURL, ps.Network, ps.Plan, ps.WssURL, ps.Chain).Scan(&qnID)
	if err != nil {
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader("UpdateProvisionedQuickNodeService"))
}

func DeactivateProvisionedQuickNodeService(ctx context.Context, quickNodeID, endpointID string, ou org_users.OrgUser) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `UPDATE provisioned_quicknode_services
    			  SET active = false
				  WHERE org_id = $1 AND quicknode_id = $2
			      RETURNING quicknode_id;
				  `
	qnID := ""
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, ou.OrgID, quickNodeID, endpointID).Scan(&qnID)
	if err != nil {
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader("DeactivateProvisionedQuickNodeService"))
}

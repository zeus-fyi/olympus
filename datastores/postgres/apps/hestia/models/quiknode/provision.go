package hestia_quicknode_models

import (
	"context"

	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type QuickNodeService struct {
	hestia_autogen_bases.ProvisionedQuickNodeServices
	ProvisionedQuicknodeServicesContractAddresses []hestia_autogen_bases.ProvisionedQuicknodeServicesContractAddresses
	ProvisionedQuicknodeServicesReferers          []hestia_autogen_bases.ProvisionedQuicknodeServicesReferers
}

func InsertProvisionedQuickNodeService(ctx context.Context, ps QuickNodeService) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `WITH cte_insert_service AS (
					  INSERT INTO provisioned_quicknode_services(quicknode_id, endpoint_id, http_url, network, plan, active, org_id, wss_url, chain)
					  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
					  ON CONFLICT (quicknode_id, endpoint_id) 
					  DO UPDATE SET 
					  http_url = EXCLUDED.http_url,
					  network = EXCLUDED.network,
					  plan = EXCLUDED.plan,
					  active = EXCLUDED.active,
					  org_id = EXCLUDED.org_id,
					  wss_url = EXCLUDED.wss_url,
					  chain = EXCLUDED.chain
					  RETURNING quicknode_id, endpoint_id
				  ), cte_delete_ca AS (
					  DELETE FROM provisioned_quicknode_services_contract_addresses
					  WHERE quicknode_id = (SELECT quicknode_id FROM cte_insert_service) AND endpoint_id = (SELECT endpoint_id FROM cte_insert_service)
				  ), cte_unnest_ca AS (
					  SELECT column1 AS contract_address
 					  FROM UNNEST($10::text[]) AS column1
				  ), cte_insert_contract_addresses AS (
					  INSERT INTO provisioned_quicknode_services_contract_addresses(quicknode_id, endpoint_id, contract_address)
					  SELECT cte_insert_service.quicknode_id, cte_insert_service.endpoint_id, cte_unnest_ca.contract_address
					  FROM cte_insert_service, cte_unnest_ca
					  WHERE cte_unnest_ca.contract_address IS NOT NULL AND cte_unnest_ca.contract_address != '' 
				  ), cte_delete_ref AS (
					  DELETE FROM provisioned_quicknode_services_referers
					  WHERE quicknode_id = (SELECT quicknode_id FROM cte_insert_service) AND endpoint_id = (SELECT endpoint_id FROM cte_insert_service)
				  ), cte_unnest_ref AS (
					  SELECT column1 AS referer
 					  FROM UNNEST($11::text[]) AS column1
				  ), cte_insert_referers AS (
					  INSERT INTO provisioned_quicknode_services_referers(quicknode_id, endpoint_id, referer)
					  SELECT cte_insert_service.quicknode_id, cte_insert_service.endpoint_id, cte_unnest_ref.referer
					  FROM cte_insert_service, cte_unnest_ref
					  WHERE cte_unnest_ref.referer IS NOT NULL AND cte_unnest_ref.referer != ''
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

func UpdateProvisionedQuickNodeService(ctx context.Context, ps QuickNodeService) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `WITH cte_insert_service AS (
					  UPDATE provisioned_quicknode_services
					  SET http_url = $4, wss_url = $5, plan = $6
					  WHERE org_id = $1 AND quicknode_id = $2 AND endpoint_id = $3
					  RETURNING quicknode_id, endpoint_id
				  ), cte_delete_ca AS (
					  DELETE FROM provisioned_quicknode_services_contract_addresses
					  WHERE quicknode_id = (SELECT quicknode_id FROM cte_insert_service) AND endpoint_id = (SELECT endpoint_id FROM cte_insert_service)
				  ), cte_unnest_ca AS (
					  SELECT column1 AS contract_address
 					  FROM UNNEST($7::text[]) AS column1
				  ), cte_insert_contract_addresses AS (
					  INSERT INTO provisioned_quicknode_services_contract_addresses(quicknode_id, endpoint_id, contract_address)
					  SELECT cte_insert_service.quicknode_id, cte_insert_service.endpoint_id, cte_unnest_ca.contract_address
					  FROM cte_insert_service, cte_unnest_ca
					  WHERE cte_unnest_ca.contract_address IS NOT NULL AND cte_unnest_ca.contract_address != ''
				  ), cte_delete_ref AS (
					  DELETE FROM provisioned_quicknode_services_referers
					  WHERE quicknode_id = (SELECT quicknode_id FROM cte_insert_service) AND endpoint_id = (SELECT endpoint_id FROM cte_insert_service)
				  ), cte_unnest_ref AS (
					  SELECT column1 AS referer
 					  FROM UNNEST($8::text[]) AS column1
				  ), cte_insert_referers AS (
					  INSERT INTO provisioned_quicknode_services_referers(quicknode_id, endpoint_id, referer)
					  SELECT cte_insert_service.quicknode_id, cte_insert_service.endpoint_id, cte_unnest_ref.referer
					  FROM cte_insert_service, cte_unnest_ref
					  WHERE cte_unnest_ref.referer IS NOT NULL AND cte_unnest_ref.referer != ''
				  ) SELECT quicknode_id FROM cte_insert_service;`
	cas := make([]string, len(ps.ProvisionedQuicknodeServicesContractAddresses))
	for _, ca := range ps.ProvisionedQuicknodeServicesContractAddresses {
		cas = append(cas, ca.ContractAddress)
	}
	refs := make([]string, len(ps.ProvisionedQuicknodeServicesReferers))
	for _, ref := range ps.ProvisionedQuicknodeServicesReferers {
		refs = append(refs, ref.Referer)
	}
	qnID := ""
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, ps.OrgID, ps.QuickNodeID, ps.EndpointID, ps.HttpURL, ps.WssURL, ps.Plan,
		pq.Array(cas), pq.Array(refs)).Scan(&qnID)
	if err != nil {
		log.Error().Err(err).Msg("UpdateProvisionedQuickNodeService: failed to execute query")
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader("UpdateProvisionedQuickNodeService"))
}

func DeactivateProvisionedQuickNodeServiceEndpoint(ctx context.Context, orgID int, quickNodeID, endpointID string) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `UPDATE provisioned_quicknode_services
    			  SET active = false
				  WHERE org_id = $1 AND quicknode_id = $2 AND endpoint_id = $3
			      RETURNING quicknode_id;
				  `
	qnID := ""
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, orgID, quickNodeID, endpointID).Scan(&qnID)
	if err != nil {
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader("DeactivateProvisionedQuickNodeServiceEndpoint"))
}

func DeprovisionQuickNodeServices(ctx context.Context, orgID int, quickNodeID string) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `UPDATE provisioned_quicknode_services
    			  SET active = false
				  WHERE org_id = $1 AND quicknode_id = $2
			      RETURNING quicknode_id;
				  `
	qnID := ""
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, orgID, quickNodeID).Scan(&qnID)
	if err != nil {
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader("DeactivateProvisionedQuickNodeServiceEndpoint"))
}

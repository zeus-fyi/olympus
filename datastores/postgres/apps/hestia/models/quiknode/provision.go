package hestia_quicknode_models

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type QuickNodeService struct {
	IsTest bool
	hestia_autogen_bases.ProvisionedQuickNodeServices
	ProvisionedQuicknodeServicesContractAddresses []hestia_autogen_bases.ProvisionedQuicknodeServicesContractAddresses
	ProvisionedQuicknodeServicesReferers          []hestia_autogen_bases.ProvisionedQuicknodeServicesReferers
}

func InsertProvisionedQuickNodeService(ctx context.Context, ps QuickNodeService) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `WITH cte_insert_org AS (
					  INSERT INTO orgs (name, metadata)
					  SELECT $1, '{}'
						WHERE NOT EXISTS (
							SELECT 1
							FROM orgs
							WHERE name = $1
						)
					   RETURNING org_id, name AS quicknode_id
					), cte_marketplace_customer AS (
						  INSERT INTO quicknode_marketplace_customer (quicknode_id, plan, is_test)
						  SELECT $1, $2, $11
						  FROM cte_insert_org
						  ON CONFLICT (quicknode_id) 
						  DO UPDATE SET 
						  plan = EXCLUDED.plan
				    ), cte_insert_service AS (
					  INSERT INTO provisioned_quicknode_services(quicknode_id, endpoint_id, http_url, network, active, wss_url, chain)
					  VALUES ($1, $3, $4, $5, $6, $7, $8)
					  ON CONFLICT (quicknode_id, endpoint_id) 
					  DO UPDATE SET 
					  http_url = EXCLUDED.http_url,
					  network = EXCLUDED.network,
					  active = EXCLUDED.active,
					  wss_url = EXCLUDED.wss_url,
					  chain = EXCLUDED.chain
					  RETURNING quicknode_id, endpoint_id
				  ), cte_delete_ca AS (
					  DELETE FROM provisioned_quicknode_services_contract_addresses
					  WHERE endpoint_id = (SELECT endpoint_id FROM cte_insert_service)
				  ), cte_unnest_ca AS (
					  SELECT column1 AS contract_address
 					  FROM UNNEST($9::text[]) AS column1
				  ), cte_insert_contract_addresses AS (
					  INSERT INTO provisioned_quicknode_services_contract_addresses(endpoint_id, contract_address)
					  SELECT cte_insert_service.endpoint_id, cte_unnest_ca.contract_address
					  FROM cte_insert_service, cte_unnest_ca
					  WHERE cte_unnest_ca.contract_address IS NOT NULL AND cte_unnest_ca.contract_address != '' 
 			 		  ON CONFLICT (endpoint_id) DO UPDATE SET contract_address = EXCLUDED.contract_address
				  ), cte_delete_ref AS (
					  DELETE FROM provisioned_quicknode_services_referers
					  WHERE endpoint_id = (SELECT endpoint_id FROM cte_insert_service)
				  ), cte_unnest_ref AS (
					  SELECT column1 AS referer
 					  FROM UNNEST($10::text[]) AS column1
				  ), cte_insert_referers AS (
					  INSERT INTO provisioned_quicknode_services_referers(endpoint_id, referer)
					  SELECT cte_insert_service.endpoint_id, cte_unnest_ref.referer
					  FROM cte_insert_service, cte_unnest_ref
					  WHERE cte_unnest_ref.referer IS NOT NULL AND cte_unnest_ref.referer != ''
					  ON CONFLICT (endpoint_id) DO UPDATE SET referer = EXCLUDED.referer
				  ) SELECT quicknode_id FROM cte_insert_service;`
	cas := make([]string, len(ps.ProvisionedQuicknodeServicesContractAddresses))
	for _, ca := range ps.ProvisionedQuicknodeServicesContractAddresses {
		cas = append(cas, ca.ContractAddress)
	}
	refs := make([]string, len(ps.ProvisionedQuicknodeServicesReferers))
	for _, ref := range ps.ProvisionedQuicknodeServicesReferers {
		refs = append(refs, ref.Referer)
	}
	result, err := apps.Pg.Exec(ctx, q.RawQuery, ps.QuickNodeID, ps.Plan, ps.EndpointID, ps.HttpURL, ps.Network, ps.Active, ps.WssURL, ps.Chain,
		pq.Array(cas), pq.Array(refs), ps.IsTest)
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
	q.RawQuery = `WITH cte_insert_plan AS (
					 UPDATE quicknode_marketplace_customer
					 SET plan = $2
					 WHERE quicknode_id = $1
 					 RETURNING quicknode_id
				  ), cte_update_service AS (
					  UPDATE provisioned_quicknode_services
					  SET http_url = $4, wss_url = $5
					  WHERE quicknode_id = $1 AND endpoint_id = $3
					  RETURNING quicknode_id, endpoint_id
				  ), cte_delete_ca AS (
					  DELETE FROM provisioned_quicknode_services_contract_addresses
					  WHERE endpoint_id = $3
				  ), cte_unnest_ca AS (
					  SELECT column1 AS contract_address
 					  FROM UNNEST($6::text[]) AS column1
				  ), cte_insert_contract_addresses AS (
					  INSERT INTO provisioned_quicknode_services_contract_addresses(endpoint_id, contract_address)
					  SELECT $3, cte_unnest_ca.contract_address
					  FROM cte_unnest_ca
					  WHERE cte_unnest_ca.contract_address IS NOT NULL AND cte_unnest_ca.contract_address != ''
					  ON CONFLICT (endpoint_id) DO UPDATE SET contract_address = EXCLUDED.contract_address	
				  ), cte_delete_ref AS (
					  DELETE FROM provisioned_quicknode_services_referers
					  WHERE endpoint_id = $3
				  ), cte_unnest_ref AS (
					  SELECT column1 AS referer
 					  FROM UNNEST($7::text[]) AS column1
				  ), cte_insert_referers AS (
					  INSERT INTO provisioned_quicknode_services_referers(endpoint_id, referer)
					  SELECT $3, cte_unnest_ref.referer
					  FROM cte_unnest_ref
					  WHERE cte_unnest_ref.referer IS NOT NULL AND cte_unnest_ref.referer != ''
					  ON CONFLICT (endpoint_id) DO UPDATE SET referer = EXCLUDED.referer
				  ) SELECT quicknode_id FROM cte_insert_plan;`
	cas := make([]string, len(ps.ProvisionedQuicknodeServicesContractAddresses))
	for _, ca := range ps.ProvisionedQuicknodeServicesContractAddresses {
		cas = append(cas, ca.ContractAddress)
	}
	refs := make([]string, len(ps.ProvisionedQuicknodeServicesReferers))
	for _, ref := range ps.ProvisionedQuicknodeServicesReferers {
		refs = append(refs, ref.Referer)
	}
	qnID := ""
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, ps.QuickNodeID, ps.Plan, ps.EndpointID, ps.HttpURL, ps.WssURL,
		pq.Array(cas), pq.Array(refs)).Scan(&qnID)
	if err != nil {
		log.Error().Err(err).Msg("UpdateProvisionedQuickNodeService: failed to execute query")
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader("UpdateProvisionedQuickNodeService"))
}

func DeactivateProvisionedQuickNodeServiceEndpoint(ctx context.Context, quickNodeID, endpointID string) (string, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `UPDATE provisioned_quicknode_services
    			  SET active = false
				  WHERE quicknode_id = $1 AND endpoint_id = $2
			      RETURNING http_url;
				  `
	httpURL := ""
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, quickNodeID, endpointID).Scan(&httpURL)
	if err != nil {
		return httpURL, err
	}
	return httpURL, err
}

func DeprovisionQuickNodeServices(ctx context.Context, quickNodeID string) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `UPDATE provisioned_quicknode_services
    			  SET active = false
				  WHERE quicknode_id = $1
			      RETURNING quicknode_id;
				  `
	qnID := ""
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, quickNodeID).Scan(&qnID)
	if err != nil {
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader("DeactivateProvisionedQuickNodeServiceEndpoint"))
}

type QuickNodeServicedEndpoints struct {
	Plan        string
	EndpointMap map[string]QuickNodeService
}

func SelectQuickNodeServicesByQid(ctx context.Context, qId string) (QuickNodeServicedEndpoints, error) {
	q := `  SELECT qps.endpoint_id, qps.http_url, qps.network, qps.wss_url, qps.chain, 
         		(SELECT plan FROM public.quicknode_marketplace_customer qmc WHERE qmc.quicknode_id = $1) AS plan,
         		ca.contract_address, ref.referer
			FROM provisioned_quicknode_services qps
			LEFT JOIN provisioned_quicknode_services_contract_addresses ca ON ca.endpoint_id = qps.endpoint_id
			LEFT JOIN provisioned_quicknode_services_referers ref ON ref.endpoint_id = qps.endpoint_id
			WHERE quicknode_id = $1 AND active = true`

	args := []interface{}{
		qId,
	}
	var qnse QuickNodeServicedEndpoints
	qs := make(map[string]QuickNodeService)
	// Execute the SQL query
	rows, err := apps.Pg.Query(ctx, q, args...)
	if err != nil {
		return qnse, err
	}
	defer rows.Close()
	for rows.Next() {
		qns := QuickNodeService{
			ProvisionedQuickNodeServices: hestia_autogen_bases.ProvisionedQuickNodeServices{
				QuickNodeID: qId,
				Active:      true,
			},
			ProvisionedQuicknodeServicesContractAddresses: []hestia_autogen_bases.ProvisionedQuicknodeServicesContractAddresses{},
			ProvisionedQuicknodeServicesReferers:          []hestia_autogen_bases.ProvisionedQuicknodeServicesReferers{},
		}
		var cadr, refa sql.NullString
		err = rows.Scan(
			&qns.EndpointID,
			&qns.HttpURL,
			&qns.Network,
			&qns.WssURL,
			&qns.Chain,
			&qnse.Plan,
			&cadr,
			&refa,
		)
		if err != nil {
			return qnse, err
		}
		if cadr.Valid {
			qns.ProvisionedQuicknodeServicesContractAddresses = append(qns.ProvisionedQuicknodeServicesContractAddresses, hestia_autogen_bases.ProvisionedQuicknodeServicesContractAddresses{
				ContractAddress: cadr.String,
			})
		}
		if refa.Valid {
			qns.ProvisionedQuicknodeServicesReferers = append(qns.ProvisionedQuicknodeServicesReferers, hestia_autogen_bases.ProvisionedQuicknodeServicesReferers{
				Referer: refa.String,
			})
		}
		if val, ok := qs[qns.EndpointID]; ok {
			val.ProvisionedQuicknodeServicesContractAddresses = append(val.ProvisionedQuicknodeServicesContractAddresses, qns.ProvisionedQuicknodeServicesContractAddresses...)
			val.ProvisionedQuicknodeServicesReferers = append(val.ProvisionedQuicknodeServicesReferers, qns.ProvisionedQuicknodeServicesReferers...)
			qs[qns.EndpointID] = val
		} else {
			qs[qns.EndpointID] = qns
		}
	}
	if err = rows.Err(); err != nil {
		return qnse, err
	}
	qnse.EndpointMap = qs
	return qnse, nil
}

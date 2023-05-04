package hestia_compute_resources

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func InsertGcpServices(ctx context.Context, service autogen_bases.GcpServicesSlice) error {
	q := sql_query_templates.NewQueryParam("InsertGcpServices", "gcp_services", "where", 1000, []string{})
	cte := sql_query_templates.CTE{Name: "InsertGcpServices"}
	cte.SubCTEs = []sql_query_templates.SubCTE{}
	cte.Params = []interface{}{}

	for _, s := range service {
		queryResourceId := fmt.Sprintf("resource_id_insert_%d", ts.UnixTimeStampNow())
		scteRe := sql_query_templates.NewSubInsertCTE(queryResourceId)
		scteRe.TableName = s.GetTableName()
		scteRe.Columns = s.GetTableColumns()
		scteRe.Values = []apps.RowValues{s.GetRowValues(queryResourceId)}
		cte.SubCTEs = append(cte.SubCTEs, scteRe)
	}
	q.RawQuery = cte.GenerateChainedCTE()
	r, err := apps.Pg.Exec(ctx, q.RawQuery, cte.Params...)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("Nodes: %s, Rows Affected: %d", q.LogHeader(Sn), rowsAffected)
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func InsertGcpServicesSKU(ctx context.Context, sku autogen_bases.GcpServicesSkus) error {
	q := `INSERT INTO gcp_services_skus (
			service_id,
			name,
			sku_id,
			description,
			service_display_name,
			resource_family,
			resource_group,
			usage_type,
			service_regions,
			pricing_info,
			service_provider_name,
			geo_taxonomy
		)
		VALUES (
			$1, -- service_id
			$2, -- name
			$3, -- sku_id
			$4, -- description
			$5, -- service_display_name
			$6, -- resource_family
			$7, -- resource_group
			$8, -- usage_type
			$9, -- service_regions
			$10, -- pricing_info
			$11, -- service_provider_name
			$12 -- geo_taxonomy
		)
		ON CONFLICT (sku_id)
		DO UPDATE SET
			service_id = EXCLUDED.service_id,
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			service_display_name = EXCLUDED.service_display_name,
			resource_family = EXCLUDED.resource_family,
			resource_group = EXCLUDED.resource_group,
			usage_type = EXCLUDED.usage_type,
			service_regions = EXCLUDED.service_regions,
			pricing_info = EXCLUDED.pricing_info,
			service_provider_name = EXCLUDED.service_provider_name,
			geo_taxonomy = EXCLUDED.geo_taxonomy;`
	r, err := apps.Pg.Exec(ctx, q, sku.GetRowValues("InsertGcpServicesSKU")...)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("InsertGcpServicesSKU: Error inserting into gcp_services_skus")
		return err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("InsertGcpServicesSKU: Rows Affected: %d", rowsAffected)
	return err
}

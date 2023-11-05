package artemis_reporting

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func getBundleProfitQ() string {
	var que = `
        INSERT INTO eth_mev_bundle_profit (
            bundle_hash, 
            revenue, 
            costs, 
            revenue_prediction
        )  
        VALUES (
            $1, $2, $3, $4
        )
        ON CONFLICT (bundle_hash) 
        DO UPDATE SET 
            revenue = EXCLUDED.revenue,
            revenue_prediction = EXCLUDED.revenue_prediction,
            costs = EXCLUDED.costs
        `
	return que
}

func InsertBundleProfit(ctx context.Context, bundleProfit artemis_autogen_bases.EthMevBundleProfit) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = getBundleProfitQ()
	_, err := apps.Pg.Exec(ctx, q.RawQuery, bundleProfit.BundleHash, bundleProfit.Revenue, bundleProfit.Costs, bundleProfit.RevenuePrediction)
	if err == pgx.ErrNoRows {
		err = nil
	}
	return misc.ReturnIfErr(err, q.LogHeader("InsertBundleProfit"))
}

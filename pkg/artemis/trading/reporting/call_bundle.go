package artemis_reporting

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/jackc/pgx/v4"
	"github.com/metachris/flashbotsrpc"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func getCallBundleSaveQ() string {
	var que = `
        WITH cte_mev_call AS (
			INSERT INTO events (event_id)
			VALUES ($1)
		RETURNING event_id
		) INSERT INTO eth_mev_call_bundle (event_id, builder_name, bundle_hash, protocol_network_id, eth_call_resp_json)
		   SELECT 
				event_id, 
			    $2,
				$3,
				$4, 
				$5::jsonb
			FROM cte_mev_call;
        `
	return que
}

func InsertCallBundleResp(ctx context.Context, builder string, protocolID int, callBundlesResp flashbotsrpc.FlashbotsCallBundleResponse) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = getCallBundleSaveQ()
	if callBundlesResp.BundleHash == "" {
		return errors.New("bundle hash is empty")
	}
	for i, _ := range callBundlesResp.Results {
		callBundlesResp.Results[i].Revert = strconv.QuoteToASCII(callBundlesResp.Results[i].Revert)
	}
	b, err := json.Marshal(callBundlesResp)
	if err != nil {
		log.Err(err).Msg("InsertCallBundleResp: error marshalling call bundle response")
		return err
	}
	ts := chronos.Chronos{}
	eventID := ts.UnixTimeStampNow()
	_, err = apps.Pg.Exec(ctx, q.RawQuery, eventID, builder, callBundlesResp.BundleHash, protocolID, b)
	if err == pgx.ErrNoRows {
		err = nil
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader("InsertBundleProfit"))
}

/*
WITH cte_mev_call AS (
    INSERT INTO events (event_id)
    VALUES (1)
    RETURNING event_id
)
INSERT INTO eth_mev_call_bundle (event_id, bundle_hash, protocol_network_id, eth_call_resp_json)
SELECT
    event_id,
    '0x1', -- This assumes that '0x1' is the correct hash you intended to insert. Adjust if necessary.
    1, -- Replace with your actual network ID value.
    '{}'::jsonb
FROM cte_mev_call;
*/

package artemis_reporting

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"unicode/utf8"

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

func sanitizeAndUnmarshal(data string, v interface{}) error {
	if err := json.Unmarshal([]byte(data), &v); err != nil {
		// If unmarshal fails, try to sanitize the string
		cleanData := sanitizeString(data)
		// Attempt to unmarshal again with sanitized data
		return json.Unmarshal([]byte(cleanData), &v)
	}
	return nil
}
func sanitizeString(s string) string {
	buf := make([]rune, 0, len(s))
	for _, r := range s {
		if r == utf8.RuneError {
			continue // Exclude RuneError values
		}
		buf = append(buf, r)
	}
	return string(buf)
}

func InsertCallBundleResp(ctx context.Context, builder string, protocolID int, callBundlesResp flashbotsrpc.FlashbotsCallBundleResponse) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = getCallBundleSaveQ()
	if callBundlesResp.BundleHash == "" {
		return errors.New("bundle hash is empty")
	}

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(callBundlesResp)
	if err != nil {
		return err
	}
	ts := chronos.Chronos{}
	eventID := ts.UnixTimeStampNow()
	_, err = apps.Pg.Exec(ctx, q.RawQuery, eventID, builder, callBundlesResp.BundleHash, protocolID, b.String())
	if err == pgx.ErrNoRows {
		err = nil
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader("InsertBundleProfit"))
}

func selectCallBundles() string {
	var que = `SELECT event_id, builder_name, bundle_hash, eth_call_resp_json
		  	   FROM eth_mev_call_bundle
			   WHERE event_id > $1 AND protocol_network_id = $2
			   ORDER BY event_id DESC
			   LIMIT 10000;
		  	   	`
	return que
}

type CallBundleHistory struct {
	EventID                                  int    `json:"eventID"`
	BuilderName                              string `json:"builderName"`
	flashbotsrpc.FlashbotsCallBundleResponse `json:"flashbotsCallBundleResponse"`
}

func SelectCallBundleHistory(ctx context.Context, minEventId, protocolNetworkID int) ([]CallBundleHistory, error) {
	rows, err := apps.Pg.Query(ctx, selectCallBundles(), minEventId, protocolNetworkID)
	if err != nil {
		return nil, err
	}
	var rw []CallBundleHistory
	defer rows.Close()
	for rows.Next() {
		cbh := CallBundleHistory{
			FlashbotsCallBundleResponse: flashbotsrpc.FlashbotsCallBundleResponse{},
		}
		respStr := bytes.Buffer{}
		rowErr := rows.Scan(&cbh.EventID, &cbh.BuilderName, &cbh.BundleHash, &respStr)
		if rowErr != nil {
			log.Err(rowErr).Msg("SelectCallBundleHistory")
			return nil, rowErr
		}
		err = json.Unmarshal(respStr.Bytes(), &cbh.FlashbotsCallBundleResponse)
		if err != nil {
			log.Err(err).Msg("SelectCallBundleHistory: error unmarshalling call bundle response")
			return nil, err
		}
		rw = append(rw, cbh)
	}
	return rw, nil
}

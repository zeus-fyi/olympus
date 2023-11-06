package artemis_reporting

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"unicode"

	"github.com/jackc/pgx/v4"
	"github.com/metachris/flashbotsrpc"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
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

func clearString(str string) string {
	str = strings.Replace(str, "�y�", "", -1)
	var result []rune
	for _, r := range str {
		// Keep only letters, digits, and underscore
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			result = append(result, r)
		}
	}
	return string(result)
}

func InsertCallBundleResp(ctx context.Context, builder string, protocolID int, callBundlesResp flashbotsrpc.FlashbotsCallBundleResponse) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = getCallBundleSaveQ()
	if callBundlesResp.BundleHash == "" {
		return errors.New("bundle hash is empty")
	}

	tmp := flashbotsrpc.FlashbotsCallBundleResponse{
		BundleGasPrice:    callBundlesResp.BundleGasPrice,
		BundleHash:        callBundlesResp.BundleHash,
		CoinbaseDiff:      callBundlesResp.CoinbaseDiff,
		EthSentToCoinbase: callBundlesResp.EthSentToCoinbase,
		GasFees:           callBundlesResp.GasFees,
		Results:           make([]flashbotsrpc.FlashbotsCallBundleResult, len(callBundlesResp.Results)),
		StateBlockNumber:  callBundlesResp.StateBlockNumber,
		TotalGasUsed:      callBundlesResp.TotalGasUsed,
	}

	for i, v := range callBundlesResp.Results {
		fr := flashbotsrpc.FlashbotsCallBundleResult{
			CoinbaseDiff:      v.CoinbaseDiff,
			EthSentToCoinbase: v.EthSentToCoinbase,
			FromAddress:       v.FromAddress,
			GasFees:           v.GasFees,
			GasPrice:          v.GasPrice,
			GasUsed:           v.GasUsed,
			ToAddress:         v.ToAddress,
			TxHash:            v.TxHash,
			Value:             v.Value,
			Error:             clearString(v.Error),
			Revert:            clearString(v.Revert),
		}
		tmp.Results[i] = fr
	}

	b, err := json.Marshal(tmp)
	if err != nil {
		return err
	}
	ts := chronos.Chronos{}
	eventID := ts.UnixTimeStampNow()
	_, err = apps.Pg.Exec(ctx, q.RawQuery, eventID, builder, callBundlesResp.BundleHash, protocolID, string(b))
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
			   LIMIT 1000;
		  	   	`
	return que
}

type CallBundleHistory struct {
	EventID                                  int    `json:"eventID"`
	SubmissionTime                           string `json:"submissionTime"`
	BuilderName                              string `json:"builderName"`
	flashbotsrpc.FlashbotsCallBundleResponse `json:"flashbotsCallBundleResponse"`
}

func SelectCallBundleHistory(ctx context.Context, minEventId, protocolNetworkID int) ([]CallBundleHistory, error) {
	rows, err := apps.Pg.Query(ctx, selectCallBundles(), minEventId, protocolNetworkID)
	if err != nil {
		return nil, err
	}
	var rw []CallBundleHistory
	ts := chronos.Chronos{}

	defer rows.Close()
	for rows.Next() {
		cbh := CallBundleHistory{
			FlashbotsCallBundleResponse: flashbotsrpc.FlashbotsCallBundleResponse{},
		}
		rowErr := rows.Scan(&cbh.EventID, &cbh.BuilderName, &cbh.BundleHash, &cbh.FlashbotsCallBundleResponse)
		if rowErr != nil {
			log.Err(rowErr).Msg("SelectCallBundleHistory")
			return nil, rowErr
		}
		// Create a big.Float representation of 1e18 for the division
		eth := new(big.Float).SetFloat64(1e18)
		// Divide the gas price by 1e18 to convert wei to ether
		bundleGasPriceWei := artemis_eth_units.NewBigFloatFromStr(cbh.FlashbotsCallBundleResponse.BundleGasPrice)
		bundleGasPrice, _ := new(big.Float).Quo(bundleGasPriceWei, eth).Float64()
		// Format the float to a string with 5 decimal places
		cbh.BundleGasPrice = fmt.Sprintf("%.5f", bundleGasPrice)

		bundleGasFeesWei := artemis_eth_units.NewBigFloatFromStr(cbh.FlashbotsCallBundleResponse.GasFees)
		bundleGasFees, _ := new(big.Float).Quo(bundleGasFeesWei, eth).Float64()
		cbh.GasFees = fmt.Sprintf("%.5f", bundleGasFees)

		ethCoinbaseWei := artemis_eth_units.NewBigFloatFromStr(cbh.FlashbotsCallBundleResponse.EthSentToCoinbase)
		ethCoinbase, _ := new(big.Float).Quo(ethCoinbaseWei, eth).Float64()
		cbh.EthSentToCoinbase = fmt.Sprintf("%.5f", ethCoinbase)

		ethCoinbaseDiffWei := artemis_eth_units.NewBigFloatFromStr(cbh.FlashbotsCallBundleResponse.CoinbaseDiff)
		ethCoinbaseDiff, _ := new(big.Float).Quo(ethCoinbaseDiffWei, eth).Float64()
		cbh.CoinbaseDiff = fmt.Sprintf("%.5f", ethCoinbaseDiff)

		cbh.SubmissionTime = ts.ConvertUnixTimeStampToDate(cbh.EventID).String()
		rw = append(rw, cbh)
	}
	return rw, nil
}

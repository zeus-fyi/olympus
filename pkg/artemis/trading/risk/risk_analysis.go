package artemis_risk_analysis

import (
	"context"
	"math/big"

	"github.com/lib/pq"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	core_entities "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/uniswap_libs/uniswap_core/entities"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

type RiskAnalysis struct {
	MaxTradeSize *big.Int            `json:"maxTradeSize"`
	Token        core_entities.Token `json:"token"`

	TotalEarnedInEth *big.Int                              `json:"totalEarnedInEth"`
	TotalFailed      int                                   `json:"totalFailed"`
	TotalSuccess     int                                   `json:"totalSuccess"`
	DateRange        string                                `json:"dateRange"`
	Historical       artemis_mev_models.HistoricalAnalysis `json:"historical"`
}

func SetTradingPermission(ctx context.Context, addresses []string, pni int) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `WITH cte_addresses AS (
						SELECT unnest($1::text[]) AS addr
					)
					UPDATE erc20_token_info 
					SET trading_enabled = $1 AND protocol_network_id = $2
					WHERE address IN (SELECT addr FROM cte_addresses)`
	if pni == 0 {
		pni = hestia_req_types.EthereumMainnetProtocolNetworkID
	}
	_, err := apps.Pg.Exec(ctx, q.RawQuery, pq.Array(addresses), pni)
	if err != nil {
		return err
	}
	return err
}

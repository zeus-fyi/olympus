package artemis_risk_analysis

import (
	"math/big"

	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	core_entities "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/uniswap_libs/uniswap_core/entities"
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

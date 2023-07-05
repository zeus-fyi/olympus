package artemis_models_uniswap

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func InsertUniswapPairInfo(ctx context.Context, pair artemis_autogen_bases.UniswapPairInfo) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `INSERT INTO uniswap_pair_info(address, factory_address, fee, version, token0, token1, protocol_network_id, trading_enabled)
				  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			;`

	protocolIDNum := 1
	if pair.ProtocolNetworkID != 0 {
		protocolIDNum = pair.ProtocolNetworkID
	}
	_, err := apps.Pg.Exec(ctx, q.RawQuery, pair.Address, pair.FactoryAddress, pair.Fee, pair.Version, pair.Token0, pair.Token1, protocolIDNum, pair.TradingEnabled)
	if err == pgx.ErrNoRows {
		err = nil
	}
	return misc.ReturnIfErr(err, q.LogHeader("InsertUniswapPairInfo"))
}

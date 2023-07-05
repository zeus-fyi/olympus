package artemis_models_uniswap

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/constants"
	uniswap_pricing "github.com/zeus-fyi/olympus/pkg/artemis/trading/pricing/uniswap"
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

func InsertStandardUniswapPairInfoFromPair(ctx context.Context, pair []accounts.Address) error {
	pools, err := uniswap_pricing.NewUniswapPools(pair)
	if err != nil {
		return err
	}
	piV2 := artemis_autogen_bases.UniswapPairInfo{
		TradingEnabled:    false,
		Address:           pools.V2Pair.PairContractAddr,
		FactoryAddress:    artemis_trading_constants.UniswapV2FactoryAddress,
		Fee:               int(pools.V2Pair.GetBaseFee()),
		Version:           "v2",
		Token0:            pools.V2Pair.Token0.String(),
		Token1:            pools.V2Pair.Token1.String(),
		ProtocolNetworkID: 1,
	}
	err = InsertUniswapPairInfo(ctx, piV2)
	if err != nil {
		return err
	}
	piV3Med := artemis_autogen_bases.UniswapPairInfo{
		TradingEnabled:    false,
		Address:           pools.V3Pairs.MediumFee.PoolAddress,
		FactoryAddress:    artemis_trading_constants.UniswapV3FactoryAddress,
		Fee:               int(pools.V3Pairs.MediumFee.Fee),
		Version:           "v3",
		Token0:            pools.V2Pair.Token0.String(),
		Token1:            pools.V2Pair.Token1.String(),
		ProtocolNetworkID: 1,
	}
	err = InsertUniswapPairInfo(ctx, piV3Med)
	if err != nil {
		return err
	}
	piV3High := artemis_autogen_bases.UniswapPairInfo{
		TradingEnabled:    false,
		Address:           pools.V3Pairs.HighFee.PoolAddress,
		FactoryAddress:    artemis_trading_constants.UniswapV3FactoryAddress,
		Fee:               int(pools.V3Pairs.HighFee.Fee),
		Version:           "v3",
		Token0:            pools.V2Pair.Token0.String(),
		Token1:            pools.V2Pair.Token1.String(),
		ProtocolNetworkID: 1,
	}
	err = InsertUniswapPairInfo(ctx, piV3High)
	if err != nil {
		return err
	}
	return nil
}

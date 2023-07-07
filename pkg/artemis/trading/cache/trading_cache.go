package artemis_trading_cache

import (
	"context"

	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
	artemis_flashbots "github.com/zeus-fyi/olympus/pkg/artemis/flashbots"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

var (
	TokenMap        map[string]artemis_autogen_bases.Erc20TokenInfo
	FlashbotsClient artemis_flashbots.FlashbotsClient
)

func InitTokenFilter(ctx context.Context) {
	_, tm, terr := artemis_mev_models.SelectERC20Tokens(ctx)
	if terr != nil {
		panic(terr)
	}
	TokenMap = tm
}

func InitFlashbotsCache(ctx context.Context) {
	web3 := artemis_network_cfgs.ArtemisEthereumMainnet
	FlashbotsClient = artemis_flashbots.InitFlashbotsClient(ctx, web3.NodeURL, hestia_req_types.Mainnet, web3.Account)
}

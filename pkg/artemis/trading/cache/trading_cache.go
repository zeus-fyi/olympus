package artemis_trading_cache

import (
	"context"

	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/dynamic_secrets"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
	artemis_flashbots "github.com/zeus-fyi/olympus/pkg/artemis/flashbots"
	"github.com/zeus-fyi/olympus/pkg/athena"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/encryption"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

var (
	TradeExecutor   accounts.Account
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

func InitFlashbotsCache(ctx context.Context, age encryption.Age) {
	web3 := artemis_network_cfgs.ArtemisEthereumMainnet
	TradeExecutor = InitAccount(ctx, age)
	FlashbotsClient = artemis_flashbots.InitFlashbotsClient(ctx, web3.NodeURL, hestia_req_types.Mainnet, &TradeExecutor)
}

// InitAccount pubkey 0x000025e60C7ff32a3470be7FE3ed1666b0E326e2
func InitAccount(ctx context.Context, age encryption.Age) accounts.Account {
	p := filepaths.Path{
		DirIn:  "keygen",
		DirOut: "keygen",
		FnIn:   "key-4.txt.age",
	}
	r, err := dynamic_secrets.ReadAddress(ctx, p, athena.AthenaS3Manager, age)
	if err != nil {
		panic(err)
	}
	acc, err := dynamic_secrets.GetAccount(r)
	if err != nil {
		panic(err)
	}
	return acc
}

package artemis_trading_cache

import (
	"context"

	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
)

var TokenMap map[string]artemis_autogen_bases.Erc20TokenInfo

func InitTokenFilter(ctx context.Context) {
	_, tm, terr := artemis_validator_service_groups_models.SelectERC20Tokens(ctx)
	if terr != nil {
		panic(terr)
	}
	TokenMap = tm
}

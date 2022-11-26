package zeus_actions

import (
	"context"

	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

func (z *ZeusActionsClient) RestartZeusPods(ctx context.Context) ([]byte, error) {
	zeusBasePar.FilterOpts = &string_utils.FilterOpts{
		StartsWith: "zeus",
	}
	z.PrintReqJson(zeusBasePar)
	resp, err := z.DeletePods(ctx, zeusBasePar)
	z.PrintRespJson(resp)
	return resp, err
}

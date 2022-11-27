package aegis_actions

import (
	"context"

	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

func (z *AegisActionsClient) RestartAegisPods(ctx context.Context) ([]byte, error) {
	aegisBasePar.FilterOpts = &string_utils.FilterOpts{
		StartsWith: "aegis",
	}
	z.PrintReqJson(aegisBasePar)
	resp, err := z.DeletePods(ctx, aegisBasePar)
	z.PrintRespJson(resp)
	return resp, err
}

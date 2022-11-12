package zeus_client

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	read_infra "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/read"
)

func (z *ZeusClient) ReadChart(ctx context.Context, tar read_infra.TopologyReadRequest) (chart_workload.NativeK8s, error) {
	respJson := chart_workload.NativeK8s{}
	z.PrintReqJson(tar)
	resp, err := z.R().
		SetResult(&respJson).
		SetBody(tar).
		Post(InfraReadChartV1Path)

	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Ctx(ctx).Err(err).Msg("ZeusClient: ReadChart")
		return respJson, err
	}
	z.PrintRespJson(resp.Body())
	return respJson, err
}

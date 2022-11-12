package zeus_client

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/compression"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
	create_infra "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/create"
)

func (z *ZeusClient) UploadChart(ctx context.Context, p structs.Path, tar create_infra.TopologyCreateRequest) (create_infra.TopologyCreateResponse, error) {
	respJson := create_infra.TopologyCreateResponse{}
	err := z.ZipK8sChartToPath(&p)
	if err != nil {
		return respJson, err
	}
	z.PrintReqJson(tar)
	resp, err := z.R().
		SetResult(&respJson).
		SetFormData(map[string]string{
			"topologyName":     tar.TopologyName,
			"chartName":        tar.ChartName,
			"chartDescription": tar.ChartDescription,
			"version":          tar.Version,
		}).
		SetFile("chart", p.V2FileOutPath()).
		Post(InfraCreateV1Path)

	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Ctx(ctx).Err(err).Msg("ZeusClient: UploadChart")
		return respJson, err
	}
	z.PrintRespJson(resp.Body())
	return respJson, err
}

func (z *ZeusClient) ZipK8sChartToPath(p *structs.Path) error {
	comp := compression.NewCompression()
	err := comp.CreateTarGzipArchiveDir(p)
	if err != nil {
		log.Err(err).Interface("path", p).Msg("ZeusClient: ZipK8sChartToPath")
		return err
	}
	return err
}

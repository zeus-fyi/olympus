package zeus_client

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
	read_topologies "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/endpoints"
)

func (z *ZeusClient) ReadTopologies(ctx context.Context) (read_topologies.ReadTopologiesMetadataGroup, error) {
	respJson := read_topologies.ReadTopologiesMetadataGroup{}
	resp, err := z.R().
		SetResult(&respJson.Slice).
		Get(zeus_endpoints.InfraReadTopologyV1Path)

	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Ctx(ctx).Err(err).Msg("ZeusClient: ReadTopologies")
		return respJson, err
	}
	z.PrintRespJson(resp.Body())
	return respJson, err
}

package athena_client

import (
	"context"
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	athena_endpoints "github.com/zeus-fyi/olympus/pkg/athena/client/endpoints"
)

func (a *AthenaClient) GetHostDiskInfo(ctx context.Context) (*disk.UsageStat, error) {
	respJson := &disk.UsageStat{}
	resp, err := a.R().
		SetResult(respJson).
		Get(athena_endpoints.InternalHostDiskV1Path)
	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Ctx(ctx).Err(err).Msg("AthenaClient: Kill")
		if resp.StatusCode() == http.StatusBadRequest {
			err = errors.New("bad request")
		}
		return respJson, err
	}
	a.PrintRespJson(resp.Body())
	return respJson, err
}

func (a *AthenaClient) GetHostMemInfo(ctx context.Context) (*mem.VirtualMemoryStat, error) {
	respJson := &mem.VirtualMemoryStat{}
	resp, err := a.R().
		SetResult(respJson).
		Get(athena_endpoints.InternalHostMemV1Path)
	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Ctx(ctx).Err(err).Msg("AthenaClient: GetHostMemInfo")
		if resp.StatusCode() == http.StatusBadRequest {
			err = errors.New("bad request")
		}
		return respJson, err
	}
	a.PrintRespJson(resp.Body())
	return respJson, err
}

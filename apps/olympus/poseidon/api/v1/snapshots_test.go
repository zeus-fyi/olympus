package v1_poseidon

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	pg_poseidon "github.com/zeus-fyi/olympus/datastores/postgres/apps/poseidon"
	poseidon_base_test "github.com/zeus-fyi/olympus/poseidon/api/test"
	zeus_client "github.com/zeus-fyi/zeus/pkg/zeus/client"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

type SnapshotDiskRequestTestSuite struct {
	poseidon_base_test.PoseidonBaseTestSuite
}

const (
	useProd               = true
	productionPoseidonURL = "https://poseidon.zeus.fyi/v1"
	localPoseidonURL      = "http://localhost:9010/v1"
	dwRoute               = "/ethereum/beacon/disk/wipe"
	ssUpload              = "/ethereum/beacon/disk/upload"
)

func (t *SnapshotDiskRequestTestSuite) TestDiskWipe() {
	t.InitLocalConfigs()
	t.Eg.POST(dwRoute, DiskWipeRequestHandler)

	start := make(chan struct{}, 1)
	go func() {
		close(start)
		_ = t.E.Start(":9010")
	}()

	<-start
	ctx := context.Background()
	defer t.E.Shutdown(ctx)

	dw := DiskWipeRequest{
		pg_poseidon.DiskWipeOrchestration{
			ClientName: "geth",
			OrchestrationJob: artemis_orchestrations.OrchestrationJob{
				Orchestrations: artemis_autogen_bases.Orchestrations{},
				Scheduled:      artemis_autogen_bases.OrchestrationsScheduledToCloudCtxNs{},
				CloudCtxNs: zeus_common_types.CloudCtxNs{
					CloudProvider: "do",
					Region:        "sfo3",
					Context:       "do-sfo3-dev-do-sfo3-zeus",
					Namespace:     "athena-beacon-goerli",
					Env:           "",
				},
			},
		},
	}

	err := SendDiskWipeRequest(ctx, t.ZeusClient, dw)
	t.Require().NoError(err)
}

func (t *SnapshotDiskRequestTestSuite) TestSnapshotUploadRequest() {
	t.InitLocalConfigs()
	t.Eg.POST(ssUpload, SnapshotUploadRequestHandler)

	start := make(chan struct{}, 1)
	go func() {
		close(start)
		_ = t.E.Start(":9010")
	}()

	<-start
	ctx := context.Background()
	defer t.E.Shutdown(ctx)

	su := SnapshotUploadRequest{
		pg_poseidon.UploadDataDirOrchestration{
			ClientName: "lighthouse",
			OrchestrationJob: artemis_orchestrations.OrchestrationJob{
				Orchestrations: artemis_autogen_bases.Orchestrations{},
				Scheduled:      artemis_autogen_bases.OrchestrationsScheduledToCloudCtxNs{},
				CloudCtxNs: zeus_common_types.CloudCtxNs{
					CloudProvider: "do",
					Region:        "sfo3",
					Context:       "do-sfo3-dev-do-sfo3-zeus",
					Namespace:     "goerli-staking",
					Env:           "",
				},
			},
		},
	}

	err := SendSnapshotUploadRequest(ctx, t.ZeusClient, su)
	t.Require().NoError(err)
}

func SendDiskWipeRequest(ctx context.Context, z zeus_client.ZeusClient, dw DiskWipeRequest) error {
	z.PrintReqJson(dw)

	url := localPoseidonURL + dwRoute

	if useProd {
		url = productionPoseidonURL + dwRoute
	}
	resp, err := z.R().
		SetBody(&dw).
		Post(url)

	if err != nil || (resp.StatusCode() != http.StatusAccepted && resp.StatusCode() != http.StatusOK) {
		log.Ctx(ctx).Err(err).Msg("ZeusClient: SendDiskWipeRequest")
		if resp.StatusCode() == http.StatusBadRequest {
			err = errors.New("bad request")
		}
		if err == nil {
			err = fmt.Errorf("non-OK status code: %d", resp.StatusCode())
		}
		return err
	}
	z.PrintRespJson(resp.Body())
	return err
}

func SendSnapshotUploadRequest(ctx context.Context, z zeus_client.ZeusClient, su SnapshotUploadRequest) error {
	z.PrintReqJson(su)

	url := localPoseidonURL + ssUpload

	if useProd {
		url = productionPoseidonURL + ssUpload
	}
	resp, err := z.R().
		SetBody(&su).
		Post(url)

	if err != nil || (resp.StatusCode() != http.StatusAccepted && resp.StatusCode() != http.StatusOK) {
		log.Ctx(ctx).Err(err).Msg("ZeusClient: SendSnapshotUploadRequest")
		if resp.StatusCode() == http.StatusBadRequest {
			err = errors.New("bad request")
		}
		if err == nil {
			err = fmt.Errorf("non-OK status code: %d", resp.StatusCode())
		}
		return err
	}
	z.PrintRespJson(resp.Body())
	return err
}

func TestSnapshotDiskRequestTestSuite(t *testing.T) {
	suite.Run(t, new(SnapshotDiskRequestTestSuite))
}

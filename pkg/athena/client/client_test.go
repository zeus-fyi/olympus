package athena_client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	athena_routines "github.com/zeus-fyi/olympus/athena/api/v1/common/routines"
	"github.com/zeus-fyi/olympus/pkg/poseidon"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"

	test_base "github.com/zeus-fyi/olympus/test"
	api_configs "github.com/zeus-fyi/olympus/test/configs"
)

var ctx context.Context

type AthenaClientTestSuite struct {
	base.TestSuite
	AthenaTestClient AthenaClient
}

func (t *AthenaClientTestSuite) SetupTest() {
	// points dir to test/configs
	tc := api_configs.InitLocalTestConfigs()
	//t.ZeusTestClient = NewDefaultZeusClient(tc.Bearer)
	t.AthenaTestClient = NewLocalAthenaClient(tc.Bearer)
	// points working dir to inside /test
	test_base.ForceDirToTestDirLocation()
}

func (t *AthenaClientTestSuite) DownloadTest() {
	br := poseidon.BucketRequest{
		BucketName: "",
		Protocol:   "",
		Network:    "",
		ClientType: "",
		ClientName: "",
	}
	err := t.AthenaTestClient.Download(ctx, br)
	t.Assert().Nil(err)
}

func (t *AthenaClientTestSuite) UploadTest() {
	br := poseidon.BucketRequest{
		BucketName: "",
		Protocol:   "",
		Network:    "",
		ClientType: "",
		ClientName: "",
	}
	err := t.AthenaTestClient.Upload(ctx, br)
	t.Assert().Nil(err)
}

func (t *AthenaClientTestSuite) ResumeTest() {
	rr := athena_routines.RoutineRequest{ClientName: ""}
	err := t.AthenaTestClient.Resume(ctx, rr)
	t.Assert().Nil(err)
}

func (t *AthenaClientTestSuite) PauseTest() {
	rr := athena_routines.RoutineRequest{ClientName: ""}
	err := t.AthenaTestClient.Pause(ctx, rr)
	t.Assert().Nil(err)
}

func TestAthenaClientTestSuite(t *testing.T) {
	suite.Run(t, new(AthenaClientTestSuite))
}

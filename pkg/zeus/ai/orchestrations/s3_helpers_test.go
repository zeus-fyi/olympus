package ai_platform_service_orchestrations

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/aegis/aws_secrets"
)

func (t *ZeusWorkerTestSuite) TestSecretsSelect() {
	user2 := 1710298581127603000

	// user2 FlowsOrgID
	ps, err := aws_secrets.GetMockingbirdPlatformSecrets(ctx, org_users.NewOrgUserWithID(user2, 0), "api-iris")
	t.Require().Nil(err)
	t.Require().NotNil(ps)
	t.Assert().NotEmpty(ps.S3AccessKey)
	t.Assert().NotEmpty(ps.S3SecretKey)
}

func (t *ZeusWorkerTestSuite) TestSaveWfStatus() {
	fnv := "CsvIteratorDebug-cycle-1-chunk-0-1712604784507256000.json"
	dbg := OpenCsvIteratorDebug(fnv)
	wfs := WfStatus{
		TotalApiRequests:    10,
		CompleteApiRequests: 9,
		TotalCsvElements:    100,
		CompleteCsvElements: 99,
	}
	err := saveWfStatus(ctx, dbg.Cp, wfs)
	t.Require().Nil(err)
	wfrs, err := getWfStatus(ctx, dbg.Cp)
	t.Require().Nil(err)
	t.Require().NotNil(wfrs)
}

package athena_client

import (
	athena_routines "github.com/zeus-fyi/olympus/athena/api/v1/common/routines"
)

var clientName = "lighthouse"

func (t *AthenaClientTestSuite) TestResume() {
	rr := athena_routines.RoutineRequest{ClientName: clientName}
	err := t.AthenaTestClient.Resume(ctx, rr)
	t.Assert().Nil(err)
}

func (t *AthenaClientTestSuite) TestSuspend() {
	rr := athena_routines.RoutineRequest{ClientName: clientName}
	err := t.AthenaTestClient.Suspend(ctx, rr)
	t.Assert().Nil(err)
}

func (t *AthenaClientTestSuite) TestKill() {
	rr := athena_routines.RoutineRequest{ClientName: clientName}
	err := t.AthenaTestClient.Kill(ctx, rr)
	t.Assert().Nil(err)
}

func (t *AthenaClientTestSuite) TestDiskWipe() {
	rr := athena_routines.RoutineRequest{ClientName: clientName}
	err := t.AthenaTestClient.DiskWipe(ctx, rr)
	t.Assert().Nil(err)
}

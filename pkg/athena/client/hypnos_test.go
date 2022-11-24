package athena_client

import athena_routines "github.com/zeus-fyi/olympus/athena/api/v1/common/routines"

func (t *AthenaClientTestSuite) TestHypnos() {
	rr := athena_routines.RoutineRequest{ClientName: clientName}
	err := t.AthenaTestClient.Hypnos(ctx, rr)
	t.Assert().Nil(err)
}

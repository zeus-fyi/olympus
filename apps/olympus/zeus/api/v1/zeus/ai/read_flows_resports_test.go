package zeus_v1_ai

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func (t *FlowsWorkerTestSuite) TestFlowExport() {
	tmpOu := t.Ou
	tmpOu.OrgID = 1685378241971196000

	_, err := ExportRunCsvRequest2(ctx, tmpOu, 1714378055182320000)
	t.Require().Nil(err)
}

func TestFlowExportTestSuite(t *testing.T) {
	suite.Run(t, new(FlowsWorkerTestSuite))
}

package zeus_demo

import (
	"github.com/zeus-fyi/olympus/pkg/ares/demo"
)

func (t *AresDemoTestSuite) TestReadTopologies() {
	demo.ChangeDirToAresDemoDir()
	resp, err := t.ZeusTestClient.ReadTopologies(ctx)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}

package zeus_demo

import (
	"github.com/zeus-fyi/olympus/pkg/ares/demo"
)

func (t *AresDemoTestSuite) TestDeploy() {
	demo.ChangeDirToAresDemoDir()
	resp, err := t.ZeusTestClient.Deploy(ctx, deployKnsReq)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}

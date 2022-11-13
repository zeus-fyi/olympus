package zeus_demo

import (
	"github.com/zeus-fyi/olympus/pkg/ares/demo"
)

func (t *AresDemoTestSuite) TestDestroyDeploy() {
	demo.ChangeDirToAresDemoDir()
	resp, err := t.ZeusTestClient.DestroyDeploy(ctx, deployDestroyKnsReq)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}

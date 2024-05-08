package zeus_v1_ai

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	resty_base "github.com/zeus-fyi/zeus/zeus/z_client/base"
)

type FlowsWorkerTestSuite struct {
	test_suites_base.TestSuite
}

var ctx = context.Background()

func (t *FlowsWorkerTestSuite) TestGetMappedColumns() {

	cm := map[string]string{
		linkedIn: "colname",
	}
	w := ExecFlowsActionsRequest{}
	w.StageContactsMap = cm

}

func (t *FlowsWorkerTestSuite) SetupTest() {
	t.InitLocalConfigs()
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
}
func hitAccount() {

	i := 0
	for {
		rcb := resty.New()
		r := resty_base.Resty{
			Client:    rcb,
			PrintReq:  false,
			PrintResp: false,
		}
		i++
		ul := "https://nat.org/"
		re, err := r.R().Get(ul)
		if err != nil {
			log.Err(err).Interface("re", re).Msg("hitAccount")
		}
		fmt.Println(i, "code", re.StatusCode())
		time.Sleep(time.Millisecond)
	}
}
func (t *FlowsWorkerTestSuite) TestFlowMultiPrompt1() {
	hitAccount()
}

func (t *FlowsWorkerTestSuite) TestFlowMultiPrompt() {
	tmpOu := t.Ou
	tmpOu.OrgID = 1685378241971196000
	wfs := []artemis_orchestrations.WorkflowTemplate{
		{
			WorkflowName: webFetchWf,
		},
		{
			WorkflowName: webFetchWf,
		},
	}
	wfc := make(map[string]int)

	for _, wfv := range wfs {
		wfc[wfv.WorkflowName] += 1
	}

	resp, err := artemis_orchestrations.GetAiOrchestrationParams(ctx, tmpOu, nil, wfs)
	t.Require().Nil(err)
	t.Require().NotEmpty(resp)

	var wfts []artemis_orchestrations.WorkflowTemplateData
	for _, v := range resp {
		for _, tv := range v.WorkflowTasks {
			if tv.AnalysisTaskName != "" {
				wfts = append(wfts, tv)
				fmt.Println(tv.AnalysisTaskName)
			}
		}
	}

	t.Assert().Equal(wfc[webFetchWf], len(wfts))
}

func TestFlowsWorkerTestSuite(t *testing.T) {
	suite.Run(t, new(FlowsWorkerTestSuite))
}

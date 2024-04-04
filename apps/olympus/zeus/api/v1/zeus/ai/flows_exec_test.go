package zeus_v1_ai

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type FlowsWorkerTestSuite struct {
	test_suites_base.TestSuite
}

var ctx = context.Background()

func (t *FlowsWorkerTestSuite) SetupTest() {
	t.InitLocalConfigs()
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
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

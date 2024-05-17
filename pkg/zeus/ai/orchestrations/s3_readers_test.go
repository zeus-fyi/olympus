package ai_platform_service_orchestrations

import (
	"encoding/json"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
)

/*

 */

//	func (t *ZeusWorkerTestSuite) Testgs3wfsCustomTaskName() {
//		cp := &MbChildSubProcessParams{
//
//		}
//		b, berr := gs3wfsCustomTaskName(ctx, cp, fmt.Sprintf("%d", r.WorkflowResultID))
//
//		if berr != nil {
//			log.Err(berr).Msg("AiAggregateAnalysisRetrievalTask: failed")
//			continue
//		}
//	}

func (t *ZeusWorkerTestSuite) TestS3WfDebugRunExport() {
	wfn := "csv-analysis-dbc1318e-c65c"
	b, err := S3WfDebugRunExport(ctx, wfn)
	t.Require().Nil(err)
	t.Require().NotNil(b.Bytes())

	mb := MbChildSubProcessParams{}
	err = json.Unmarshal(b.Bytes(), &mb)
	t.Require().Nil(err)
}

func (t *ZeusWorkerTestSuite) TestGetGlobalEntitiesFromRef() {
	ueh := "b4d0c637a8768434cc90142d15c76ea1959ce3cfaba037fafad7232d0c9415fab4d0c637a8768434cc90142d15c76ea1959ce3cfaba037fafad7232d0c9415fa"
	refs := []artemis_entities.EntitiesFilter{
		{
			Nickname: ueh,
			Platform: "flows",
		},
	}
	ue, err := GetGlobalEntitiesFromRef(ctx, t.Ou, refs)
	t.Require().Nil(err)
	t.Require().NotEmpty(ue)
	for _, v := range ue {
		t.Assert().NotZero(len(v.MdSlice))
		for _, mv := range v.MdSlice {
			t.Assert().NotNil(mv.TextData)
			t.Assert().NotEmpty(mv.Labels)
		}
	}
}

func (t *ZeusWorkerTestSuite) TestS3WfCycleStageRead() {
	wsi := t.mockCsvMerge()
	cp := &MbChildSubProcessParams{
		WfExecParams: wsi.WorkflowExecParams,
		Ou:           t.Ou,
		Tc: TaskContext{
			TaskName: "validate-emails",
			TaskType: "analysis",
		},
	}
	res, err := gs3wfs(ctx, cp)
	t.Require().Nil(err)
	t.Require().NotNil(res)

	t.Require().NotNil(res.PromptReduction)
	t.Require().NotNil(res.PromptReduction.PromptReductionSearchResults)
	t.Require().NotEmpty(res.PromptReduction.PromptReductionSearchResults.OutSearchGroups)
	for _, v := range res.PromptReduction.PromptReductionSearchResults.OutSearchGroups {
		t.Assert().Len(v.ApiResponseResults, 2)
	}
}

func (t *ZeusWorkerTestSuite) TestS3HelperDownloadWfData() {
	cp := &MbChildSubProcessParams{}
	wfi, err := gs3wfs(ctx, cp)
	t.Require().Nil(err)
	t.Require().NotNil(wfi)
}

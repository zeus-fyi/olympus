package ai_platform_service_orchestrations

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
)

const (
	validemailRetQp = "validemail-query-params"
)

func (t *ZeusWorkerTestSuite) TestWfCsv() {
	ueh := "b4d0c637a8768434cc90142d15c76ea1959ce3cfaba037fafad7232d0c9415fab4d0c637a8768434cc90142d15c76ea1959ce3cfaba037fafad7232d0c9415fa"
	csvSourceEntity, csvContacts := t.getContactCsvMock()
	_, emRow := ts(csvContacts)
	t.Require().NotEmpty(emRow)
	b, err := json.Marshal(emRow)
	t.Require().Nil(err)
	csvMergeInEntity := artemis_entities.UserEntity{
		Nickname: ueh,
		Platform: "flows",
		MdSlice: []artemis_entities.UserEntityMetadata{
			{
				JsonData: b,
				TextData: aws.String("Email"),
				Labels:   artemis_entities.CreateMdLabels([]string{csvGlobalMergeRetLabel(validemailRetQp)}),
			},
		},
	}
	expa := artemis_orchestrations.WorkflowExecParams{
		WorkflowOverrides: artemis_orchestrations.WorkflowOverrides{
			WorkflowRunName: "test-wf",
			IsUsingFlows:    true,
			WorkflowEntityRefs: []artemis_entities.EntitiesFilter{
				{
					Nickname: ueh,
					Platform: "flows",
					Labels:   []string{csvSrcGlobalLabel, csvGlobalMergeRetLabel(validemailRetQp)},
				},
			},
			WorkflowEntities: []artemis_entities.UserEntity{
				csvMergeInEntity,
			},
		},
	}
	wsi := &WorkflowStageIO{
		WorkflowExecParams: expa,
		WorkflowStageReference: artemis_orchestrations.WorkflowStageReference{
			RunCycle: 0,
		},
		WorkflowStageInfo: WorkflowStageInfo{
			PromptReduction: &PromptReduction{
				PromptReductionSearchResults: &PromptReductionSearchResults{
					OutSearchGroups: []*hera_search.SearchResultGroup{
						{
							RetrievalName: aws.String(validemailRetQp),
							ApiResponseResults: []hera_search.SearchResult{
								{
									WebResponse: hera_search.WebResponse{
										Body: echo.Map{
											"Email": "alex@zeus.fyi",
											"Tag":   false,
											"Score": 25,
										},
									},
								},
								{
									WebResponse: hera_search.WebResponse{
										Body: echo.Map{
											"Email": "leevar@gmail.com",
											"Tag":   true,
											"Score": 90,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	ur, err := FindAndMergeMatchingNicknamesByLabel(
		csvSourceEntity,
		[]artemis_entities.UserEntity{csvMergeInEntity},
		wsi,
		csvGlobalMergeRetLabel(validemailRetQp),
	)
	t.Require().Nil(err)
	t.Require().NotEmpty(ur)
	t.Assert().NotEmpty(ur.MdSlice)

	//cp := &MbChildSubProcessParams{
	//	WfExecParams: expa,
	//	Ou:           t.Ou,
	//	Tc: TaskContext{
	//		TaskName: "validate-emails",
	//		TaskType: "analysis",
	//	},
	//}
	//za := NewZeusAiPlatformActivities()
	//wr := &artemis_orchestrations.AIWorkflowAnalysisResult{}
	//res, err := za.SaveCsvTaskOutput(ctx, cp, wr)
	//t.Require().Nil(err)
	//t.Assert().NotEmpty(res)
}

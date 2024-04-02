package ai_platform_service_orchestrations

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
)

func (t *ZeusWorkerTestSuite) TestGetGlobalEntitiesFromRef() {
	ueh := "48123fb7c2365017fec6634c0650dd0ed07c7986796114f9d6d1154f4aaf9acb48123fb7c2365017fec6634c0650dd0ed07c7986796114f9d6d1154f4aaf9acb"
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
	}
}

func (t *ZeusWorkerTestSuite) TestWsiOut() {
	wsi := WorkflowStageIO{
		WorkflowExecParams:     artemis_orchestrations.WorkflowExecParams{},
		WorkflowStageReference: artemis_orchestrations.WorkflowStageReference{},
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
	m := map[string]bool{
		validemailRetQp: true,
	}
	sgs := wsi.GetSearchGroupsOutByRetNameMatch(m)
	t.Require().NotEmpty(sgs)

}

func (t *ZeusWorkerTestSuite) TestPayloadToCsvString() {
	constactsCsvStr := "First Name,Last Name,Company,LinkedIn,Email,Website\nAlex,George,Zeusfyi,https://www.linkedin.com/in/alexandersgeorge/,alex@zeus.fyi,http://www.bsci.com\nLevar,Williams,APrime Technology,https://www.linkedin.com/in/leevarwilliams/,leevar@gmail.com,http://www.natroxwoundcare.com\nAlex,George,Zeusfyi,https://www.linkedin.com/in/alexandersgeorge/,alex@zeus.fyi,http://www.shockwavemedical.com\nLevar,Williams,APrime Technology,https://www.linkedin.com/in/leevarwilliams/,leevar@gmail.com,http://www.ottobock.com\n"
	csvContacts, err := ParseCsvStringToMap(constactsCsvStr)
	t.Require().Nil(err)
	t.Require().NotEmpty(csvContacts)
	cs, err := PayloadToCsvString(csvContacts)
	t.Require().Nil(err)
	fmt.Println(cs)
	//t.Assert().Equal(constactsCsvStr, cs)
}

const (
	validemailRetQp = "validemail-query-params"
)

func (t *ZeusWorkerTestSuite) TestCsvMerge() {
	constactsCsvStr := "First Name,Last Name,Company,LinkedIn,Email,Website\nAlex,George,Zeusfyi,https://www.linkedin.com/in/alexandersgeorge/,alex@zeus.fyi,http://www.bsci.com\nLevar,Williams,APrime Technology,https://www.linkedin.com/in/leevarwilliams/,leevar@gmail.com,http://www.natroxwoundcare.com\nAlex,George,Zeusfyi,https://www.linkedin.com/in/alexandersgeorge/,alex@zeus.fyi,http://www.shockwavemedical.com\nLevar,Williams,APrime Technology,https://www.linkedin.com/in/leevarwilliams/,leevar@gmail.com,http://www.ottobock.com\n"
	csvContacts, err := ParseCsvStringToMap(constactsCsvStr)
	t.Require().Nil(err)
	t.Require().NotEmpty(csvContacts)
	//emRow := map[string][]int{
	//	"alex@zeus.fyi":    []int{0, 2},
	//	"leevar@gmail.com": []int{1, 3},
	//}

	colName, emRow := ts(csvContacts)
	t.Require().NotEmpty(emRow)
	fmt.Println(emRow)

	csvStr := "Tag,Free,Role,Email,Score,State,Domain,Reason,IsValid,MXRecord,AcceptAll,Disposable,EmailAdditionalInfo\n\"false\",\"false\",\"false\",\"alex@zeus.fyi\",\"60\",\"Deliverable\",\"zeus.fyi\",\"ACCEPTED EMAIL\",\"true\",\"aspmx.l.google.com.\",\"true\",\"false\",\"\"\n\"false\",\"true\",\"false\",\"leevar@gmail.com\",\"95\",\"Deliverable\",\"gmail.com\",\"ACCEPTED EMAIL\",\"true\",\"gmail-smtp-in.l.google.com.\",\"false\",\"false\",\"\""
	csvData, err := ParseCsvStringToMap(csvStr)
	t.Require().Nil(err)
	t.Require().NotEmpty(csvData)
	mergedCsv, err := appendCsvData(csvContacts, csvData, colName, emRow)
	t.Require().Nil(err)
	fmt.Println(mergedCsv)
}
func (t *ZeusWorkerTestSuite) TestWfCsv() {
	constactsCsvStr := "First Name,Last Name,Company,LinkedIn,Email,Website\nAlex,George,Zeusfyi,https://www.linkedin.com/in/alexandersgeorge/,alex@zeus.fyi,http://www.bsci.com\nLevar,Williams,APrime Technology,https://www.linkedin.com/in/leevarwilliams/,leevar@gmail.com,http://www.natroxwoundcare.com\nAlex,George,Zeusfyi,https://www.linkedin.com/in/alexandersgeorge/,alex@zeus.fyi,http://www.shockwavemedical.com\nLevar,Williams,APrime Technology,https://www.linkedin.com/in/leevarwilliams/,leevar@gmail.com,http://www.ottobock.com\n"
	csvSourceEntity := artemis_entities.UserEntity{
		Nickname: "test",
		Platform: "flows",
		MdSlice: []artemis_entities.UserEntityMetadata{
			{
				TextData: aws.String(constactsCsvStr),
				Labels:   artemis_entities.CreateMdLabels([]string{"csv:source", fmt.Sprintf("csv:merge:ret:%s", validemailRetQp)}),
			},
		},
	}
	csvContacts, err := ParseCsvStringToMap(constactsCsvStr)
	t.Require().Nil(err)
	t.Require().NotEmpty(csvContacts)
	_, emRow := ts(csvContacts)
	b, err := json.Marshal(emRow)
	t.Require().Nil(err)
	csvMergeInEntity := artemis_entities.UserEntity{
		Nickname: "test",
		Platform: "flows",
		MdSlice: []artemis_entities.UserEntityMetadata{
			{
				JsonData: b,
				TextData: aws.String("Email"),
				Labels:   artemis_entities.CreateMdLabels([]string{"csv:merge", fmt.Sprintf("csv:merge:ret:%s", validemailRetQp)}),
			},
		},
	}
	wsi := &WorkflowStageIO{
		WorkflowExecParams:     artemis_orchestrations.WorkflowExecParams{},
		WorkflowStageReference: artemis_orchestrations.WorkflowStageReference{},
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

	ur, err := FindAndMergeMatchingNicknamesByLabel(csvSourceEntity, []artemis_entities.UserEntity{csvMergeInEntity}, wsi, fmt.Sprintf("csv:merge:ret:%s", validemailRetQp))
	t.Require().Nil(err)
	t.Require().NotEmpty(ur)
	t.Assert().NotEmpty(ur.MdSlice)
}

func ts(csvContacts []map[string]string) (string, map[string][]int) {
	seen := make(map[string]bool)
	var pls []map[string]interface{}
	emRow := make(map[string][]int)
	var colName string
	for r, cv := range csvContacts {
		for cn, emv := range cv {
			tv := strings.ToLower(cn)
			if strings.Contains(tv, "email") && len(emv) > 0 && strings.Contains(emv, "@") {
				if len(colName) > 0 && colName != cn {
					panic("duplicate col")
				}
				colName = cn
				etm := emRow[emv]
				etm = append(etm, r)
				emRow[emv] = etm
				if _, ok := seen[emv]; ok {
					continue
				}
				pl := make(map[string]interface{})
				pl["email"] = emv
				pls = append(pls, pl)
			}
			seen[emv] = true
		}
	}
	return colName, emRow
}

package ai_platform_service_orchestrations

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	utils_csv "github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/csv"
)

// mergeCsvs

func (t *ZeusWorkerTestSuite) TestCoreCsvMerge() {
	srcCsvEntity := artemis_entities.UserEntity{}
	mergeInCsvEntities := []artemis_entities.UserEntity{}
	wsi := &WorkflowStageIO{}
	_, err := mergeCsvs(srcCsvEntity, mergeInCsvEntities, wsi)
	t.Require().Nil(err)
}

func (t *ZeusWorkerTestSuite) TestAppendCsv() {
	t.mockCsvMerge()
	//var inputCsvData, csvData []map[string]string
	//mergedCsv, err := appendCsvData(inputCsvData, csvData, "", nil)
	//t.Require().Nil(err)
	//t.Require().NotEmpty(mergedCsv)
}

/*
	1. create entity source: with csv str
	2. merged entity to have:
		a. textField| colName
		b. jsonField| emRow
		c. chained w/  csvRowsAppend []map[string]interface{}
*/

/*
func appendCsvData(inputCsvData, csvData []map[string]string, colName string, emRow map[string][]int) ([]map[string]string, error) {
	// Iterate through csvData to find and merge matching rows
	for _, dataRow := range csvData {
		email := dataRow[colName]
		if indices, ok := emRow[email]; ok {
			// If a matching row is found, merge the data
			for _, index := range indices {
				for key, value := range dataRow {
					inputCsvData[index][key] = value
				}
			}
		}
	}
	return inputCsvData, nil
}
*/

func (t *ZeusWorkerTestSuite) mockCsvMerge() *WorkflowStageIO {
	ueh := "b4d0c637a8768434cc90142d15c76ea1959ce3cfaba037fafad7232d0c9415fab4d0c637a8768434cc90142d15c76ea1959ce3cfaba037fafad7232d0c9415fa"
	csvSourceEntity, csvContacts := t.getContactCsvMock()
	t.Require().NotEmpty(csvSourceEntity)
	cn, emRow := ts(csvContacts)

	b, err := utils_csv.NewCsvMergeEntityFromSrcBin(cn, emRow)
	t.Require().Nil(err)
	csvMergeInEntity := artemis_entities.UserEntity{
		Nickname: ueh,
		Platform: "flows",
		MdSlice: []artemis_entities.UserEntityMetadata{
			{
				JsonData: b,
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
	// to debug
	ur, err := FindAndMergeMatchingNicknamesByLabelPrefix(
		csvSourceEntity,
		[]artemis_entities.UserEntity{csvMergeInEntity},
		wsi,
		csvGlobalMergeRetLabel(validemailRetQp),
	)
	t.Require().Nil(err)
	t.Require().NotEmpty(ur)
	t.Assert().NotEmpty(ur.MdSlice)
	for _, v := range ur.MdSlice {
		t.Require().NotNil(v.TextData)
		fmt.Println(*v.TextData)
	}
	return wsi
}

func (t *ZeusWorkerTestSuite) TestCsvMerge() {
	constactsCsvStr := "First Name,Last Name,Company,LinkedIn,Email,Website\nAlex,George,Zeusfyi,https://www.linkedin.com/in/alexandersgeorge/,alex@zeus.fyi,http://www.bsci.com\nLevar,Williams,APrime Technology,https://www.linkedin.com/in/leevarwilliams/,leevar@gmail.com,http://www.natroxwoundcare.com\nAlex,George,Zeusfyi,https://www.linkedin.com/in/alexandersgeorge/,alex@zeus.fyi,http://www.shockwavemedical.com\nLevar,Williams,APrime Technology,https://www.linkedin.com/in/leevarwilliams/,leevar@gmail.com,http://www.ottobock.com\n"
	csvContacts, err := utils_csv.ParseCsvStringToMap(constactsCsvStr)
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
	csvData, err := utils_csv.ParseCsvStringToMap(csvStr)
	t.Require().Nil(err)
	t.Require().NotEmpty(csvData)
	mergedCsv, err := utils_csv.AppendCsvData(csvContacts, csvData, colName, emRow)
	t.Require().Nil(err)
	fmt.Println(mergedCsv)
}

func (t *ZeusWorkerTestSuite) getContactCsvMock() (artemis_entities.UserEntity, []map[string]string) {
	constactsCsvStr := "First Name,Last Name,Company,LinkedIn,Email,Website\nAlex,George,Zeusfyi,https://www.linkedin.com/in/alexandersgeorge/,alex@zeus.fyi,http://www.bsci.com\nLevar,Williams,APrime Technology,https://www.linkedin.com/in/leevarwilliams/,leevar@gmail.com,http://www.natroxwoundcare.com\nAlex,George,Zeusfyi,https://www.linkedin.com/in/alexandersgeorge/,alex@zeus.fyi,http://www.shockwavemedical.com\nLevar,Williams,APrime Technology,https://www.linkedin.com/in/leevarwilliams/,leevar@gmail.com,http://www.ottobock.com\n"
	csvSourceEntity := artemis_entities.UserEntity{
		Nickname: "b4d0c637a8768434cc90142d15c76ea1959ce3cfaba037fafad7232d0c9415fab4d0c637a8768434cc90142d15c76ea1959ce3cfaba037fafad7232d0c9415fa",
		Platform: "flows",
		MdSlice: []artemis_entities.UserEntityMetadata{
			{
				TextData: aws.String(constactsCsvStr),
				Labels:   artemis_entities.CreateMdLabels([]string{csvSrcGlobalLabel, csvGlobalMergeRetLabel(validemailRetQp)}),
			},
		},
	}
	csvContacts, err := utils_csv.ParseCsvStringToMap(constactsCsvStr)
	t.Require().Nil(err)
	t.Require().NotEmpty(csvContacts)
	return csvSourceEntity, csvContacts
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
	csvContacts, err := utils_csv.ParseCsvStringToMap(constactsCsvStr)
	t.Require().Nil(err)
	t.Require().NotEmpty(csvContacts)
	cs, err := utils_csv.PayloadToCsvString(csvContacts)
	t.Require().Nil(err)
	fmt.Println(cs)
	//t.Assert().Equal(constactsCsvStr, cs)
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

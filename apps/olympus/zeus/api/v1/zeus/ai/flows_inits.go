package zeus_v1_ai

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	utils_csv "github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/csv"
	ai_platform_service_orchestrations "github.com/zeus-fyi/olympus/pkg/zeus/ai/orchestrations"
)

const (
	// wfs
	// 	prev webFetchWf = "website-analysis-wf"
	// new 	webFetchWf = "scraped-website-analysis-wf"
	// 	webFetchWf = "scraped-website-analysis-wf"
	zenrowsWebFetchWf = "zenrows-website-analysis-wf"
	webFetchWf        = zenrowsWebFetchWf
	//webFetchWf        = "scraped-website-analysis-wf"
	emailVdWf = "validate-emails-wf"
	googWf    = "google-query-regex-index-wf"

	googCsvWf = "google-search-csv-wf"

	liWf    = "linkedin-rapid-api-profiles-wf"
	liBizWf = "linkedin-rapid-api-biz-profiles-wf"

	// wf identifier stage
	linkedIn       = "linkedIn"
	linkedInBiz    = "linkedInBiz"
	googleSearch   = "googleSearch"
	validateEmails = "validateEmails"
	websiteScrape  = "websiteScrape"

	// ret override
	validemailRetQp  = "validemail-query-params"
	linkedInRetQp    = "linkedin-profile"
	linkedInBizRetQp = "linkedin-biz-profile"

	zenrowsScrapedWbsRet = "zenrows-website-analysis"
	//scrapedWbsRet        = "scraped-website-analysis"
	scrapedWbsRet = zenrowsScrapedWbsRet
	googleRetName = "google-query-params"
	// task
	wbsTaskName    = "website-analysis"
	googleTaskName = "google-search-analysis"

	// work labels
	csvSrcGlobalLabel      = "csv:global:source"
	csvSrcGlobalMergeLabel = "csv:global:merge"
)

// note this requires being first to process, others add to workflow slice on add

func (w *ExecFlowsActionsRequest) CustomCsvWorkflows(uef *artemis_entities.EntitiesFilter) error {
	if len(w.Workflows) <= 0 {
		return nil
	}
	for _, wfv := range w.Workflows {
		var colName string
		//seen := make(map[string]bool)
		emRow := make(map[string][]int)
		var pls []map[string]interface{}
		for r, cvs := range w.ContactsCsv {
			for cname, colValue := range cvs {
				// "company",
				if v, ok := w.StageContactsMap[cname]; ok && v == wfv.WorkflowName {
					//fmt.Println(r, v, colName, colValue, seen)
					tnm := cname
					if newVarName, ok1 := w.StageContactsOverrideMap[cname]; ok1 && len(newVarName) > 0 {
						tnm = newVarName
					}
					if len(emRow) <= 0 {
						emRow[colValue] = []int{r}
					} else {
						etm := emRow[colValue]
						etm = append(etm, r)
						emRow[colValue] = etm
					}
					pl := make(map[string]interface{})
					w.ContactsCsv[r][cname] = colValue
					// only for payload not override csv value
					pl[tnm] = w.ContactsCsv[r][cname]
					pls = append(pls, pl)
					//seen[v] = true
				}
			}
		}
		if len(pls) == 0 {
			log.Warn().Msg("no profiles found")
			return nil
		}
		w.InitMaps()
		for _, tv := range wfv.Tasks {
			if tv.ResponseFormat == "csv" && len(tv.TaskName) > 0 && len(tv.RetrievalName) > 0 {
				err := w.createCsvMergeEntity2(wfv.WorkflowName, tv.TaskName, tv.RetrievalName, uef, colName, emRow, pls)
				if err != nil {
					log.Err(err).Msg("createCsvMergeEntity: failed to marshal")
					return err
				}
				if v, ok := w.CommandPrompts[wfv.WorkflowName]; ok || (len(tv.TaskName) > 0 && len(tv.RetrievalName) > 0) {
					prompts := w.getPrompts()
					if len(prompts) <= 0 {
						if v == "" {
							v = "Can you tell me their role and responsibilities?"
						}
						prompts = []string{v}
					} else {
						tmp := w.TaskOverrides[tv.TaskName]
						tmp.SystemPromptExt = v
						w.TaskOverrides[tv.TaskName] = tmp
					}
					w.createWfTaskPromptOverrides(wfv.WorkflowName, tv.TaskName, w.getPromptsMap(wfv.WorkflowName))
					//w.createWfSchemaFieldOverride(liWf, wbsTaskName, "summary", prompts)
				}
			}
		}
	}
	return nil
}

func (w *ExecFlowsActionsRequest) SetupFlow(ctx context.Context, ou org_users.OrgUser) (*artemis_entities.EntitiesFilter, error) {
	uef := &artemis_entities.EntitiesFilter{
		Platform: "flows",
		Labels:   []string{csvSrcGlobalLabel},
	}

	w.InitMaps()
	headersCsv, err := w.ConvertToCsvStrToMap()
	if err != nil {
		log.Err(err).Interface("w", w).Msg("ConvertToCsvStrToMap failed")
		return nil, err
	}
	err = w.CustomCsvWorkflows(uef)
	if err != nil {
		log.Err(err).Interface("w", w).Msg("CustomCsvWorkflows failed")
		return nil, err
	}
	// this should add the email label
	err = w.EmailsValidatorSetup(uef)
	if err != nil {
		log.Err(err).Interface("w", w).Msg("EmailsValidatorSetup failed")
		return nil, err
	}
	// this should add the goog label, etc
	err = w.GoogleSearchSetup(uef)
	if err != nil {
		log.Err(err).Interface("w", w).Msg("GoogleSearchSetup failed")
		return nil, err
	}
	err = w.LinkedInScraperSetup(uef)
	if err != nil {
		log.Err(err).Interface("w", w).Msg("LinkedInScraperSetup failed")
		return nil, err
	}
	err = w.LinkedInBizScraperSetup(uef)
	if err != nil {
		log.Err(err).Interface("w", w).Msg("LinkedInBizScraperSetup failed")
		return nil, err
	}
	err = w.ScrapeRegularWebsiteSetup(uef)
	if err != nil {
		log.Err(err).Interface("w", w).Msg("ScrapeRegularWebsiteSetup failed")
		return nil, err
	}
	uef, err = w.SaveCsvImports(ctx, ou, uef, headersCsv)
	if err != nil {
		log.Err(err).Interface("w", w).Msg("EmailsValidatorSetup failed")
		return nil, err
	}
	if uef != nil && uef.Nickname != "" && uef.Platform != "" {
		w.WorkflowEntityRefs = append(w.WorkflowEntityRefs, *uef)
	}
	return uef, err
}

func (w *ExecFlowsActionsRequest) SaveCsvImports(ctx context.Context, ou org_users.OrgUser, uef *artemis_entities.EntitiesFilter, he []string) (*artemis_entities.EntitiesFilter, error) {
	if uef == nil || len(w.FlowsActionsRequest.ContactsCsvStr) == 0 {
		return nil, nil
	}
	cs, err := utils_csv.PayloadToCsvString(w.ContactsCsv)
	if err != nil {
		log.Err(err).Interface("w.ContactsCsv", w.ContactsCsvStr).Msg("SaveCsvImports: PayloadToCsvString: failed")
		return nil, err
	}
	cs, err = utils_csv.SortCSV(cs, he)
	if err != nil {
		log.Err(err).Interface("w.ContactsCsv", w.ContactsCsvStr).Msg("SaveCsvImports: PayloadToCsvString: failed")
		return nil, err
	}
	w.ContactsCsvStr = cs
	usre := &artemis_entities.UserEntity{
		Platform: uef.Platform,
		MdSlice: []artemis_entities.UserEntityMetadata{
			{
				TextData: aws.String(w.ContactsCsvStr),
				Labels:   artemis_entities.CreateMdLabels(uef.Labels),
			},
		},
	}
	nn, err := artemis_entities.HashParams(ou.OrgID, []interface{}{usre.MdSlice})
	if err != nil {
		log.Err(err).Msg("workingRunCycleStagePath: failed to hash wsr io")
		return nil, err
	}
	uef.Nickname = nn
	usre.Nickname = uef.Nickname
	if len(usre.Nickname) <= 0 || len(usre.Platform) <= 0 || len(usre.MdSlice) <= 0 {
		return nil, fmt.Errorf("no entities name")
	}
	for k, ue := range w.WorkflowEntities {
		ue.Nickname = nn
		w.WorkflowEntities[k] = ue
	}
	_, err = ai_platform_service_orchestrations.S3GlobalOrgUpload(ctx, ou, usre)
	if err != nil {
		log.Err(err).Msg("SaveImport: error")
		return nil, err
	}
	log.Info().Interface("ue", usre.Nickname).Msg("entity hash")
	return uef, nil
}

//func (w *ExecFlowsActionsRequest) AddPromptInject(uef *artemis_entities.EntitiesFilter, wfName string) error {
//	if w.PromptsCsvStr == "" {
//		return fmt.Errorf("no prompt inject csv provided")
//	}
//
//	switch wfName {
//	case webFetchWf, liWf, liBizWf, googWf:
//		ue := artemis_entities.UserEntity{
//			Nickname: uef.Nickname,
//			Platform: "prompts",
//			MdSlice: []artemis_entities.UserEntityMetadata{
//				{
//					TextData: aws.String(w.PromptsCsvStr),
//					Labels:   artemis_entities.CreateMdLabels([]string{"csv:prompts:wf"}),
//				},
//			},
//		}
//		w.WorkflowEntitiesOverrides[wfName] = append(w.WorkflowEntitiesOverrides[wfName], ue)
//	}
//	return nil
//}

func (w *ExecFlowsActionsRequest) ScrapeRegularWebsiteSetup(uef *artemis_entities.EntitiesFilter) error {
	if v, ok := w.Stages[websiteScrape]; !ok || !v {
		return nil
	}
	var colName string
	seen := make(map[string]bool)
	emRow := make(map[string][]int)
	var pls []map[string]interface{}
	for r, cvs := range w.ContactsCsv {
		for cname, colValue := range cvs {
			v, ok1 := w.StageContactsMap[cname]
			if (strings.Contains(cname, "web") || strings.Contains(cname, "url") || strings.Contains(cname, "link") || strings.Contains(cname, "site") || (ok1 && v == websiteScrape)) && len(colValue) > 0 {
				if len(colName) > 0 && colName != cname {
					log.Info().Interface("colName", colName).Interface("cname", cname).Msg("ScrapeRegularWebsiteSetup")
					return fmt.Errorf("duplicate web column")
				}
				colName = cname
				uv, err := convertToHTTPS(colValue)
				if err != nil {
					log.Err(err).Msg("failed to convert url to https")
					continue
				}
				if len(emRow) <= 0 {
					emRow[uv] = []int{r}
				} else {
					etm := emRow[uv]
					etm = append(etm, r)
					emRow[uv] = etm
				}
				if _, ok := seen[uv]; ok {
					w.ContactsCsv[r][cname] = uv
					continue
				}
				// skips linkedin
				if strings.HasPrefix(uv, "https://www.linkedin.com") || strings.HasPrefix(uv, "https://linkedin.com") {
					continue
				}
				pl := make(map[string]interface{})
				w.ContactsCsv[r][cname] = uv
				pl["url"] = w.ContactsCsv[r][cname]
				pls = append(pls, pl)
				seen[uv] = true
			}
		}
	}
	if len(pls) == 0 {
		log.Warn().Msg("no urls found")
		return nil
	}
	w.InitMaps()
	err := w.createCsvMergeEntity(webFetchWf, wbsTaskName, scrapedWbsRet, uef, colName, emRow, pls)
	if err != nil {
		log.Err(err).Msg("createCsvMergeEntity: failed to marshal")
		return err
	}
	prompts := w.getPromptsMap(websiteScrape)
	if v, ok := w.CommandPrompts[websiteScrape]; ok && v != "" {
		if len(prompts) <= 0 {
			//prompts = []string{v}
		} else {
			w.TaskOverrides[wbsTaskName] = artemis_orchestrations.TaskOverride{
				SystemPromptExt: v,
			}
		}
	}

	w.createWfTaskPromptOverrides(webFetchWf, wbsTaskName, prompts)
	//w.createWfSchemaFieldOverride(webFetchWf, wbsTaskName, "summary", prompts)
	return nil
}

func (w *ExecFlowsActionsRequest) createWfTaskPromptOverrides(wfn, tn string, overrides map[string]string) {
	if _, exists := w.WfTaskOverrides[wfn]; !exists {
		w.WfTaskOverrides[wfn] = make(map[string]artemis_orchestrations.TaskOverride)
	}
	if _, exists := w.WfTaskOverrides[wfn][tn]; exists {
		rps := w.WfTaskOverrides[wfn][tn]
		for k, v := range overrides {
			rps.ReplacePrompts[k] = v
		}
		w.WfTaskOverrides[wfn][tn] = rps
	} else {
		w.WfTaskOverrides[wfn][tn] = artemis_orchestrations.TaskOverride{
			ReplacePrompts:  overrides,
			SystemPromptExt: "",
		}
	}
}

// Can you tell me what this person does in their current role; and the company they work at now?

func (w *ExecFlowsActionsRequest) EmailsValidatorSetup(uef *artemis_entities.EntitiesFilter) error {
	if v, ok := w.Stages[validateEmails]; !ok || !v {
		return nil
	}
	var colName string
	seen := make(map[string]bool)
	var pls []map[string]interface{}
	emRow := make(map[string][]int)
	for r, cvs := range w.ContactsCsv {
		for cname, colValue := range cvs {
			tv := strings.ToLower(cname)
			v, ok := w.StageContactsMap[cname]
			if ((ok && v == validateEmails) || strings.Contains(tv, "email")) && len(colValue) > 0 && strings.Contains(colValue, "@") {
				if len(colName) > 0 && colName != cname {
					log.Info().Interface("colName", colName).Interface("cname", cname).Msg("EmailsValidatorSetup")
					return fmt.Errorf(fmt.Sprintf("Email csv input has duplicate web column: expecting: %s actual: %s", colName, cname))
				}
				colName = cname
				if len(emRow) <= 0 {
					emRow[colValue] = []int{r}
				} else {
					etm := emRow[colValue]
					etm = append(etm, r)
					emRow[colValue] = etm
				}
				if _, ok1 := seen[colValue]; ok1 {
					continue
				}
				pl := make(map[string]interface{})
				pl["email"] = colValue
				pls = append(pls, pl)
				seen[colValue] = true
			}
		}
	}
	if len(pls) == 0 {
		log.Warn().Msg("no emails found")
		return nil
	}
	w.InitMaps()
	w.WfRetrievalOverrides[emailVdWf] = map[string]artemis_orchestrations.RetrievalOverride{
		validemailRetQp: artemis_orchestrations.RetrievalOverride{Payloads: pls},
	}
	emLabel := csvGlobalMergeRetLabel(validemailRetQp)
	labels := artemis_entities.CreateMdLabels([]string{
		fmt.Sprintf("wf:%s", emailVdWf),
		emLabel,
	})
	uef.Labels = append(uef.Labels, emLabel)
	b, err := utils_csv.NewCsvMergeEntityFromSrcBin(colName, emRow)
	if err != nil {
		log.Err(err).Msg("failed to marshal emRow")
		return err
	}
	usre := artemis_entities.UserEntity{
		Nickname: uef.Nickname,
		Platform: uef.Platform,
		MdSlice: []artemis_entities.UserEntityMetadata{
			{
				JsonData: b,
				Labels:   labels,
			},
		},
	}
	log.Info().Interface("labels", labels).Msg("EmailsValidatorSetup")
	w.WorkflowEntitiesOverrides[emailVdWf] = append(w.WorkflowEntitiesOverrides[emailVdWf], usre)
	w.CustomBasePeriodStepSize = 24
	w.CustomBasePeriodStepSizeUnit = "hours"
	w.CustomBasePeriod = true
	w.Workflows = append(w.Workflows, artemis_orchestrations.WorkflowTemplate{
		WorkflowName: emailVdWf,
	})
	return nil
}

func (w *ExecFlowsActionsRequest) GoogleSearchSetup(uef *artemis_entities.EntitiesFilter) error {
	if v, ok := w.Stages[googleSearch]; !ok || !v {
		return nil
	}
	var colName string
	//seen := make(map[string]bool)
	prompts := w.getPrompts()
	emRow := make(map[string][]int)
	var pls []map[string]interface{}
	for r, cvs := range w.ContactsCsv {
		pl := make(map[string]interface{})
		var colVal string
		for cname, colValue := range cvs {
			// "company",
			if v, ok := w.StageContactsMap[cname]; ok && v == googleSearch {
				//fmt.Println(r, v, colName, colValue, seen)
				colVal = colValue
				tnm := cname
				if newVarName, ok1 := w.StageContactsOverrideMap[cname]; ok1 && len(newVarName) > 0 {
					tnm = newVarName
				}
				w.ContactsCsv[r][cname] = colValue
				// only for payload not override csv value
				pl[tnm] = w.ContactsCsv[r][cname]
				//seen[v] = true
			}
		}
		if len(emRow) <= 0 {
			emRow[colVal] = []int{r}
		} else {
			etm := emRow[colVal]
			etm = append(etm, r)
			emRow[colVal] = etm
		}
		pls = append(pls, pl)
	}
	if len(pls) == 0 {
		log.Warn().Msg("no profiles found")
		return nil
	}
	w.InitMaps()
	err := w.createCsvMergeEntity4(googCsvWf, googleTaskName, googleRetName, uef, colName, emRow, pls)
	if err != nil {
		log.Err(err).Msg("createCsvMergeEntity: failed to marshal")
		return err
	}
	if v, ok := w.CommandPrompts[googleSearch]; ok {
		if len(prompts) <= 0 {
			if v == "" {
				v = "Can you tell me their role and responsibilities?"
			}
			prompts = []string{v}
		} else {
			tmp := w.TaskOverrides[googleTaskName]
			tmp.SystemPromptExt = v
			w.TaskOverrides[googleTaskName] = tmp
		}
		w.createWfTaskPromptOverrides(googCsvWf, googleTaskName, w.getPromptsMap(googleSearch))
	}
	return nil
}

// ReplaceAndPassParams replaces placeholders in the route with URL-encoded values from the provided map.
func ReplaceAndPassParams(route string, params map[string]interface{}) (string, []string, error) {
	// Compile a regular expression to find {param} patterns
	re, err := regexp.Compile(`\{([^\{\}]+)\}`)
	if err != nil {
		log.Err(err).Msg("failed to compile regular expression")
		return "", nil, err // Return an error if the regular expression compilation fails
	}
	var qps []string
	// Replace each placeholder with the corresponding URL-encoded value from the map
	replacedRoute := re.ReplaceAllStringFunc(route, func(match string) string {
		// Extract the parameter name from the match, excluding the surrounding braces
		paramName := match[1 : len(match)-1]
		// Look up the paramName in the params map
		if value, ok := params[paramName]; ok {
			// Delete the matched entry from the map
			if rs, rok := value.(string); rok {
				qps = append(qps, rs)
			}
			delete(params, paramName)
			// If the value exists, convert it to a string and URL-encode it
			return fmt.Sprint(value)
		}
		// If no matching paramName is found in the map, return the match unchanged
		return match
	})

	return replacedRoute, qps, nil
}

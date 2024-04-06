package zeus_v1_ai

import (
	"context"
	"encoding/json"
	"fmt"
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
	webFetchWf = "website-analysis-wf"
	emailVdWf  = "validate-emails-wf"
	googWf     = "google-query-regex-index-wf"
	liWf       = "linkedin-rapid-api-profiles-wf"
	liBizWf    = "linkedin-rapid-api-biz-profiles-wf"

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
	// task
	wbsTaskName      = "website-analysis"
	linkedInTaskName = "linkedin-profile-summary"

	// work labels
	csvSrcGlobalLabel      = "csv:global:source"
	csvSrcGlobalMergeLabel = "csv:global:merge"
)

func (w *ExecFlowsActionsRequest) SetupFlow(ctx context.Context, ou org_users.OrgUser) (*artemis_entities.EntitiesFilter, error) {
	uef := &artemis_entities.EntitiesFilter{
		Platform: "flows",
		Labels:   []string{csvSrcGlobalLabel},
	}

	w.InitMaps()
	err := w.ConvertToCsvStrToMap()
	if err != nil {
		log.Err(err).Interface("w", w).Msg("EmailsValidatorSetup failed")
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
	uef, err = w.SaveCsvImports(ctx, ou, uef)
	if err != nil {
		log.Err(err).Interface("w", w).Msg("EmailsValidatorSetup failed")
		return nil, err
	}
	if uef != nil && uef.Nickname != "" && uef.Platform != "" {
		w.WorkflowEntityRefs = append(w.WorkflowEntityRefs, *uef)
	}
	return uef, err
}

func (w *ExecFlowsActionsRequest) SaveCsvImports(ctx context.Context, ou org_users.OrgUser, uef *artemis_entities.EntitiesFilter) (*artemis_entities.EntitiesFilter, error) {
	if uef == nil || len(w.FlowsActionsRequest.ContactsCsvStr) == 0 {
		return nil, nil
	}
	cs, err := utils_csv.PayloadToCsvString(w.ContactsCsv)
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
			if (strings.Contains(cname, "web") || strings.Contains(cname, "url") || strings.Contains(cname, "link") || strings.Contains(cname, "site")) && len(colValue) > 0 {
				if len(colName) > 0 && colName != cname {
					log.Info().Interface("colName", colName).Interface("cname", cname).Msg("EmailsValidatorSetup")
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
	err := w.createCsvMergeEntity(webFetchWf, wbsTaskName, wbsTaskName, uef, colName, emRow, pls)
	if err != nil {
		log.Err(err).Msg("createCsvMergeEntity: failed to marshal")
		return err
	}
	prompts := w.getPrompts()
	if v, ok := w.CommandPrompts[websiteScrape]; ok && v != "" {
		if len(prompts) <= 0 {
			prompts = []string{v}
		} else {
			w.TaskOverrides[wbsTaskName] = artemis_orchestrations.TaskOverride{
				SystemPromptExt: v,
			}
		}
	}
	w.createWfSchemaFieldOverride(webFetchWf, wbsTaskName, "summary", prompts)
	return nil
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
			if strings.Contains(tv, "email") && len(colValue) > 0 && strings.Contains(colValue, "@") {
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
				if _, ok := seen[colValue]; ok {
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
	b, err := json.Marshal(w.ContactsCsv)
	if err != nil {
		log.Err(err).Msg("failed to marshal gs")
		return err
	}
	if w.TaskOverrides == nil {
		w.TaskOverrides = make(map[string]artemis_orchestrations.TaskOverride)
	}
	w.TaskOverrides["zeusfyi-verbatim"] = artemis_orchestrations.TaskOverride{ReplacePrompt: string(b)}
	if v, ok := w.CommandPrompts[googleSearch]; ok && v != "" {
		//if w.SchemaFieldOverrides == nil {
		//	w.SchemaFieldOverrides = make(map[string]map[string]artemis_orchestrations.)
		//	w.SchemaFieldOverrides["google-results-agg"] = map[string]string{
		//		"summary": v,
		//	}
		//}
	}
	w.Workflows = append(w.Workflows, artemis_orchestrations.WorkflowTemplate{
		WorkflowName: googWf,
	})
	return nil
}

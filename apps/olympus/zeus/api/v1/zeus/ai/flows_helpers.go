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
	validemailRetQp = "validemail-query-params"
)

// S3GlobalOrgImports

func (w *ExecFlowsActionsRequest) TestCsvParser() error {
	for _, r := range w.FlowsActionsRequest.ContactsCsv {
		fmt.Println("r", r)
		for _, c := range r {
			fmt.Println("c", c)
		}
	}
	return fmt.Errorf("fake err")
}

func (w *ExecFlowsActionsRequest) SetupFlow(ctx context.Context, ou org_users.OrgUser) (*artemis_entities.EntitiesFilter, error) {
	uef := &artemis_entities.EntitiesFilter{
		Platform: "flows",
		Labels:   []string{"csv:source", "flows:global"},
	}
	err := w.ConvertToCsvStrToMap()
	if err != nil {
		log.Err(err).Interface("w", w).Msg("EmailsValidatorSetup failed")
		return nil, err
	}
	uef, err = w.SaveCsvImports(ctx, ou, uef)
	if err != nil {
		log.Err(err).Interface("w", w).Msg("EmailsValidatorSetup failed")
		return nil, err
	}
	err = w.EmailsValidatorSetup(uef)
	if err != nil {
		log.Err(err).Interface("w", w).Msg("EmailsValidatorSetup failed")
		return nil, err
	}
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
	err = w.ScrapeRegularWebsiteSetup(uef)
	if err != nil {
		log.Err(err).Interface("w", w).Msg("ScrapeRegularWebsiteSetup failed")
		return nil, err
	}
	if uef != nil && uef.Nickname != "" && uef.Platform != "" {
		w.WorkflowEntityRefs = append(w.WorkflowEntityRefs, *uef)
	}
	err = w.TestCsvParser()
	return uef, err
}

func (w *ExecFlowsActionsRequest) SaveCsvImports(ctx context.Context, ou org_users.OrgUser, uef *artemis_entities.EntitiesFilter) (*artemis_entities.EntitiesFilter, error) {
	if uef == nil || len(w.FlowsActionsRequest.ContactsCsvStr) == 0 {
		return nil, nil
	}
	lbs := append(uef.Labels, "csv:contacts")
	usre := &artemis_entities.UserEntity{
		Nickname: uef.Nickname,
		Platform: uef.Platform,
		MdSlice: []artemis_entities.UserEntityMetadata{
			{
				TextData: aws.String(w.ContactsCsvStr),
				Labels:   artemis_entities.CreateMdLabels(lbs),
			},
		},
	}
	if len(usre.Nickname) <= 0 || len(usre.Platform) <= 0 || len(usre.MdSlice) <= 0 {
		return nil, fmt.Errorf("no entities name")
	}
	_, err := ai_platform_service_orchestrations.S3GlobalOrgImports(ctx, ou, usre)
	if err != nil {
		log.Err(err).Msg("SaveImport: error")
		return nil, err
	}
	log.Info().Interface("ue", usre.Nickname).Msg("entity hash")
	return uef, nil
}

func (w *ExecFlowsActionsRequest) ScrapeRegularWebsiteSetup(uef *artemis_entities.EntitiesFilter) error {
	if v, ok := w.Stages[websiteScrape]; !ok || !v {
		return nil
	}
	seen := make(map[string]bool)
	var pls []map[string]interface{}
	for _, cv := range w.ContactsCsv {
		for em, emv := range cv {
			//if _, ok := seen[tv]; ok {
			//	continue
			//}
			if (strings.Contains(em, "web") || strings.Contains(em, "url") || strings.Contains(em, "link") || strings.Contains(em, "site")) && len(emv) > 0 {
				pl := make(map[string]interface{})
				uv, err := convertToHTTPS(emv)
				if err != nil {
					log.Err(err).Msg("failed to convert url to https")
					continue
				}
				if _, ok := seen[uv]; ok {
					continue
				}
				if strings.HasPrefix(uv, "https://www.linkedin.com") || strings.HasPrefix(uv, "https://linkedin.com") {
					continue
				}
				seen[uv] = true
				pl["url"] = uv
				pls = append(pls, pl)
			}
			//seen[tv] = true
		}
	}
	if len(pls) == 0 {
		log.Warn().Msg("no urls found")
		return nil
	}
	if w.TaskOverrides == nil {
		w.TaskOverrides = make(map[string]artemis_orchestrations.TaskOverride)
	}
	if w.RetrievalOverrides == nil {
		w.RetrievalOverrides = make(map[string]artemis_orchestrations.RetrievalOverride)
	}
	if v, ok := w.CommandPrompts[websiteScrape]; ok && v != "" {
		if w.SchemaFieldOverrides == nil {
			w.SchemaFieldOverrides = make(map[string]map[string]string)
			w.SchemaFieldOverrides["website-analysis"] = map[string]string{
				"summary": v,
			}
		}
	}
	w.RetrievalOverrides["website-analysis"] = artemis_orchestrations.RetrievalOverride{Payloads: pls}
	w.Workflows = append(w.Workflows, artemis_orchestrations.WorkflowTemplate{
		WorkflowName: webFetchWf,
	})
	return nil
}

func (w *ExecFlowsActionsRequest) EmailsValidatorSetup(uef *artemis_entities.EntitiesFilter) error {
	if v, ok := w.Stages[validateEmails]; !ok || !v {
		return nil
	}
	var colName string
	seen := make(map[string]bool)
	var pls []map[string]interface{}
	emRow := make(map[string][]int)
	for r, cv := range w.ContactsCsv {
		for cn, emv := range cv {
			tv := strings.ToLower(cn)
			if strings.Contains(tv, "email") && len(emv) > 0 && strings.Contains(emv, "@") {
				if len(colName) > 0 && colName != cn {
					return fmt.Errorf("duplicate email column")
				}
				colName = emv
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
	if len(pls) == 0 {
		log.Warn().Msg("no emails found")
		return nil
	}
	b, err := json.Marshal(emRow)
	if err != nil {
		log.Err(err).Msg("failed to marshal gs")
		return err
	}
	w.RetrievalOverrides[validemailRetQp] = artemis_orchestrations.RetrievalOverride{Payloads: pls}
	labels := artemis_entities.CreateMdLabels([]string{
		fmt.Sprintf("wf:%s", emailVdWf),
		"csv:merge",
		fmt.Sprintf("merge:ret:%s", validemailRetQp),
	})
	usre := artemis_entities.UserEntity{
		Nickname: uef.Nickname,
		Platform: uef.Platform,
		MdSlice: []artemis_entities.UserEntityMetadata{
			{
				TextData: aws.String(colName),
				JsonData: b,
				Labels:   labels,
			},
		},
	}
	w.WorkflowEntities = append(w.WorkflowEntities, usre)
	if w.TaskOverrides == nil {
		w.TaskOverrides = make(map[string]artemis_orchestrations.TaskOverride)
	}
	if w.RetrievalOverrides == nil {
		w.RetrievalOverrides = make(map[string]artemis_orchestrations.RetrievalOverride)
	}
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
		if w.SchemaFieldOverrides == nil {
			w.SchemaFieldOverrides = make(map[string]map[string]string)
			w.SchemaFieldOverrides["google-results-agg"] = map[string]string{
				"summary": v,
			}
		}
	}
	w.Workflows = append(w.Workflows, artemis_orchestrations.WorkflowTemplate{
		WorkflowName: googWf,
	})
	return nil
}

// Can you tell me what this person does in their current role; and the company they work at now?

func (w *ExecFlowsActionsRequest) LinkedInScraperSetup(uef *artemis_entities.EntitiesFilter) error {
	if v, ok := w.Stages[linkedIn]; !ok || !v {
		return nil
	}
	b, err := json.Marshal(w.ContactsCsv)
	if err != nil {
		log.Err(err).Msg("failed to marshal li")
		return err
	}
	if w.TaskOverrides == nil {
		w.TaskOverrides = make(map[string]artemis_orchestrations.TaskOverride)
	}
	w.TaskOverrides["linkedin-profiles-rapid-api-qps"] = artemis_orchestrations.TaskOverride{ReplacePrompt: string(b)}
	if v, ok := w.CommandPrompts[linkedIn]; ok && v != "" {
		if w.SchemaFieldOverrides == nil {
			w.SchemaFieldOverrides = make(map[string]map[string]string)
			w.SchemaFieldOverrides["results-agg"] = map[string]string{
				"summary": v,
			}
		}
	}
	w.Workflows = append(w.Workflows, artemis_orchestrations.WorkflowTemplate{
		WorkflowName: liWf,
	})
	return nil
}

func (w *ExecFlowsActionsRequest) LinkedInBizScraperSetup() error {
	if v, ok := w.Stages[linkedInBiz]; !ok || !v {
		return nil
	}
	b, err := json.Marshal(w.ContactsCsv)
	if err != nil {
		log.Err(err).Msg("failed to marshal linkedInBiz")
		return err
	}
	if w.TaskOverrides == nil {
		w.TaskOverrides = make(map[string]artemis_orchestrations.TaskOverride)
	}
	w.TaskOverrides["linkedin-biz-profiles-rapid-api-qps"] = artemis_orchestrations.TaskOverride{ReplacePrompt: string(b)}
	if v, ok := w.CommandPrompts[linkedInBiz]; ok && v != "" {
		if w.SchemaFieldOverrides == nil {
			w.SchemaFieldOverrides = make(map[string]map[string]string)
			w.SchemaFieldOverrides["results-agg"] = map[string]string{
				"summary": v,
			}
		}
	}
	w.Workflows = append(w.Workflows, artemis_orchestrations.WorkflowTemplate{
		WorkflowName: liBizWf,
	})
	return nil
}

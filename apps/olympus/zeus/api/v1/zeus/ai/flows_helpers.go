package zeus_v1_ai

import (
	"encoding/json"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
)

const (
	siteFetchWf = "analysis-website-wf"
	emailVdWf   = "validate-emails-wf"
	googWf      = "google-query-regex-index-wf"
	liWf        = "linkedin-rapid-api-profiles-wf"
	liBizWf     = "linkedin-rapid-api-biz-profiles-wf"
)

func (w *ExecFlowsActionsRequest) ScrapeRegularWebsiteSetup() error {
	if v, ok := w.Stages["websiteScrape"]; !ok || !v {
		return nil
	}

	var pls []map[string]interface{}
	for _, cv := range w.ContactsCsv {
		for em, emv := range cv {
			tv := strings.ToLower(em)
			if (strings.Contains(tv, "web") || strings.Contains(tv, "url") || strings.Contains(tv, "link") || strings.Contains(tv, "site")) && len(emv) > 0 {
				pl := make(map[string]interface{})
				pl["url"] = emv
				pls = append(pls, pl)
			}
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
	if v, ok := w.CommandPrompts["websiteScrape"]; ok && v != "" {
		if w.SchemaFieldOverrides == nil {
			w.SchemaFieldOverrides = make(map[string]map[string]string)
			w.SchemaFieldOverrides["website-analysis"] = map[string]string{
				"summary": v,
			}
		}
	}
	w.RetrievalOverrides["website-analysis"] = artemis_orchestrations.RetrievalOverride{Payloads: pls}
	w.Workflows = append(w.Workflows, artemis_orchestrations.WorkflowTemplate{
		WorkflowName: siteFetchWf,
	})
	return nil
}

func (w *ExecFlowsActionsRequest) EmailsValidatorSetup() error {
	if v, ok := w.Stages["validateEmails"]; !ok || !v {
		return nil
	}

	var pls []map[string]interface{}
	for _, cv := range w.ContactsCsv {
		for em, emv := range cv {
			if strings.Contains(strings.ToLower(em), "email") && len(emv) > 0 {
				pl := make(map[string]interface{})
				pl["email"] = emv
				pls = append(pls, pl)
			}
		}
	}
	if w.TaskOverrides == nil {
		w.TaskOverrides = make(map[string]artemis_orchestrations.TaskOverride)
	}
	if w.RetrievalOverrides == nil {
		w.RetrievalOverrides = make(map[string]artemis_orchestrations.RetrievalOverride)
	}
	w.RetrievalOverrides["validemail-query-params"] = artemis_orchestrations.RetrievalOverride{Payloads: pls}
	w.Workflows = append(w.Workflows, artemis_orchestrations.WorkflowTemplate{
		WorkflowName: emailVdWf,
	})
	return nil
}

func (w *ExecFlowsActionsRequest) GoogleSearchSetup() error {
	if v, ok := w.Stages["googleSearch"]; !ok || !v {
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
	if v, ok := w.CommandPrompts["googleSearch"]; ok && v != "" {
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

func (w *ExecFlowsActionsRequest) LinkedInScraperSetup() error {
	if v, ok := w.Stages["linkedIn"]; !ok || !v {
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
	if v, ok := w.CommandPrompts["linkedIn"]; ok && v != "" {
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
	if v, ok := w.Stages["linkedInBiz"]; !ok || !v {
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
	if v, ok := w.CommandPrompts["linkedInBiz"]; ok && v != "" {
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

package zeus_v1_ai

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	utils_csv "github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/csv"
)

func (w *ExecFlowsActionsRequest) LinkedInScraperSetup(uef *artemis_entities.EntitiesFilter) error {
	if v, ok := w.Stages[linkedIn]; !ok || !v {
		return nil
	}
	var colName string
	seen := make(map[string]bool)
	emRow := make(map[string][]int)
	var pls []map[string]interface{}
	for r, cvs := range w.ContactsCsv {
		for cname, colValue := range cvs {
			// "company",
			v, ok := w.StageContactsMap[cname]
			if (strings.Contains(strings.ToLower(cname), "linkedin") || (ok && v == linkedIn)) && strings.Contains(strings.ToLower(colValue), "linkedin.com/in") && len(colValue) > 0 {
				if len(colName) > 0 && colName != cname {
					log.Info().Interface("colName", colName).Interface("cname", cname).Msg("LinkedInScraperSetup")
					return fmt.Errorf(fmt.Sprintf("LinkedIn csv input has duplicate web column: expecting: %s actual: %s", colName, cname))
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
				if _, ok1 := seen[uv]; ok1 {
					w.ContactsCsv[r][cname] = uv
					continue
				}
				if strings.Contains(strings.ToLower(uv), "linkedin.com/in") {
					pl := make(map[string]interface{})
					w.ContactsCsv[r][cname] = uv
					pl["linkedin_url"] = w.ContactsCsv[r][cname]
					pls = append(pls, pl)
					seen[uv] = true
				}
			}
		}
	}
	if len(pls) == 0 {
		log.Warn().Msg("no profiles found")
		return nil
	}
	w.InitMaps()
	err := w.createCsvMergeEntity(liWf, wbsTaskName, linkedInRetQp, uef, colName, emRow, pls)
	if err != nil {
		log.Err(err).Msg("createCsvMergeEntity: failed to marshal")
		return err
	}
	if v, ok := w.CommandPrompts[linkedIn]; ok {
		prompts := w.getPrompts()
		if len(prompts) <= 0 {
			if v == "" {
				v = "Can you tell me their role and responsibilities?"
			}
			prompts = []string{v}
		} else {
			tmp := w.TaskOverrides[wbsTaskName]
			tmp.SystemPromptExt = v
			w.TaskOverrides[wbsTaskName] = tmp
		}
		w.createWfTaskPromptOverrides(liWf, wbsTaskName, w.getPromptsMap(linkedIn))
		//w.createWfSchemaFieldOverride(liWf, wbsTaskName, "summary", prompts)
	}
	return nil
}

func (w *ExecFlowsActionsRequest) LinkedInBizScraperSetup(uef *artemis_entities.EntitiesFilter) error {
	if v, ok := w.Stages[linkedInBiz]; !ok || !v {
		return nil
	}
	var colName string
	seen := make(map[string]bool)
	emRow := make(map[string][]int)
	var pls []map[string]interface{}
	for r, cvs := range w.ContactsCsv {
		for cname, colValue := range cvs {
			// "company",
			v, ok := w.StageContactsMap[cname]
			if (strings.Contains(strings.ToLower(cname), "linkedin") || (ok && v == linkedInBiz)) && strings.Contains(strings.ToLower(colValue), "linkedin.com/company") && len(colValue) > 0 {
				if len(colName) > 0 && colName != cname {
					log.Info().Interface("colName", colName).Interface("cname", cname).Msg("LinkedInBizScraperSetup")
					return fmt.Errorf(fmt.Sprintf("LinkedInBiz csv input has duplicate web column: expecting: %s actual: %s", colName, cname))
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
				if _, ok1 := seen[uv]; ok1 {
					w.ContactsCsv[r][cname] = uv
					continue
				}
				if strings.Contains(strings.ToLower(uv), "linkedin.com/company") {
					pl := make(map[string]interface{})
					w.ContactsCsv[r][cname] = uv
					pl["linkedin_url"] = w.ContactsCsv[r][cname]
					pls = append(pls, pl)
					seen[uv] = true
				}
			}
		}
	}
	if len(pls) == 0 {
		log.Warn().Msg("no profiles found")
		return nil
	}
	w.InitMaps()
	err := w.createCsvMergeEntity(liBizWf, wbsTaskName, linkedInBizRetQp, uef, colName, emRow, pls)
	if err != nil {
		log.Err(err).Msg("createCsvMergeEntity: failed to marshal")
		return err
	}
	prompts := w.getPrompts()
	if v, ok := w.CommandPrompts[linkedInBiz]; ok {
		if len(prompts) <= 0 {
			if v == "" {
				v = "Can you tell me their role and responsibilities?"
			}
			prompts = []string{v}
		} else {
			tmp := w.TaskOverrides[wbsTaskName]
			tmp.SystemPromptExt = v
			w.TaskOverrides[wbsTaskName] = tmp
		}
		w.createWfTaskPromptOverrides(liBizWf, wbsTaskName, w.getPromptsMap(linkedInBiz))
		//w.createWfSchemaFieldOverride(liBizWf, wbsTaskName, "summary", prompts)
	}
	return nil
}

func (w *ExecFlowsActionsRequest) createWfSchemaFieldOverride(wfn, schN, fN string, overrides []string) {
	if _, exists := w.WfSchemaFieldOverrides[wfn]; !exists {
		w.WfSchemaFieldOverrides[wfn] = make(map[string]map[string][]string)
	}
	if _, exists := w.WfSchemaFieldOverrides[wfn][schN]; !exists {
		w.WfSchemaFieldOverrides[wfn][schN] = make(map[string][]string)
	}
	// Check if the field name exists; if it does, append the overrides, else create a new entry
	if existingOverrides, exists := w.WfSchemaFieldOverrides[wfn][schN][fN]; exists {
		w.WfSchemaFieldOverrides[wfn][schN][fN] = append(existingOverrides, overrides...)
	} else {
		w.WfSchemaFieldOverrides[wfn][schN][fN] = overrides
	}
}

// createCsvMergeEntity is standardizing csv payload driven setups
func (w *ExecFlowsActionsRequest) createCsvMergeEntity(wfn, tn, retN string, uef *artemis_entities.EntitiesFilter, colName string, emRow map[string][]int, pls []map[string]interface{}) error {
	wsbLabel := csvGlobalMergeAnalysisTaskLabel(tn)
	labels := artemis_entities.CreateMdLabels([]string{
		fmt.Sprintf("wf:%s", wfn),
		wsbLabel,
	})
	uef.Labels = append(uef.Labels, wsbLabel)
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
	w.WfRetrievalOverrides[wfn] = map[string]artemis_orchestrations.RetrievalOverride{
		retN: artemis_orchestrations.RetrievalOverride{Payloads: pls},
	}
	w.WorkflowEntitiesOverrides[wfn] = append(w.WorkflowEntitiesOverrides[wfn], usre)
	w.Workflows = append(w.Workflows, artemis_orchestrations.WorkflowTemplate{
		WorkflowName: wfn,
	})
	return err
}

// createCsvMergeEntity is standardizing csv payload driven setups
func (w *ExecFlowsActionsRequest) createCsvMergeEntity2(wfn, tn, retN string, uef *artemis_entities.EntitiesFilter, colName string, emRow map[string][]int, pls []map[string]interface{}) error {
	wsbLabel := csvGlobalMergeAnalysisTaskLabel(tn)
	labels := artemis_entities.CreateMdLabels([]string{
		fmt.Sprintf("wf:%s", wfn),
		wsbLabel,
	})
	uef.Labels = append(uef.Labels, wsbLabel)
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
	w.WfRetrievalOverrides[wfn] = map[string]artemis_orchestrations.RetrievalOverride{
		retN: artemis_orchestrations.RetrievalOverride{Payloads: pls},
	}
	w.WorkflowEntitiesOverrides[wfn] = append(w.WorkflowEntitiesOverrides[wfn], usre)
	return err
}

func (w *ExecFlowsActionsRequest) createCsvMergeEntity4(wfn, tn, retN string, uef *artemis_entities.EntitiesFilter, colName string, emRow map[string][]int, pls []map[string]interface{}) error {
	wsbLabel := csvGlobalMergeAnalysisTaskLabel(tn)
	labels := artemis_entities.CreateMdLabels([]string{
		fmt.Sprintf("wf:%s", wfn),
		wsbLabel,
	})
	uef.Labels = append(uef.Labels, wsbLabel)
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

	var nps []map[string]interface{}
	for _, pl := range pls {
		prompts := w.getPromptsMap(googleSearch)
		for _, pv := range prompts {
			tmp := make(map[string]interface{})
			for k, v := range pl {
				tmp[k] = v
			}
			nv, _, nerr := ReplaceAndPassParams(pv, tmp)
			if nerr != nil {
				log.Err(nerr).Msg("failed to ReplaceAndPassParams")
				return nerr
			}
			nmv := make(map[string]interface{})
			nmv["q"] = url.QueryEscape(nv)
			nps = append(nps, nmv)
		}
	}
	w.WfRetrievalOverrides[wfn] = map[string]artemis_orchestrations.RetrievalOverride{
		retN: artemis_orchestrations.RetrievalOverride{Payloads: nps},
	}
	w.WorkflowEntitiesOverrides[wfn] = append(w.WorkflowEntitiesOverrides[wfn], usre)
	w.Workflows = append(w.Workflows, artemis_orchestrations.WorkflowTemplate{
		WorkflowName: wfn,
	})
	return err
}

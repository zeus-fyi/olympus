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

func csvGlobalRetLabel() string {
	return fmt.Sprintf("%s:ret", csvSrcGlobalMergeLabel)
}
func csvGlobalAnalysisTaskLabel() string {
	return fmt.Sprintf("%s:analysis:task", csvSrcGlobalMergeLabel)
}

func csvGlobalMergeAnalysisTaskLabel(tn string) string {
	return fmt.Sprintf("%s:%s", csvGlobalAnalysisTaskLabel(), tn)
}

func csvGlobalMergeRetLabel(rn string) string {
	return fmt.Sprintf("%s:%s", csvGlobalRetLabel(), rn)
}

func doesNotContain(sin string, notAllowed []string) bool {
	for _, ns := range notAllowed {
		if strings.Contains(strings.ToLower(sin), strings.ToLower(ns)) {
			return false
		}
	}
	return true
}

func isValidURL(inputURL string) (*url.URL, error) {
	u, err := url.Parse(inputURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("URL must be http or https, got: %s", u.Scheme)
	}
	if u.Host == "" {
		return nil, fmt.Errorf("URL must have a host")
	}
	return u, nil
}
func convertToHTTPS(inputURL string) (string, error) {
	u, err := isValidURL(inputURL)
	if err != nil {
		return "", err
	}

	if u.Scheme == "http" {
		u.Scheme = "https"
	}
	return u.String(), nil
}

func (w *ExecFlowsActionsRequest) TestCsvParser() error {
	for _, r := range w.FlowsActionsRequest.ContactsCsv {
		fmt.Println("r", r)
		for _, c := range r {
			fmt.Println("c", c)
		}
	}
	return fmt.Errorf("fake err")
}

func (w *ExecFlowsActionsRequest) InitMaps() {
	tmp := make(map[string]string)
	for k, v := range w.StagePromptMap {
		tmp[strings.TrimSpace(k)] = strings.TrimSpace(v)
	}
	w.StagePromptMap = tmp

	if w.WorkflowEntitiesOverrides == nil {
		w.WorkflowEntitiesOverrides = make(map[string][]artemis_entities.UserEntity)
	}
	// task level
	if w.TaskOverrides == nil {
		w.TaskOverrides = make(map[string]artemis_orchestrations.TaskOverride)
	}
	if w.WfTaskOverrides == nil {
		w.WfTaskOverrides = make(map[string]artemis_orchestrations.TaskOverrides)
	}
	if w.SchemaFieldOverrides == nil {
		w.SchemaFieldOverrides = make(map[string]map[string][]string)
	}
	// wf level
	if w.WfRetrievalOverrides == nil {
		tmp1 := make(map[string]artemis_orchestrations.RetrievalOverrides)
		w.WfRetrievalOverrides = tmp1
	}
	if w.WfSchemaFieldOverrides == nil {
		w.WfSchemaFieldOverrides = make(map[string]artemis_orchestrations.SchemaOverrides)
	}
}

func (w *ExecFlowsActionsRequest) getPromptsMap(stage string) map[string]string {
	prompts := make(map[string]string)
	for _, cvs := range w.PromptsCsv {
		for cn, colValue := range cvs {
			v, ok := w.StagePromptMap[cn]
			if ok && (v == stage || strings.ToLower(v) == "default") {
				cnStage := fmt.Sprintf("%s_%s", stage, cn)
				prompts[cnStage] = colValue
			}
		}
	}
	return prompts
}

func (w *ExecFlowsActionsRequest) getPrompts() []string {
	var prompts []string
	for _, cvs := range w.PromptsCsv {
		for _, colValue := range cvs {
			prompts = append(prompts, colValue)
		}
	}
	return prompts
}

func (w *ExecFlowsActionsRequest) ConvertToCsvStrToMap() ([]string, error) {
	var headersCsv []string
	if len(w.FlowsActionsRequest.ContactsCsvStr) > 0 {
		headers, err := utils_csv.ParseCsvStringOrderedHeaders(w.FlowsActionsRequest.ContactsCsvStr)
		if err != nil {
			log.Err(err).Msg("ConvertToCsvStrToMap: error")
			return nil, err
		}
		headersCsv = headers
		cv, err := utils_csv.ParseCsvStringToMap(w.FlowsActionsRequest.ContactsCsvStr)
		if err != nil {
			log.Err(err).Msg("SaveCsvImports: ContactsCsvStr: error")
			return nil, err
		}
		// Check if PreviewCount is set and limit the number of records accordingly
		if w.PreviewCount > 0 && len(cv) > w.PreviewCount {
			cv = cv[:w.PreviewCount]
		}
		tm, err := utils_csv.PayloadToCsvString(cv)
		if err != nil {
			log.Err(err).Msg("SaveCsvImports: ContactsCsvStr: error")
			return nil, err
		}
		w.ContactsCsvStr = tm
		w.ContactsCsv = cv
	}
	if len(w.FlowsActionsRequest.PromptsCsvStr) > 0 {
		pcv, err := utils_csv.ParseCsvStringToMap(w.FlowsActionsRequest.PromptsCsvStr)
		if err != nil {
			log.Err(err).Msg("SaveCsvImports: PromptsCsvStr: error")
			return nil, err
		}
		var tslice []map[string]string
		for _, vi := range pcv {
			for k, v := range vi {
				if v == "" {
					delete(vi, k)
				}
			}
			if len(vi) > 0 {
				tslice = append(tslice, vi)
			}
		}
		pcv = tslice
		for r, mv := range pcv {
			tmp := make(map[string]string)
			for k, v := range mv {
				tmp[strings.TrimSpace(k)] = strings.TrimSpace(v)
			}
			pcv[r] = tmp
		}
		tm, err := utils_csv.PayloadToCsvString(pcv)
		if err != nil {
			log.Err(err).Msg("SaveCsvImports: ContactsCsvStr: error")
			return nil, err
		}
		w.PromptsCsvStr = tm
		w.PromptsCsv = pcv
	}
	return headersCsv, nil
}

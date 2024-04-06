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
	if w.WorkflowEntitiesOverrides == nil {
		w.WorkflowEntitiesOverrides = make(map[string][]artemis_entities.UserEntity)
	}
	// task level
	if w.TaskOverrides == nil {
		w.TaskOverrides = make(map[string]artemis_orchestrations.TaskOverride)
	}
	if w.SchemaFieldOverrides == nil {
		w.SchemaFieldOverrides = make(map[string]map[string][]string)
	}
	// wf level
	if w.WfRetrievalOverrides == nil {
		tmp := make(map[string]artemis_orchestrations.RetrievalOverrides)
		w.WfRetrievalOverrides = tmp
	}
	if w.WfSchemaFieldOverrides == nil {
		w.WfSchemaFieldOverrides = make(map[string]artemis_orchestrations.SchemaOverrides)
	}
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

func (w *ExecFlowsActionsRequest) ConvertToCsvStrToMap() error {
	if len(w.FlowsActionsRequest.ContactsCsvStr) > 0 {
		cv, err := utils_csv.ParseCsvStringToMap(w.FlowsActionsRequest.ContactsCsvStr)
		if err != nil {
			log.Err(err).Msg("SaveCsvImports: ContactsCsvStr: error")
			return err
		}
		// Check if PreviewCount is set and limit the number of records accordingly
		if w.PreviewCount > 0 && len(cv) > w.PreviewCount {
			cv = cv[:w.PreviewCount]
		}
		w.ContactsCsv = cv
	}
	if len(w.FlowsActionsRequest.PromptsCsvStr) > 0 {
		pcv, err := utils_csv.ParseCsvStringToMap(w.FlowsActionsRequest.PromptsCsvStr)
		if err != nil {
			log.Err(err).Msg("SaveCsvImports: PromptsCsvStr: error")
			return err
		}
		w.PromptsCsv = pcv
	}

	return nil
}

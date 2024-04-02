package zeus_v1_ai

import (
	"fmt"
	"net/url"

	"github.com/rs/zerolog/log"
)

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

func (w *ExecFlowsActionsRequest) ConvertToCsvStrToMap() error {
	if len(w.FlowsActionsRequest.ContactsCsvStr) > 0 {
		cv, err := parseCsvStringToMap(w.FlowsActionsRequest.ContactsCsvStr)
		if err != nil {
			log.Err(err).Msg("SaveCsvImports: error")
			return err
		}
		w.ContactsCsv = cv
	}
	if len(w.FlowsActionsRequest.PromptsCsvStr) > 0 {
		pcv, err := parseCsvStringToMap(w.FlowsActionsRequest.PromptsCsvStr)
		if err != nil {
			log.Err(err).Msg("SaveCsvImports: error")
			return err
		}
		w.PromptsCsv = pcv
	}
	return nil
}

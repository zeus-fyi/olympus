package ai_platform_service_orchestrations

import (
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
)

func getPayloads(cp *MbChildSubProcessParams) []map[string]interface{} {
	var echoReqs []map[string]interface{}
	if cp.WfExecParams.WorkflowOverrides.RetrievalOverrides != nil {
		if v, ok := cp.WfExecParams.WorkflowOverrides.RetrievalOverrides[cp.Tc.Retrieval.RetrievalName]; ok {
			for _, pl := range v.Payloads {
				echoReqs = append(echoReqs, pl)
			}
		}
	}
	return echoReqs
}

func getRetOpt(cp *MbChildSubProcessParams, echoReqs []map[string]interface{}) string {
	retOpt := "default"
	if cp.Tc.Retrieval.WebFilters != nil && cp.Tc.Retrieval.WebFilters.PayloadPreProcessing != nil && len(echoReqs) > 0 {
		retOpt = *cp.Tc.Retrieval.WebFilters.PayloadPreProcessing
	}
	return retOpt
}

func FixRegexInput(input string) string {
	if len(input) > 0 {
		// Check if the first character is a backtick and replace it with a double quote
		if input[0] == '`' {
			input = "\"" + input[1:]
		}
		// Check if the last character is a backtick and replace it with a double quote
		if input[len(input)-1] == '`' {
			input = input[:len(input)-1] + "\""
		}
	}
	return input
}

func setDontRetryCodes(retInst artemis_orchestrations.RetrievalItem) []int {
	if retInst.WebFilters.DontRetryStatusCodes != nil {
		return retInst.WebFilters.DontRetryStatusCodes
	}
	return nil
}

func getRegexPatterns(retInst artemis_orchestrations.RetrievalItem) []string {
	var regexPatterns []string
	for _, rgp := range retInst.WebFilters.RegexPatterns {
		regexPatterns = append(regexPatterns, FixRegexInput(rgp))
	}
	return regexPatterns
}

func getRestMethod(retInst artemis_orchestrations.RetrievalItem) string {
	restMethod := http.MethodGet
	if retInst.WebFilters.EndpointREST != nil {
		restMethod = strings.ToLower(*retInst.WebFilters.EndpointREST)
		switch restMethod {
		case "post", "POST":
			restMethod = http.MethodPost
		case "put", "PUT":
			restMethod = http.MethodPut
		case "delete":
			restMethod = http.MethodDelete
		case "patch":
			restMethod = http.MethodPatch
		case "get":
			restMethod = http.MethodGet
		default:
			log.Info().Str("restMethod", restMethod).Msg("ApiCallRequestTask: rest method")
		}
	}
	return restMethod
}

func setHeaders(retInst artemis_orchestrations.RetrievalItem, r RouteTask) RouteTask {
	if retInst.WebFilters.RequestHeaders != nil {
		if r.Headers == nil {
			r.Headers = make(http.Header)
		}
		for k, v := range retInst.WebFilters.RequestHeaders {
			r.Headers.Set(k, v)
		}
	}
	return r
}

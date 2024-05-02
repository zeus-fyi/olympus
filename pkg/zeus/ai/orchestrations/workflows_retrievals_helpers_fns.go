package ai_platform_service_orchestrations

import (
	"fmt"
	"net/url"
	"regexp"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func getRoutesByRule(cp *MbChildSubProcessParams, routes []iris_models.RouteInfo) []iris_models.RouteInfo {
	if cp.Tc.Retrieval.RetrievalItemInstruction.WebFilters != nil &&
		cp.Tc.Retrieval.RetrievalItemInstruction.WebFilters.LbStrategy != nil &&
		*cp.Tc.Retrieval.RetrievalItemInstruction.WebFilters.LbStrategy != lbStrategyPollTable && len(routes) > 1 {
		routes = routes[0:1]
	}
	return routes
}

func getOrgRetIfFlows(cp *MbChildSubProcessParams) org_users.OrgUser {
	tmpOu := cp.Ou
	if cp.WfExecParams.WorkflowOverrides.IsUsingFlows {
		tmpOu.OrgID = FlowsOrgID
	}
	return tmpOu
}

func getFanOutRetPolicy(cp *MbChildSubProcessParams, ao workflow.ActivityOptions) workflow.ActivityOptions {
	cao := ao
	cao.HeartbeatTimeout = time.Minute * 10
	cao.RetryPolicy = GetRetryPolicy(cp.Tc.Retrieval, time.Hour*24)
	cao.RetryPolicy.MaximumAttempts = 10000
	return cao
}

func getPlatform(cp *MbChildSubProcessParams) string {
	platform := cp.Tc.Retrieval.RetrievalPlatform
	if cp.Tc.TriggerActionsApproval.TriggerAction == apiApproval {
		platform = apiApproval
	}
	if cp.Tc.EvalID <= 0 || cp.Tc.TriggerActionsApproval.ApprovalID <= 0 {
		switch platform {
		case twitterPlatform, redditPlatform, discordPlatform, telegramPlatform:
		default:
			platform = webPlatform
		}
	}
	return platform
}

func getRetActRetryPolicy(cp *MbChildSubProcessParams) workflow.ActivityOptions {
	ao := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Hour * 24, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			BackoffCoefficient: 2.5,
			MaximumInterval:    time.Minute * 5,
			MaximumAttempts:    100,
		},
	}
	if cp.Tc.RetryPolicy != nil {
		ao.RetryPolicy = cp.Tc.RetryPolicy
	}
	return ao
}

func GetRetryPolicy(ret artemis_orchestrations.RetrievalItem, maxRunTime time.Duration) *temporal.RetryPolicy {
	if ret.WebFilters == nil {
		return nil
	}
	retry := &temporal.RetryPolicy{
		MaximumInterval: maxRunTime,
	}
	if ret.WebFilters.BackoffCoefficient != nil && *ret.WebFilters.BackoffCoefficient >= 1 {
		retry.BackoffCoefficient = *ret.WebFilters.BackoffCoefficient
	}
	if ret.WebFilters.MaxRetries != nil && *ret.WebFilters.MaxRetries >= 0 {
		retry.MaximumAttempts = int32(*ret.WebFilters.MaxRetries)
	}
	return retry
}

// ExtractParams takes a string of comma-separated regex patterns and a target string.
// It applies each regex pattern to the target string and accumulates all matched groups from each pattern into a single slice.
func ExtractParams(regexStrs []string, strContent []byte) ([]string, error) {
	// Split the regexStr into individual patterns
	var combinedParams []string
	for _, pattern := range regexStrs {
		// Compile and execute each pattern
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err // Return an error if the regular expression compilation fails
		}
		// Find all matches and extract the parameter names
		matches := re.FindAll(strContent, -1)
		for _, match := range matches {
			combinedParams = append(combinedParams, string(match))
		}
	}

	return combinedParams, nil
}

// ReplaceParams replaces placeholders in the route with URL-encoded values from the provided map.
func ReplaceParams(route string, params echo.Map) (string, error) {
	// Compile a regular expression to find {param} patterns
	re, err := regexp.Compile(`\{([^\{\}]+)\}`)
	if err != nil {
		log.Err(err).Msg("failed to compile regular expression")
		return "", err // Return an error if the regular expression compilation fails
	}

	// Replace each placeholder with the corresponding URL-encoded value from the map
	replacedRoute := re.ReplaceAllStringFunc(route, func(match string) string {
		// Extract the parameter name from the match, excluding the surrounding braces
		paramName := match[1 : len(match)-1]
		// Look up the paramName in the params map
		if value, ok := params[paramName]; ok {
			// Delete the matched entry from the map
			delete(params, paramName)
			// If the value exists, convert it to a string and URL-encode it
			return url.QueryEscape(fmt.Sprint(value))
		}
		// If no matching paramName is found in the map, return the match unchanged
		return match
	})

	return replacedRoute, nil
}

// ReplaceAndPassParams replaces placeholders in the route with URL-encoded values from the provided map.
func ReplaceAndPassParams(route string, params echo.Map) (string, []string, error) {
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
			return url.QueryEscape(fmt.Sprint(value))
		}
		// If no matching paramName is found in the map, return the match unchanged
		return match
	})

	return replacedRoute, qps, nil
}

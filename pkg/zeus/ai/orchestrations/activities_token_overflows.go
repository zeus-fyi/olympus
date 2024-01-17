package ai_platform_service_orchestrations

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

const (
	modelGpt4              = "gpt-4"
	modelGpt4TurboPreview  = "gpt-4-1106-preview"
	modelGpt4Vision        = "gpt-4-vision-preview"
	modelGpt432k           = "gpt-4-32k"
	modelGpt40613          = "gpt-4-0613"
	modelGpt432k0613       = "gpt-4-32k-0613"
	modelGpt35Turbo1106    = "gpt-3.5-turbo-1106"
	modelGpt35Turbo        = "gpt-3.5-turbo"
	modelGpt35Turbo16k     = "gpt-3.5-turbo-16k"
	modelGpt35TurboInstr   = "gpt-3.5-turbo-instruct"
	modelGpt35Turbo0613    = "gpt-3.5-turbo-0613"
	modelGpt35Turbo16k0613 = "gpt-3.5-turbo-16k-0613"
	modelGpt35Turbo0301    = "gpt-3.5-turbo-0301"
)

type PromptReduction struct {
	MarginBuffer                 float64                       `json:"marginBuffer,omitempty"`
	TokenOverflowStrategy        string                        `json:"tokenOverflowStrategy"`
	PromptReductionSearchResults *PromptReductionSearchResults `json:"promptReductionSearchResults,omitempty"`
	PromptReductionText          *PromptReductionText          `json:"promptReductionText,omitempty"`
}

type PromptReductionText struct {
	Model              string   `json:"model"`
	InPromptBody       string   `json:"inPromptBody"`
	OutPromptChunks    []string `json:"outPromptChunks,omitempty"`
	OutPromptTruncated string   `json:"outPromptTruncated,omitempty"`
}

type PromptReductionSearchResults struct {
	InPromptBody    string                           `json:"inPromptBody"`
	InSearchGroup   *hera_search.SearchResultGroup   `json:"inSearchGroup,omitempty"`
	OutSearchGroups []*hera_search.SearchResultGroup `json:"outSearchGroups,omitempty"`
}

const (
	OverflowStrategyTruncate = "truncate"
	OverflowStrategyDeduce   = "deduce"
)

func (z *ZeusAiPlatformActivities) TokenOverflowReduction(ctx context.Context, ou org_users.OrgUser, pr *PromptReduction) (*PromptReduction, error) {
	if pr == nil {
		return nil, nil
	}
	err := TokenOverflowSearchResults(ctx, pr)
	if err != nil {
		log.Err(err).Msg("TokenOverflowReduction: TokenOverflowSearchResults")
		return nil, err
	}
	err = TokenOverflowString(ctx, pr)
	if err != nil {
		log.Err(err).Msg("TokenOverflowReduction: TokenOverflowString")
		return nil, err
	}
	return pr, nil
}

func TokenOverflowSearchResults(ctx context.Context, pr *PromptReduction) error {
	if pr.PromptReductionSearchResults == nil || pr.PromptReductionSearchResults.InSearchGroup == nil || pr.PromptReductionSearchResults.InSearchGroup.SearchResults == nil {
		return nil
	}
	compressedSearchStr := hera_search.FormatSearchResultsV3(pr.PromptReductionSearchResults.InSearchGroup.SearchResults)
	needsReduction, err := CheckTokenContextMargin(ctx, pr.PromptReductionSearchResults.InSearchGroup.Model, compressedSearchStr, pr.MarginBuffer)
	if err != nil {
		log.Err(err).Msg("TokenOverflowSearchResults: CheckTokenContextMargin")
		return err
	}
	if !needsReduction {
		return nil
	}
	sr := pr.PromptReductionSearchResults.InSearchGroup.SearchResults
	switch pr.TokenOverflowStrategy {
	case OverflowStrategyDeduce:
		pr.PromptReductionSearchResults.OutSearchGroups = ChunkSearchResults(sr)
	case OverflowStrategyTruncate:
		pr.PromptReductionSearchResults.OutSearchGroups = TruncateSearchResults(sr)
	}
	return nil
}

func ChunkSearchResults(srs []hera_search.SearchResult) []*hera_search.SearchResultGroup {
	return nil
}

func TruncateSearchResults(srs []hera_search.SearchResult) []*hera_search.SearchResultGroup {
	return nil
}

func TokenOverflowString(ctx context.Context, pr *PromptReduction) error {
	if pr.PromptReductionText == nil || len(pr.PromptReductionText.InPromptBody) <= 0 {
		return nil
	}
	needsReduction, err := CheckTokenContextMargin(ctx, pr.PromptReductionSearchResults.InSearchGroup.Model, pr.PromptReductionText.InPromptBody, pr.MarginBuffer)
	if err != nil {
		log.Err(err).Msg("TokenOverflowSearchResults: CheckTokenContextMargin")
		return err
	}
	if !needsReduction {
		return nil
	}
	switch pr.TokenOverflowStrategy {
	case OverflowStrategyDeduce:
		pr.PromptReductionText.OutPromptChunks = ChunkPromptToSlices(pr.PromptReductionText.InPromptBody, pr.MarginBuffer)
	case OverflowStrategyTruncate:
		pr.PromptReductionText.OutPromptTruncated = TruncateString(pr.PromptReductionText.InPromptBody, pr.MarginBuffer)
	}
	return nil
}

func ChunkPromptToSlices(strIn string, marginBuffer float64) []string {
	return nil
}

func TruncateString(strIn string, marginBuffer float64) string {
	marginBuffer = validateMarginBufferLimits(marginBuffer)
	// Calculate the maximum length allowed based on the marginBuffer
	maxLength := int(float64(len(strIn)) * marginBuffer)
	// If the string is longer than the maximum length, truncate it
	if len(strIn) > maxLength {
		return strIn[:maxLength]
	}
	return strIn
}

func CheckTokenContextMargin(ctx context.Context, model, promptStr string, marginBuffer float64) (bool, error) {
	tokenLimit := GetModelTokenContextLimit(model)
	if tokenLimit == 0 {
		return false, fmt.Errorf("CheckTokenContextMargin: missing model in search group")
	}
	tokenEstimate, err := GetTokenCountEstimate(ctx, model, promptStr)
	if err != nil {
		log.Err(err).Msg("TokenOverflowReduction: GetTokenCountEstimate")
		return false, err
	}
	marginBuffer = validateMarginBufferLimits(marginBuffer)
	// Calculate the threshold using the margin buffer
	threshold := int(float64(tokenLimit) * marginBuffer)
	return tokenEstimate > threshold, nil
}

func validateMarginBufferLimits(marginBuffer float64) float64 {
	if marginBuffer < 0.1 {
		return 0.5
	}
	if marginBuffer > 0.80 {
		return 0.80
	}
	return marginBuffer
}

func GetModelTokenContextLimit(m string) int {
	switch m {
	case modelGpt4Vision, modelGpt4TurboPreview:
		return 128000
	case modelGpt4, modelGpt40613:
		return 8192
	case modelGpt432k, modelGpt432k0613:
		return 32768
	case modelGpt35Turbo1106, modelGpt35Turbo16k, modelGpt35Turbo16k0613:
		return 16385
	case modelGpt35Turbo, modelGpt35TurboInstr, modelGpt35Turbo0613, modelGpt35Turbo0301:
		return 4096
	default:
		return 0 // or some default value if model not listed
	}
}

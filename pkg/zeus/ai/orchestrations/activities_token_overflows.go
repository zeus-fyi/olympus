package ai_platform_service_orchestrations

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
)

const (
	modelGpt4              = "gpt-4"
	modelGpt4JanPreview    = "gpt-4-0125-preview"
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
	modelGpt35JanPreview   = "gpt-3.5-turbo-0125"
)

type PromptReduction struct {
	MarginBuffer          float64 `json:"marginBuffer,omitempty"`
	Model                 string  `json:"model"`
	TokenOverflowStrategy string  `json:"tokenOverflowStrategy"`

	AIWorkflowAnalysisResults    []artemis_orchestrations.AIWorkflowAnalysisResult `json:"dataInAnalysisResults,omitempty"`
	DataInAnalysisAggregation    []InputDataAnalysisToAgg                          `json:"dataInAnalysisAggregation,omitempty"`
	PromptReductionSearchResults *PromptReductionSearchResults                     `json:"promptReductionSearchResults,omitempty"`
	PromptReductionText          *PromptReductionText                              `json:"promptReductionText,omitempty"`
}

type PromptReductionText struct {
	InPromptSystem     string   `json:"inPromptSystem"`
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

func (z *ZeusAiPlatformActivities) TokenOverflowReduction(ctx context.Context, inputID int) (int, error) {
	if inputID <= 0 {
		return 0, nil
	}
	wio, werr := gws(ctx, inputID)
	if werr != nil {
		log.Err(werr).Msg("TokenOverflowReduction: failed to select workflow io")
		return 0, werr
	}
	if wio.PromptReduction == nil {
		return 0, nil
	}
	pr := wio.PromptReduction
	if pr.DataInAnalysisAggregation != nil {
		pr.PromptReductionSearchResults = &PromptReductionSearchResults{
			InSearchGroup: &hera_search.SearchResultGroup{
				SearchResults:         make([]hera_search.SearchResult, 0),
				ApiResponseResults:    make([]hera_search.SearchResult, 0),
				FilteredSearchResults: make([]hera_search.SearchResult, 0),
			},
		}
		for _, d := range pr.DataInAnalysisAggregation {
			if d.SearchResultGroup != nil && d.ChatCompletionQueryResponse != nil && d.ChatCompletionQueryResponse.JsonResponseResults != nil {
				payloadMaps := artemis_orchestrations.CreateMapInterfaceFromAssignedSchemaFields(d.ChatCompletionQueryResponse.JsonResponseResults)
				switch d.SearchResultGroup.PlatformName {
				case twitterPlatform, discordPlatform, redditPlatform, telegramPlatform:
					tmpMap := make(map[string]map[string]interface{})
					for ind, pv := range payloadMaps {
						for keyName, payloadValue := range payloadMaps[ind] {
							if keyName == "msg_id" {
								msgStrID, ok := payloadValue.(string)
								if ok {
									if tmpMap[msgStrID] == nil {
										tmpMap[msgStrID] = make(map[string]interface{})
									}
									tmpMap[msgStrID] = pv
								}
							}
						}
					}
					for _, sv := range d.SearchResultGroup.SearchResults {
						if sv.TwitterMetadata != nil && sv.TwitterMetadata.TweetStrID != "" {
							if item, ok := tmpMap[sv.TwitterMetadata.TweetStrID]; ok && item != nil {
								sv.WebResponse.Body = item
								pr.PromptReductionSearchResults.InSearchGroup.SearchResults = append(pr.PromptReductionSearchResults.InSearchGroup.SearchResults, sv)
							}
						} else if item, ok := tmpMap[fmt.Sprintf("%d", sv.UnixTimestamp)]; ok && item != nil {
							sv.WebResponse.Body = item
							pr.PromptReductionSearchResults.InSearchGroup.SearchResults = append(pr.PromptReductionSearchResults.InSearchGroup.SearchResults, sv)
						}
					}
				}
			} else if d.SearchResultGroup != nil && d.SearchResultGroup.SearchResults != nil {
				if pr.PromptReductionSearchResults == nil {
					pr.PromptReductionSearchResults = &PromptReductionSearchResults{
						InSearchGroup: &hera_search.SearchResultGroup{
							SearchResults:         make([]hera_search.SearchResult, 0),
							ApiResponseResults:    make([]hera_search.SearchResult, 0),
							FilteredSearchResults: make([]hera_search.SearchResult, 0),
						},
					}
				}
				if d.SearchResultGroup.FilteredSearchResults != nil {
					pr.PromptReductionSearchResults.InSearchGroup.FilteredSearchResults = append(pr.PromptReductionSearchResults.InSearchGroup.FilteredSearchResults, d.SearchResultGroup.FilteredSearchResults...)
				}
				if d.SearchResultGroup.ApiResponseResults != nil {
					pr.PromptReductionSearchResults.InSearchGroup.ApiResponseResults = append(pr.PromptReductionSearchResults.InSearchGroup.ApiResponseResults, d.SearchResultGroup.ApiResponseResults...)
				}
				pr.PromptReductionSearchResults.InSearchGroup.SearchResults = append(pr.PromptReductionSearchResults.InSearchGroup.SearchResults, d.SearchResultGroup.SearchResults...)
			}
		}
	} else {
		for _, wr := range pr.AIWorkflowAnalysisResults {
			sv, err := artemis_orchestrations.GenerateContentTextFromOpenAIResp([]artemis_orchestrations.AIWorkflowAnalysisResult{wr})
			if err != nil {
				log.Err(err).Msg("TokenOverflowReduction: GenerateContentTextFromOpenAIResp")
				continue
			}
			hs := hera_search.SearchResult{
				Value: sv,
			}
			pr.PromptReductionSearchResults.InSearchGroup.SearchResults = append(pr.PromptReductionSearchResults.InSearchGroup.SearchResults, hs)
		}
	}
	log.Info().Interface("pr.MarginBuffer", pr.MarginBuffer).Msg("TokenOverflowReduction")
	err := TokenOverflowSearchResults(ctx, pr)
	if err != nil {
		log.Err(err).Msg("TokenOverflowReduction: TokenOverflowSearchResults")
		return 0, err
	}
	err = TokenOverflowString(ctx, pr)
	if err != nil {
		log.Err(err).Msg("TokenOverflowReduction: TokenOverflowString")
		return 0, err
	}

	tmp := pr.PromptReductionSearchResults
	if tmp != nil && tmp.OutSearchGroups != nil && len(tmp.OutSearchGroups) > 0 {
		log.Info().Interface("pr.TokenOverflowStrategy", pr.TokenOverflowStrategy).Interface("len(tmp.OutSearchGroups)", len(tmp.OutSearchGroups)).Msg("TokenOverflowReductioDone")
	}
	if pr.PromptReductionText != nil {
		log.Info().Interface("pr.TokenOverflowStrategy", pr.TokenOverflowStrategy).Interface("pr.PromptReductionText.OutPromptChunks", len(pr.PromptReductionText.OutPromptChunks)).Msg("TokenOverflowReductionDone")
	}
	wio.PromptReduction = pr
	_, werr = sws(ctx, &wio)
	if werr != nil {
		log.Err(werr).Msg("TokenOverflowReduction: failed to save workflow io")
		return 0, werr
	}
	chunkIterator := getChunkIteratorLen(pr)
	return chunkIterator, nil
}

func TokenOverflowSearchResults(ctx context.Context, pr *PromptReduction) error {
	if pr.PromptReductionSearchResults == nil || pr.PromptReductionSearchResults.InSearchGroup == nil {
		log.Info().Msg("TokenOverflowSearchResults: no search results")
		return nil
	}
	if pr.PromptReductionSearchResults.InSearchGroup.ApiResponseResults == nil && pr.PromptReductionSearchResults.InSearchGroup.SearchResults == nil {
		log.Info().Msg("TokenOverflowSearchResults: no search results or api responses")
		return nil
	}

	switch pr.TokenOverflowStrategy {
	case OverflowStrategyDeduce:
		err := ChunkSearchResults(ctx, pr)
		if err != nil {
			log.Err(err).Msg("TokenOverflowSearchResults: ChunkSearchResults")
			return err
		}
	case OverflowStrategyTruncate:
		err := TruncateSearchResults(ctx, pr)
		if err != nil {
			log.Err(err).Msg("TokenOverflowSearchResults: ChunkSearchResults")
			return err
		}
	}
	return nil
}

func TruncateSearchResults(ctx context.Context, pr *PromptReduction) error {
	err := ChunkSearchResults(ctx, pr)
	if err != nil {
		log.Err(err).Msg("TruncateSearchResults: ChunkSearchResults")
		return err
	}
	if pr.PromptReductionSearchResults != nil && len(pr.PromptReductionSearchResults.OutSearchGroups) > 0 {
		// Keep only the first element and remove the rest.
		pr.PromptReductionSearchResults.OutSearchGroups = pr.PromptReductionSearchResults.OutSearchGroups[:1]
	}
	return nil
}

func ChunkSearchResults(ctx context.Context, pr *PromptReduction) error {
	marginBuffer := validateMarginBufferLimits(pr.MarginBuffer)
	model := pr.Model
	var compressedSearchStr string
	if pr.PromptReductionSearchResults.InSearchGroup.ApiResponseResults != nil && len(pr.PromptReductionSearchResults.InSearchGroup.ApiResponseResults) > 0 {
		compressedSearchStr += hera_search.FormatApiSearchResultSliceToString(pr.PromptReductionSearchResults.InSearchGroup.ApiResponseResults)
	}
	if pr.PromptReductionSearchResults.InSearchGroup.SearchResults != nil {
		compressedSearchStr += hera_search.FormatSearchResultsV5(pr.PromptReductionSearchResults.InSearchGroup.SearchResults)
	}
	if pr.PromptReductionText != nil {
		compressedSearchStr += pr.PromptReductionText.InPromptSystem
		compressedSearchStr += pr.PromptReductionText.InPromptBody
	}

	needsReduction, tokenEstimate, err := CheckTokenContextMargin(ctx, model, compressedSearchStr, marginBuffer)
	if err != nil {
		log.Err(err).Interface("tokenEstimate", tokenEstimate).Msg("TokenOverflowSearchResults: CheckTokenContextMargin")
		return err
	}
	if !needsReduction {
		pr.PromptReductionSearchResults.InSearchGroup.SearchResultChunkTokenEstimate = &tokenEstimate
		pr.PromptReductionSearchResults.OutSearchGroups = []*hera_search.SearchResultGroup{
			pr.PromptReductionSearchResults.InSearchGroup,
		}
		return nil
	}
	totalSearchResults := pr.PromptReductionSearchResults.InSearchGroup.SearchResults

	if len(totalSearchResults) <= 0 && len(compressedSearchStr) > 0 && needsReduction {
		// Treat compressedSearchStr as if it was an input string that can be chunked
		marginBuffer = validateMarginBufferLimits(pr.MarginBuffer)
		chunks, cerr := ChunkPromptToSlices(ctx, pr.Model, compressedSearchStr, marginBuffer)
		if cerr != nil {
			log.Err(cerr).Msg("TokenOverflowSearchResults: ChunkPromptToSlices for compressedSearchStr")
			return cerr
		}

		// Assuming that ChunkPromptToSlices does not only chunk but also ensures each chunk is within token limits
		if len(chunks) > 0 {
			// Update the PromptReductionText to reflect the chunking of compressedSearchStr
			pr.PromptReductionText = &PromptReductionText{
				InPromptBody:    compressedSearchStr, // Original compressed string
				OutPromptChunks: chunks,              // Chunks after processing
			}
		} else {
			// Fallback to original string if no chunks were created (should not happen due to checks)
			pr.PromptReductionText = &PromptReductionText{
				InPromptBody:       compressedSearchStr,
				OutPromptTruncated: compressedSearchStr, // Or handle accordingly
			}
		}
		return nil
	}

	splitIteration := 2
	for needsReduction && (splitIteration < len(totalSearchResults)) {
		chunks := splitSliceIntoChunks(totalSearchResults, splitIteration)
		var tokenEstimates []int
		needsReduction, tokenEstimates, err = validateChunkTokenLimits(ctx, model, marginBuffer, chunks)
		if err != nil {
			log.Err(err).Msg("TokenOverflowSearchResults: validateChunkTokenLimits")
			return err
		}
		if !needsReduction {
			pr.PromptReductionSearchResults.OutSearchGroups = make([]*hera_search.SearchResultGroup, len(chunks))
			for i, chunk := range chunks {
				pr.PromptReductionSearchResults.OutSearchGroups[i] = createChunk(pr.PromptReductionSearchResults.InSearchGroup, chunk)
				pr.PromptReductionSearchResults.OutSearchGroups[i].SearchResultChunkTokenEstimate = &tokenEstimates[i]
			}
			return nil
		}
		splitIteration++
	}

	return fmt.Errorf("TokenOverflowSearchResults: failed to reduce search results")
}

func validateChunkTokenLimits(ctx context.Context, model string, marginBuffer float64, srs [][]hera_search.SearchResult) (bool, []int, error) {
	var tokenEstimates []int
	for _, chunk := range srs {
		compressedSearchStr := hera_search.FormatSearchResultsV3(chunk)
		needsReduction, tokenEstimate, err := CheckTokenContextMargin(ctx, model, compressedSearchStr, marginBuffer)
		if err != nil {
			log.Err(err).Interface("tokenEstimate", tokenEstimate).Msg("TokenOverflowSearchResults: CheckTokenContextMargin")
			return false, nil, err
		}
		tokenEstimates = append(tokenEstimates, tokenEstimate)
		if needsReduction {
			return true, nil, nil
		}
	}
	if len(tokenEstimates) != len(srs) {
		return false, nil, fmt.Errorf("validateChunkTokenLimits: tokenEstimates length mismatch")
	}
	return false, tokenEstimates, nil
}
func splitSliceIntoChunks[T any](s []T, chunkCount int) [][]T {
	if chunkCount <= 0 {
		// Handle invalid chunk count
		return nil
	}
	length := len(s)
	var chunks [][]T
	chunkSize := length / chunkCount
	remainder := length % chunkCount
	start := 0
	for i := 0; i < chunkCount; i++ {
		end := start + chunkSize
		if remainder > 0 {
			end++ // Distribute the remainder among the first few chunks
			remainder--
		}
		// Slice the chunk
		if end > length {
			end = length
		}
		chunks = append(chunks, s[start:end])
		start = end
	}
	return chunks
}

func createChunk(originalGroup *hera_search.SearchResultGroup, chunk []hera_search.SearchResult) *hera_search.SearchResultGroup {
	newGroup := *originalGroup
	newGroup.SearchResults = chunk
	return &newGroup
}

func TokenOverflowString(ctx context.Context, pr *PromptReduction) error {
	if pr.PromptReductionText == nil || len(pr.PromptReductionText.InPromptBody) <= 0 {
		return nil
	}
	model := pr.Model
	margin := validateMarginBufferLimits(pr.MarginBuffer)
	needsReduction, _, err := CheckTokenContextMargin(ctx, model, pr.PromptReductionText.InPromptBody, margin)
	if err != nil {
		log.Err(err).Msg("TokenOverflowSearchResults: CheckTokenContextMargin")
		return err
	}
	if !needsReduction {
		pr.PromptReductionText.OutPromptTruncated = pr.PromptReductionText.InPromptBody
		pr.PromptReductionText.OutPromptChunks = []string{pr.PromptReductionText.InPromptBody}
		return nil
	}
	var chunks []string
	switch pr.TokenOverflowStrategy {
	case OverflowStrategyDeduce:
		chunks, err = ChunkPromptToSlices(ctx, model, pr.PromptReductionText.InPromptBody, margin)
		if err != nil {
			log.Err(err).Msg("TokenOverflowSearchResults: ChunkPromptToSlices")
			return err
		}
		pr.PromptReductionText.OutPromptChunks = chunks
	case OverflowStrategyTruncate:
		chunks, err = ChunkPromptToSlices(ctx, model, pr.PromptReductionText.InPromptBody, margin)
		if err != nil {
			log.Err(err).Msg("TokenOverflowSearchResults: ChunkPromptToSlices")
			return err
		}
		if len(chunks) > 0 {
			pr.PromptReductionText.OutPromptTruncated = chunks[0]
		}
	}
	return nil
}

func ChunkPromptToSlices(ctx context.Context, model, strIn string, marginBuffer float64) ([]string, error) {
	splitIteration := 2
	for {
		chunks := splitStringIntoChunks(strIn, splitIteration)
		allChunksValid := true
		var validChunks []string
		for _, chunk := range chunks {
			needsReduction, _, err := CheckTokenContextMargin(ctx, model, chunk, marginBuffer)
			if err != nil {
				log.Err(err).Msg("TokenOverflowSearchResults: CheckTokenContextMargin")
				return nil, err
			}
			if needsReduction {
				allChunksValid = false
				break
			}
			validChunks = append(validChunks, chunk)
		}
		if allChunksValid {
			return validChunks, nil
		}
		splitIteration++
		// Avoid infinite loop: stop if splitIteration exceeds the length of the string
		if splitIteration > len(strIn) {
			return nil, fmt.Errorf("unable to reduce token overflow with current margin buffer")
		}
	}
}

func splitStringIntoChunks(str string, chunkCount int) []string {
	if chunkCount <= 0 {
		return nil
	}
	var chunks []string
	chunkSize := len(str) / chunkCount
	remainder := len(str) % chunkCount
	start := 0
	for i := 0; i < chunkCount; i++ {
		end := start + chunkSize
		if remainder > 0 {
			end++
			remainder--
		}
		if end > len(str) {
			end = len(str)
		}
		chunks = append(chunks, str[start:end])
		start = end
	}
	return chunks
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

func CheckTokenContextMargin(ctx context.Context, model, promptStr string, marginBuffer float64) (bool, int, error) {
	tokenLimit := GetModelTokenContextLimit(model)
	if tokenLimit == 0 {
		return false, -1, fmt.Errorf("CheckTokenContextMargin: missing model in search group")
	}
	tokenEstimate, err := GetTokenCountEstimate(ctx, model, promptStr)
	if err != nil {
		log.Err(err).Msg("TokenOverflowReduction: GetTokenCountEstimate")
		return false, tokenEstimate, err
	}
	if tokenEstimate < 0 {
		return false, tokenEstimate, fmt.Errorf("CheckTokenContextMargin: failed to estimate token count")
	}
	marginBuffer = validateMarginBufferLimits(marginBuffer)
	// Calculate the threshold using the margin buffer
	threshold := int(float64(tokenLimit) * marginBuffer)
	return tokenEstimate > threshold, tokenEstimate, nil
}

func validateMarginBufferLimits(marginBuffer float64) float64 {
	if marginBuffer < 0.01 {
		return 0.5
	}
	if marginBuffer >= 0.01 && marginBuffer < 0.2 {
		return 0.2
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
		return 4096 // or some default value if model not listed
	}
}

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

func overflowSetup(ctx context.Context, cp *MbChildSubProcessParams, promptExt *PromptReduction) (*PromptReduction, *WorkflowStageIO, error) {
	var pr *PromptReduction
	wio, werr := gs3wfs(ctx, cp)
	if werr != nil {
		log.Err(werr).Msg("TokenOverflowReduction: overflowSetup failed to select workflow io")
		return nil, nil, werr
	}
	if wio.WorkflowStageInfo.PromptReduction != nil && wio.WorkflowStageInfo.PromptReduction.DataInAnalysisAggregation != nil {
		for _, d := range wio.WorkflowStageInfo.PromptReduction.DataInAnalysisAggregation {
			if d.ChatCompletionQueryResponse != nil && d.ChatCompletionQueryResponse.RegexSearchResults != nil {
				wio.PromptReduction.PromptReductionSearchResults = &PromptReductionSearchResults{
					InSearchGroup: &hera_search.SearchResultGroup{
						RegexSearchResults: d.ChatCompletionQueryResponse.RegexSearchResults,
					},
				}
			}
		}
	}
	if wio.PromptReduction == nil {
		return nil, nil, nil
	}
	pr = wio.PromptReduction
	if promptExt != nil && promptExt.PromptReductionText != nil {
		pr.PromptReductionText = promptExt.PromptReductionText
	}
	wio.PromptReduction = pr
	if pr.DataInAnalysisAggregation != nil {
		pr.PromptReductionSearchResults = &PromptReductionSearchResults{
			InSearchGroup: &hera_search.SearchResultGroup{
				SearchResults:         make([]hera_search.SearchResult, 0),
				ApiResponseResults:    make([]hera_search.SearchResult, 0),
				FilteredSearchResults: make([]hera_search.SearchResult, 0),
				RegexSearchResults:    make([]hera_search.SearchResult, 0),
			},
		}
		for _, d := range pr.DataInAnalysisAggregation {
			if d.ChatCompletionQueryResponse != nil && pr != nil && d.ChatCompletionQueryResponse.RegexSearchResults != nil {
				log.Info().Msg("TokenOverflowReduction: ChatCompletionQueryResponse.RegexSearchResults")
				sk := &hera_search.SearchResultGroup{
					RegexSearchResults: d.ChatCompletionQueryResponse.RegexSearchResults,
				}
				pr.PromptReductionSearchResults = &PromptReductionSearchResults{
					InSearchGroup: sk,
				}
			} else if d.SearchResultGroup != nil && d.ChatCompletionQueryResponse != nil && d.ChatCompletionQueryResponse.JsonResponseResults != nil {
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
							RegexSearchResults:    make([]hera_search.SearchResult, 0),
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
	return pr, wio, nil
}
